package server

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/protocol"
)

// Server represents an MCP server instance
type Server struct {
	name         string
	capabilities protocol.ServerCapabilities
	info         protocol.Implementation
	session      *Session
	tools        map[string]Tool
	resources    map[string]Resource
	prompts      map[string]Prompt
	mu           sync.RWMutex
}

// Session represents a connection between client and server
type Session struct {
	ctx          context.Context
	cancel       context.CancelFunc
	server       *Server
	initialized  bool
	capabilities protocol.ClientCapabilities
	clientInfo   protocol.Implementation
	mu           sync.RWMutex
}

// Tool represents a function that can be called by the LLM
type Tool struct {
	Handler     interface{}
	Description string
	IsAsync     bool
}

// Resource represents a data source that can be accessed by the LLM
type Resource struct {
	Handler     interface{}
	Description string
	Pattern     string
}

// Prompt represents a template for LLM interactions
type Prompt struct {
	Handler     interface{}
	Description string
}

// NewServer creates a new MCP server instance
func NewServer(name string, opts ...ServerOption) *Server {
	s := &Server{
		name:      name,
		tools:     make(map[string]Tool),
		resources: make(map[string]Resource),
		prompts:   make(map[string]Prompt),
		info: protocol.Implementation{
			Name:    name,
			Version: protocol.LatestProtocolVersion,
		},
		capabilities: protocol.ServerCapabilities{
			Tools:     &protocol.ToolsCapability{},
			Resources: &protocol.ResourcesCapability{},
			Prompts:   &protocol.PromptsCapability{},
			Logging:   &protocol.LoggingCapability{},
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// NewSession creates a new session for a client connection
func NewSession(ctx context.Context, server *Server) *Session {
	ctx, cancel := context.WithCancel(ctx)
	return &Session{
		ctx:    ctx,
		cancel: cancel,
		server: server,
	}
}

// HandleRequest processes an incoming JSON-RPC request
func (s *Session) HandleRequest(req *protocol.JSONRPCRequest) (*protocol.JSONRPCResponse, error) {
	s.mu.RLock()
	initialized := s.initialized
	s.mu.RUnlock()

	// Handle initialization request
	if req.Method == "initialize" {
		if initialized {
			return nil, fmt.Errorf("server already initialized")
		}
		return s.handleInitialize(req)
	}

	// All other requests require initialization
	if !initialized {
		return nil, fmt.Errorf("server not initialized")
	}

	// Handle other requests based on method
	switch req.Method {
	case "ping":
		return s.handlePing(req)
	case "tools/list":
		return s.handleListTools(req)
	case "tools/call":
		return s.handleCallTool(req)
	case "resources/list":
		return s.handleListResources(req)
	case "resources/read":
		return s.handleReadResource(req)
	case "prompts/list":
		return s.handleListPrompts(req)
	case "prompts/get":
		return s.handleGetPrompt(req)
	default:
		return nil, fmt.Errorf("unknown method: %s", req.Method)
	}
}

// HandleNotification processes an incoming JSON-RPC notification
func (s *Session) HandleNotification(notif *protocol.JSONRPCNotification) error {
	s.mu.RLock()
	initialized := s.initialized
	s.mu.RUnlock()

	if !initialized && notif.Method != "notifications/initialized" {
		return fmt.Errorf("server not initialized")
	}

	switch notif.Method {
	case "notifications/initialized":
		return s.handleInitialized(notif)
	case "notifications/cancelled":
		return s.handleCancelled(notif)
	default:
		return fmt.Errorf("unknown notification method: %s", notif.Method)
	}
}

// handleInitialize processes the initialize request
func (s *Session) handleInitialize(req *protocol.JSONRPCRequest) (*protocol.JSONRPCResponse, error) {
	var params protocol.InitializeRequestParams
	if err := json.Unmarshal(req.Params.(json.RawMessage), &params); err != nil {
		return nil, fmt.Errorf("invalid initialize params: %w", err)
	}

	s.mu.Lock()
	s.capabilities = params.Capabilities
	s.clientInfo = params.ClientInfo
	s.initialized = true
	s.mu.Unlock()

	result := protocol.InitializeResult{
		ProtocolVersion: protocol.LatestProtocolVersion,
		Capabilities:    s.server.capabilities,
		ServerInfo:      s.server.info,
	}

	return &protocol.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}, nil
}

// handleInitialized processes the initialized notification
func (s *Session) handleInitialized(notif *protocol.JSONRPCNotification) error {
	// Nothing to do here, but we could trigger any post-initialization tasks
	return nil
}

// handlePing processes ping requests
func (s *Session) handlePing(req *protocol.JSONRPCRequest) (*protocol.JSONRPCResponse, error) {
	return &protocol.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  struct{}{},
	}, nil
}

// handleCancelled processes cancellation notifications
func (s *Session) handleCancelled(notif *protocol.JSONRPCNotification) error {
	var params protocol.CancelledNotificationParams
	if err := json.Unmarshal(notif.Params.(json.RawMessage), &params); err != nil {
		return fmt.Errorf("invalid cancellation params: %w", err)
	}

	// TODO: Implement request cancellation
	return nil
}

// Close ends the session
func (s *Session) Close() error {
	s.cancel()
	return nil
}

// WithImplementation sets the server implementation details
func WithImplementation(impl protocol.Implementation) ServerOption {
	return func(s *Server) {
		s.info = impl
	}
}

// AddTool adds a tool to the server
func (s *Server) AddTool(name string, handler interface{}, description string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tools[name]; exists {
		return fmt.Errorf("tool %s already exists", name)
	}

	s.tools[name] = Tool{
		Handler:     handler,
		Description: description,
		IsAsync:     false,
	}
	return nil
}

// AddAsyncTool adds an asynchronous tool to the server
func (s *Server) AddAsyncTool(name string, handler interface{}, description string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tools[name]; exists {
		return fmt.Errorf("tool %s already exists", name)
	}

	s.tools[name] = Tool{
		Handler:     handler,
		Description: description,
		IsAsync:     true,
	}
	return nil
}

// AddResource adds a resource to the server
func (s *Server) AddResource(pattern string, handler interface{}, description string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.resources[pattern]; exists {
		return fmt.Errorf("resource %s already exists", pattern)
	}

	s.resources[pattern] = Resource{
		Handler:     handler,
		Description: description,
		Pattern:     pattern,
	}
	return nil
}

// AddPrompt adds a prompt to the server
func (s *Server) AddPrompt(name string, handler interface{}, description string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.prompts[name]; exists {
		return fmt.Errorf("prompt %s already exists", name)
	}

	s.prompts[name] = Prompt{
		Handler:     handler,
		Description: description,
	}
	return nil
}
