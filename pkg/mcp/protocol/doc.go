// Package protocol defines the core types and messages used in the MCP protocol.
//
// The protocol package includes:
//   - JSON-RPC message types
//   - MCP-specific types (tools, resources, prompts)
//   - Capabilities and initialization types
//   - Content types and message roles
//
// JSON-RPC Messages:
//
//	type JSONRPCRequest struct {
//	    JSONRPC string      `json:"jsonrpc"`
//	    ID      interface{} `json:"id"`
//	    Method  string      `json:"method"`
//	    Params  interface{} `json:"params,omitempty"`
//	}
//
//	type JSONRPCResponse struct {
//	    JSONRPC string      `json:"jsonrpc"`
//	    ID      interface{} `json:"id"`
//	    Result  interface{} `json:"result,omitempty"`
//	    Error   *Error      `json:"error,omitempty"`
//	}
//
// MCP Types:
//
//	// Tool definition
//	type Tool struct {
//	    Name        string `json:"name"`
//	    Description string `json:"description"`
//	}
//
//	// Resource definition
//	type Resource struct {
//	    URI         string `json:"uri"`
//	    Name        string `json:"name"`
//	    Description string `json:"description"`
//	}
//
//	// Prompt definition
//	type Prompt struct {
//	    Name        string `json:"name"`
//	    Description string `json:"description"`
//	}
//
// Content Types:
//
//	// Text content
//	type TextContent struct {
//	    Type string `json:"type"`
//	    Text string `json:"text"`
//	}
//
//	// Image content
//	type ImageContent struct {
//	    Type   string `json:"type"`
//	    Source string `json:"source"`
//	}
//
// Message Roles:
//
//	const (
//	    RoleSystem    = "system"
//	    RoleUser      = "user"
//	    RoleAssistant = "assistant"
//	)
//
// Server Capabilities:
//
//	type ServerCapabilities struct {
//	    SupportsAsync bool `json:"supportsAsync"`
//	}
//
// The protocol package provides all the type definitions needed to implement
// the MCP protocol, ensuring type safety and proper serialization.
package protocol
