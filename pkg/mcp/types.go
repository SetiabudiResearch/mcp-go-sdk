package mcp

import (
	"context"
	"image"
)

// Server represents an MCP server instance
type Server struct {
	name         string
	tools        map[string]Tool
	resources    map[string]Resource
	prompts      map[string]Prompt
	dependencies []string
}

// ServerOption is a function that configures a Server
type ServerOption func(*Server)

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

// Message represents a message in an LLM conversation
type Message struct {
	Role    string
	Content string
}

// Image represents an image that can be sent to or received from the LLM
type Image struct {
	Data   image.Image
	Format string
}

// Context provides access to MCP capabilities during tool and resource execution
type Context struct {
	ctx context.Context
}

// NewUserMessage creates a new message with the user role
func NewUserMessage(content string) Message {
	return Message{
		Role:    "user",
		Content: content,
	}
}

// NewAssistantMessage creates a new message with the assistant role
func NewAssistantMessage(content string) Message {
	return Message{
		Role:    "assistant",
		Content: content,
	}
}

// NewImage creates a new Image instance
func NewImage(img image.Image, format string) *Image {
	return &Image{
		Data:   img,
		Format: format,
	}
}

// WithDependencies configures the server with additional dependencies
func WithDependencies(deps []string) ServerOption {
	return func(s *Server) {
		s.dependencies = deps
	}
}
