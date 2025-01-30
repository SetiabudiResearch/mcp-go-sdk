package fastmcp

import (
	"context"
	"fmt"
	"log"

	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/protocol"
	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/server"
	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/transport"
)

// FastMCP provides a simplified interface for creating MCP servers
type FastMCP struct {
	name    string
	server  *server.Server
	options []server.ServerOption
}

// New creates a new FastMCP instance with default capabilities
func New(name string, options ...server.ServerOption) *FastMCP {
	// Default capabilities
	defaultOptions := []server.ServerOption{
		server.WithCapabilities(protocol.ServerCapabilities{
			Tools: &protocol.ToolsCapability{
				ListChanged: boolPtr(true),
			},
			Resources: &protocol.ResourcesCapability{
				Subscribe:   boolPtr(true),
				ListChanged: boolPtr(true),
			},
			Prompts: &protocol.PromptsCapability{
				ListChanged: boolPtr(true),
			},
			Logging: &protocol.LoggingCapability{},
		}),
	}

	// Append user options after defaults
	options = append(defaultOptions, options...)

	return &FastMCP{
		name:    name,
		options: options,
	}
}

// Tool registers a synchronous tool with the server
func (f *FastMCP) Tool(name string, handler interface{}, description string) *FastMCP {
	if f.server == nil {
		f.server = server.NewServer(f.name, f.options...)
	}
	if err := f.server.AddTool(name, handler, description); err != nil {
		log.Printf("Warning: Failed to add tool %s: %v", name, err)
	}
	return f
}

// AsyncTool registers an asynchronous tool with the server
func (f *FastMCP) AsyncTool(name string, handler interface{}, description string) *FastMCP {
	if f.server == nil {
		f.server = server.NewServer(f.name, f.options...)
	}
	if err := f.server.AddAsyncTool(name, handler, description); err != nil {
		log.Printf("Warning: Failed to add async tool %s: %v", name, err)
	}
	return f
}

// Resource registers a resource with the server
func (f *FastMCP) Resource(pattern string, handler interface{}, description string) *FastMCP {
	if f.server == nil {
		f.server = server.NewServer(f.name, f.options...)
	}
	if err := f.server.AddResource(pattern, handler, description); err != nil {
		log.Printf("Warning: Failed to add resource %s: %v", pattern, err)
	}
	return f
}

// Prompt registers a prompt with the server
func (f *FastMCP) Prompt(name string, handler interface{}, description string) *FastMCP {
	if f.server == nil {
		f.server = server.NewServer(f.name, f.options...)
	}
	if err := f.server.AddPrompt(name, handler, description); err != nil {
		log.Printf("Warning: Failed to add prompt %s: %v", name, err)
	}
	return f
}

// RunStdio starts the server with stdio transport
func (f *FastMCP) RunStdio() error {
	if f.server == nil {
		return fmt.Errorf("no server configured")
	}
	session := server.NewSession(context.Background(), f.server)
	t := transport.NewStdioTransport(session)
	return t.Start()
}

// RunWebSocket starts the server with WebSocket transport
func (f *FastMCP) RunWebSocket(addr string) error {
	if f.server == nil {
		return fmt.Errorf("no server configured")
	}
	session := server.NewSession(context.Background(), f.server)
	t := transport.NewWebSocketTransport(session, transport.WithAddress(addr))
	return t.Start()
}

// RunSSE starts the server with SSE transport
func (f *FastMCP) RunSSE(addr string) error {
	if f.server == nil {
		return fmt.Errorf("no server configured")
	}
	session := server.NewSession(context.Background(), f.server)
	t := transport.NewSSETransport(session, transport.WithAddress(addr))
	return t.Start()
}

// Server returns the underlying server instance
func (f *FastMCP) Server() *server.Server {
	return f.server
}

// Helper function to create a bool pointer
func boolPtr(b bool) *bool {
	return &b
}
