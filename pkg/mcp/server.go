package mcp

import (
	"fmt"
	"reflect"
)

// NewServer creates a new MCP server instance
func NewServer(name string, opts ...ServerOption) *Server {
	s := &Server{
		name:      name,
		tools:     make(map[string]Tool),
		resources: make(map[string]Resource),
		prompts:   make(map[string]Prompt),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// AddTool adds a synchronous tool to the server
func (s *Server) AddTool(name string, handler interface{}, description string) error {
	if err := validateHandler(handler); err != nil {
		return fmt.Errorf("invalid tool handler: %w", err)
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
	if err := validateHandler(handler); err != nil {
		return fmt.Errorf("invalid async tool handler: %w", err)
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
	if err := validateHandler(handler); err != nil {
		return fmt.Errorf("invalid resource handler: %w", err)
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
	if err := validateHandler(handler); err != nil {
		return fmt.Errorf("invalid prompt handler: %w", err)
	}

	s.prompts[name] = Prompt{
		Handler:     handler,
		Description: description,
	}
	return nil
}

// Start starts the MCP server
func (s *Server) Start() error {
	// TODO: Implement server startup logic
	return nil
}

// validateHandler checks if a handler function has a valid signature
func validateHandler(handler interface{}) error {
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	t := reflect.TypeOf(handler)
	if t.Kind() != reflect.Func {
		return fmt.Errorf("handler must be a function")
	}

	// Check return types
	if t.NumOut() < 1 {
		return fmt.Errorf("handler must return at least one value")
	}

	// For async handlers, check error return
	if t.NumOut() == 2 && !t.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return fmt.Errorf("second return value must be error")
	}

	return nil
}

// Context methods

// Info logs an informational message
func (c *Context) Info(msg string) {
	// TODO: Implement logging
}

// ReportProgress reports progress of a long-running operation
func (c *Context) ReportProgress(progress float64) {
	// TODO: Implement progress reporting
}

// ReadResource reads data from a resource
func (c *Context) ReadResource(uri string) ([]byte, string, error) {
	// TODO: Implement resource reading
	return nil, "", nil
}
