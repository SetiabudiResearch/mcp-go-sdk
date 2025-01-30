package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/protocol"
	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/server"
)

// SSETransport implements a Server-Sent Events transport for MCP
type SSETransport struct {
	session *server.Session
	clients map[string]chan []byte
	mu      sync.RWMutex
	opts    Options
	srv     *http.Server
}

// NewSSETransport creates a new SSE transport
func NewSSETransport(session *server.Session, options ...Option) HTTPTransport {
	opts := defaultOptions()
	for _, opt := range options {
		opt(&opts)
	}

	return &SSETransport{
		session: session,
		clients: make(map[string]chan []byte, opts.BufferSize),
		opts:    opts,
	}
}

// Start starts the SSE transport on the default address
func (t *SSETransport) Start() error {
	return t.StartHTTP(t.opts.Address)
}

// StartHTTP starts the SSE transport on the given address
func (t *SSETransport) StartHTTP(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/events", t.handleSSE)
	mux.HandleFunc("/", t.handleRequest)

	t.srv = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return t.srv.ListenAndServe()
}

// Stop stops the transport
func (t *SSETransport) Stop(ctx context.Context) error {
	if t.srv != nil {
		return t.srv.Shutdown(ctx)
	}
	return nil
}

// handleSSE handles SSE connections
func (t *SSETransport) handleSSE(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a channel for this client
	clientChan := make(chan []byte, t.opts.BufferSize)
	clientID := r.RemoteAddr

	// Register the client
	t.mu.Lock()
	t.clients[clientID] = clientChan
	t.mu.Unlock()

	// Clean up when the connection is closed
	defer func() {
		t.mu.Lock()
		delete(t.clients, clientID)
		close(clientChan)
		t.mu.Unlock()
	}()

	// Create a context that's cancelled when the client disconnects
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Start the event loop
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-clientChan:
			if !ok {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}

// handleRequest handles incoming JSON-RPC requests
func (t *SSETransport) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request
	var msg struct {
		JSONRPC string              `json:"jsonrpc"`
		ID      *protocol.RequestID `json:"id,omitempty"`
		Method  string              `json:"method"`
		Params  json.RawMessage     `json:"params,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		t.writeError(w, nil, -32700, "Parse error", err)
		return
	}

	// Handle the message
	if msg.ID != nil {
		// This is a request
		req := &protocol.JSONRPCRequest{
			JSONRPC: msg.JSONRPC,
			ID:      *msg.ID,
			Method:  msg.Method,
			Params:  msg.Params,
		}
		t.handleJSONRPCRequest(w, req)
	} else {
		// This is a notification
		notif := &protocol.JSONRPCNotification{
			JSONRPC: msg.JSONRPC,
			Method:  msg.Method,
			Params:  msg.Params,
		}
		t.handleJSONRPCNotification(w, notif)
	}
}

// handleJSONRPCRequest processes a JSON-RPC request and writes the response
func (t *SSETransport) handleJSONRPCRequest(w http.ResponseWriter, req *protocol.JSONRPCRequest) {
	resp, err := t.session.HandleRequest(req)
	if err != nil {
		t.writeErrorWithID(w, req.ID, -32603, "Internal error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleJSONRPCNotification processes a JSON-RPC notification
func (t *SSETransport) handleJSONRPCNotification(w http.ResponseWriter, notif *protocol.JSONRPCNotification) {
	if err := t.session.HandleNotification(notif); err != nil {
		// Log the error but don't send a response for notifications
		fmt.Printf("Error handling notification: %v\n", err)
	}
	w.WriteHeader(http.StatusNoContent)
}

// writeError writes a JSON-RPC error response with no ID
func (t *SSETransport) writeError(w http.ResponseWriter, id *protocol.RequestID, code int, message string, err error) {
	errResp := &protocol.JSONRPCError{
		JSONRPC: "2.0",
		Error: protocol.ErrorData{
			Code:    code,
			Message: message,
			Data:    err.Error(),
		},
	}

	if id != nil {
		errResp.ID = *id
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errResp)
}

// writeErrorWithID writes a JSON-RPC error response with a specific ID
func (t *SSETransport) writeErrorWithID(w http.ResponseWriter, id protocol.RequestID, code int, message string, err error) {
	errResp := &protocol.JSONRPCError{
		JSONRPC: "2.0",
		ID:      id,
		Error: protocol.ErrorData{
			Code:    code,
			Message: message,
			Data:    err.Error(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errResp)
}

// SendNotification sends a notification to all connected clients
func (t *SSETransport) SendNotification(method string, params interface{}) error {
	notif := &protocol.JSONRPCNotification{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(notif)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, ch := range t.clients {
		select {
		case ch <- data:
		default:
			// Skip clients that aren't ready to receive
		}
	}

	return nil
}
