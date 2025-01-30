package server

import "github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/protocol"

// ServerOption configures a Server
type ServerOption func(*Server)

// WithCapabilities sets the server capabilities
func WithCapabilities(caps protocol.ServerCapabilities) ServerOption {
	return func(s *Server) {
		s.capabilities = caps
	}
}

// WithDefaultCapabilities sets up the default server capabilities
func WithDefaultCapabilities() ServerOption {
	return WithCapabilities(protocol.ServerCapabilities{
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
	})
}

// WithExperimentalCapabilities adds experimental capabilities
func WithExperimentalCapabilities(caps map[string]map[string]interface{}) ServerOption {
	return func(s *Server) {
		if s.capabilities.Experimental == nil {
			s.capabilities.Experimental = make(map[string]map[string]interface{})
		}
		for k, v := range caps {
			s.capabilities.Experimental[k] = v
		}
	}
}

// Helper function to create a bool pointer
func boolPtr(b bool) *bool {
	return &b
}
