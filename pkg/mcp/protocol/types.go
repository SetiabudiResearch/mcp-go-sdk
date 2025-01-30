package protocol

import (
	"encoding/json"
	"net/url"
)

// Latest protocol version
const LatestProtocolVersion = "2024-11-05"

// Common types
type ProgressToken interface{} // string or int
type Cursor string
type Role string
type RequestID interface{} // string or int

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Meta represents metadata for requests and notifications
type Meta struct {
	ProgressToken ProgressToken `json:"progressToken,omitempty"`
}

// RequestParams represents base parameters for requests
type RequestParams struct {
	Meta *Meta `json:"_meta,omitempty"`
}

// NotificationParams represents base parameters for notifications
type NotificationParams struct {
	Meta *Meta `json:"_meta,omitempty"`
}

// Request represents a base JSON-RPC request
type Request struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// PaginatedRequest represents a request that supports pagination
type PaginatedRequest struct {
	Request
	Cursor *Cursor `json:"cursor,omitempty"`
}

// Notification represents a base JSON-RPC notification
type Notification struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// Result represents a base JSON-RPC result
type Result struct {
	Meta map[string]interface{} `json:"_meta,omitempty"`
}

// PaginatedResult represents a result that supports pagination
type PaginatedResult struct {
	Result
	NextCursor *Cursor `json:"nextCursor,omitempty"`
}

// JSONRPCRequest represents a JSON-RPC request that expects a response
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      RequestID   `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// JSONRPCNotification represents a JSON-RPC notification
type JSONRPCNotification struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// JSONRPCResponse represents a successful JSON-RPC response
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      RequestID   `json:"id"`
	Result  interface{} `json:"result"`
}

// ErrorData represents error information for JSON-RPC error responses
type ErrorData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// JSONRPCError represents a JSON-RPC error response
type JSONRPCError struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      RequestID `json:"id"`
	Error   ErrorData `json:"error"`
}

// Implementation describes the name and version of an MCP implementation
type Implementation struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Capabilities

type RootsCapability struct {
	ListChanged *bool `json:"listChanged,omitempty"`
}

type SamplingCapability struct{}

type ClientCapabilities struct {
	Experimental map[string]map[string]interface{} `json:"experimental,omitempty"`
	Sampling     *SamplingCapability               `json:"sampling,omitempty"`
	Roots        *RootsCapability                  `json:"roots,omitempty"`
}

// Initialize types

type InitializeRequestParams struct {
	RequestParams
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      Implementation     `json:"clientInfo"`
}

type InitializeRequest struct {
	Method string                  `json:"method"`
	Params InitializeRequestParams `json:"params"`
}

type InitializeResult struct {
	Result
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      Implementation     `json:"serverInfo"`
	Instructions    *string            `json:"instructions,omitempty"`
}

// Content types

type Annotations struct {
	Audience []Role   `json:"audience,omitempty"`
	Priority *float64 `json:"priority,omitempty"`
}

type TextContent struct {
	Type        string       `json:"type"`
	Text        string       `json:"text"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

type ImageContent struct {
	Type        string       `json:"type"`
	Data        string       `json:"data"`
	MimeType    string       `json:"mimeType"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

type ResourceContents struct {
	URI      *url.URL `json:"uri"`
	MimeType *string  `json:"mimeType,omitempty"`
}

type TextResourceContents struct {
	ResourceContents
	Text string `json:"text"`
}

type BlobResourceContents struct {
	ResourceContents
	Blob string `json:"blob"`
}

type EmbeddedResource struct {
	Type        string       `json:"type"`
	Resource    interface{}  // TextResourceContents or BlobResourceContents
	Annotations *Annotations `json:"annotations,omitempty"`
}

// Message types

type SamplingMessage struct {
	Role    Role        `json:"role"`
	Content interface{} // TextContent or ImageContent
}

type PromptMessage struct {
	Role    Role        `json:"role"`
	Content interface{} // TextContent, ImageContent, or EmbeddedResource
}

// Helper functions

func NewTextContent(text string) TextContent {
	return TextContent{
		Type: "text",
		Text: text,
	}
}

func NewImageContent(data, mimeType string) ImageContent {
	return ImageContent{
		Type:     "image",
		Data:     data,
		MimeType: mimeType,
	}
}

func NewEmbeddedResource(resource interface{}, annotations *Annotations) EmbeddedResource {
	return EmbeddedResource{
		Type:        "resource",
		Resource:    resource,
		Annotations: annotations,
	}
}
