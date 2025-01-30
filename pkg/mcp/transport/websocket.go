package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/protocol"
	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/server"
	"github.com/gorilla/websocket"
)

// WebSocketTransport implements a WebSocket-based transport for MCP
type WebSocketTransport struct {
	session  *server.Session
	upgrader websocket.Upgrader
	clients  map[string]*websocket.Conn
	mu       sync.RWMutex
	opts     Options
	srv      *http.Server
}

// NewWebSocketTransport creates a new WebSocket transport
func NewWebSocketTransport(session *server.Session, options ...Option) HTTPTransport {
	opts := defaultOptions()
	for _, opt := range options {
		opt(&opts)
	}

	return &WebSocketTransport{
		session: session,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in this example
			},
		},
		clients: make(map[string]*websocket.Conn),
		opts:    opts,
	}
}

// Start starts the WebSocket transport on the default address
func (t *WebSocketTransport) Start() error {
	return t.StartHTTP(t.opts.Address)
}

// StartHTTP starts the WebSocket transport on the given address
func (t *WebSocketTransport) StartHTTP(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", t.handleWebSocket)

	t.srv = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return t.srv.ListenAndServe()
}

// Stop stops the transport
func (t *WebSocketTransport) Stop(ctx context.Context) error {
	t.mu.Lock()
	for _, conn := range t.clients {
		conn.Close()
	}
	t.clients = make(map[string]*websocket.Conn)
	t.mu.Unlock()

	if t.srv != nil {
		return t.srv.Shutdown(ctx)
	}
	return nil
}

// handleWebSocket handles WebSocket connections
func (t *WebSocketTransport) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := t.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade connection: %v\n", err)
		return
	}

	clientID := r.RemoteAddr
	t.mu.Lock()
	t.clients[clientID] = conn
	t.mu.Unlock()

	defer func() {
		conn.Close()
		t.mu.Lock()
		delete(t.clients, clientID)
		t.mu.Unlock()
	}()

	for {
		// Read message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Error reading message: %v\n", err)
			}
			return
		}

		if messageType != websocket.TextMessage {
			continue
		}

		// Parse the message
		var msg struct {
			JSONRPC string              `json:"jsonrpc"`
			ID      *protocol.RequestID `json:"id,omitempty"`
			Method  string              `json:"method"`
			Params  json.RawMessage     `json:"params,omitempty"`
		}
		if err := json.Unmarshal(message, &msg); err != nil {
			t.writeError(conn, nil, -32700, "Parse error", err)
			continue
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
			t.handleRequest(conn, req)
		} else {
			// This is a notification
			notif := &protocol.JSONRPCNotification{
				JSONRPC: msg.JSONRPC,
				Method:  msg.Method,
				Params:  msg.Params,
			}
			t.handleNotification(conn, notif)
		}
	}
}

// handleRequest processes a request and writes the response
func (t *WebSocketTransport) handleRequest(conn *websocket.Conn, req *protocol.JSONRPCRequest) {
	resp, err := t.session.HandleRequest(req)
	if err != nil {
		t.writeErrorWithID(conn, req.ID, -32603, "Internal error", err)
		return
	}

	t.writeResponse(conn, resp)
}

// handleNotification processes a notification
func (t *WebSocketTransport) handleNotification(conn *websocket.Conn, notif *protocol.JSONRPCNotification) {
	if err := t.session.HandleNotification(notif); err != nil {
		// Log the error but don't send a response for notifications
		fmt.Printf("Error handling notification: %v\n", err)
	}
}

// writeResponse writes a JSON-RPC response to the WebSocket connection
func (t *WebSocketTransport) writeResponse(conn *websocket.Conn, resp *protocol.JSONRPCResponse) {
	if err := conn.WriteJSON(resp); err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}
}

// writeError writes a JSON-RPC error response with no ID
func (t *WebSocketTransport) writeError(conn *websocket.Conn, id *protocol.RequestID, code int, message string, err error) {
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

	if err := conn.WriteJSON(errResp); err != nil {
		fmt.Printf("Error writing error response: %v\n", err)
	}
}

// writeErrorWithID writes a JSON-RPC error response with a specific ID
func (t *WebSocketTransport) writeErrorWithID(conn *websocket.Conn, id protocol.RequestID, code int, message string, err error) {
	errResp := &protocol.JSONRPCError{
		JSONRPC: "2.0",
		ID:      id,
		Error: protocol.ErrorData{
			Code:    code,
			Message: message,
			Data:    err.Error(),
		},
	}

	if err := conn.WriteJSON(errResp); err != nil {
		fmt.Printf("Error writing error response: %v\n", err)
	}
}

// SendNotification sends a notification to all connected clients
func (t *WebSocketTransport) SendNotification(method string, params interface{}) error {
	notif := &protocol.JSONRPCNotification{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	var lastErr error
	for _, conn := range t.clients {
		if err := conn.WriteJSON(notif); err != nil {
			lastErr = err
			fmt.Printf("Error sending notification to client: %v\n", err)
		}
	}

	return lastErr
}
