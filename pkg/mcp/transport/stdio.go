package transport

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/protocol"
	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/server"
)

// StdioTransport implements a stdio-based transport for MCP
type StdioTransport struct {
	session *server.Session
	reader  *bufio.Reader
	writer  *bufio.Writer
	mu      sync.Mutex
	opts    Options
}

// NewStdioTransport creates a new stdio transport
func NewStdioTransport(session *server.Session, options ...Option) Transport {
	opts := defaultOptions()
	for _, opt := range options {
		opt(&opts)
	}

	return &StdioTransport{
		session: session,
		reader:  bufio.NewReader(os.Stdin),
		writer:  bufio.NewWriter(os.Stdout),
		opts:    opts,
	}
}

// Start starts the transport
func (t *StdioTransport) Start() error {
	for {
		// Read a line from stdin
		line, err := t.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("failed to read from stdin: %w", err)
		}

		// Parse the message
		var msg struct {
			JSONRPC string              `json:"jsonrpc"`
			ID      *protocol.RequestID `json:"id,omitempty"`
			Method  string              `json:"method"`
			Params  json.RawMessage     `json:"params,omitempty"`
		}
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			t.writeError(nil, 0, "Parse error", err)
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
			t.handleRequest(req)
		} else {
			// This is a notification
			notif := &protocol.JSONRPCNotification{
				JSONRPC: msg.JSONRPC,
				Method:  msg.Method,
				Params:  msg.Params,
			}
			t.handleNotification(notif)
		}
	}
}

// handleRequest processes a request and writes the response
func (t *StdioTransport) handleRequest(req *protocol.JSONRPCRequest) {
	resp, err := t.session.HandleRequest(req)
	if err != nil {
		t.writeErrorWithID(req.ID, -32603, "Internal error", err)
		return
	}

	t.writeResponse(resp)
}

// handleNotification processes a notification
func (t *StdioTransport) handleNotification(notif *protocol.JSONRPCNotification) {
	if err := t.session.HandleNotification(notif); err != nil {
		// Log the error but don't send a response for notifications
		fmt.Fprintf(os.Stderr, "Error handling notification: %v\n", err)
	}
}

// writeResponse writes a JSON-RPC response to stdout
func (t *StdioTransport) writeResponse(resp *protocol.JSONRPCResponse) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := json.NewEncoder(t.writer).Encode(resp); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing response: %v\n", err)
		return
	}

	t.writer.Flush()
}

// writeError writes a JSON-RPC error response to stdout with no ID
func (t *StdioTransport) writeError(id *protocol.RequestID, code int, message string, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

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

	if err := json.NewEncoder(t.writer).Encode(errResp); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing error response: %v\n", err)
		return
	}

	t.writer.Flush()
}

// writeErrorWithID writes a JSON-RPC error response to stdout with a specific ID
func (t *StdioTransport) writeErrorWithID(id protocol.RequestID, code int, message string, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	errResp := &protocol.JSONRPCError{
		JSONRPC: "2.0",
		ID:      id,
		Error: protocol.ErrorData{
			Code:    code,
			Message: message,
			Data:    err.Error(),
		},
	}

	if err := json.NewEncoder(t.writer).Encode(errResp); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing error response: %v\n", err)
		return
	}

	t.writer.Flush()
}

// SendNotification sends a notification to the client
func (t *StdioTransport) SendNotification(method string, params interface{}) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	notif := &protocol.JSONRPCNotification{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	if err := json.NewEncoder(t.writer).Encode(notif); err != nil {
		return fmt.Errorf("failed to write notification: %w", err)
	}

	return t.writer.Flush()
}
