package protocol

import "net/url"

// Tool represents a tool that can be called by the client
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// CallToolRequestParams represents parameters for calling a tool
type CallToolRequestParams struct {
	RequestParams
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// CallToolResult represents the result of a tool call
type CallToolResult struct {
	Result
	Content []interface{} `json:"content"`
	IsError bool          `json:"isError"`
}

// ListToolsResult represents the result of listing tools
type ListToolsResult struct {
	PaginatedResult
	Tools []Tool `json:"tools"`
}

// Resource represents a resource that can be read by the client
type Resource struct {
	URI         string       `json:"uri"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	MimeType    string       `json:"mimeType,omitempty"`
	Size        *int64       `json:"size,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

// ReadResourceRequestParams represents parameters for reading a resource
type ReadResourceRequestParams struct {
	RequestParams
	URI *url.URL `json:"uri"`
}

// ReadResourceResult represents the result of reading a resource
type ReadResourceResult struct {
	Result
	Contents []interface{} `json:"contents"`
}

// ListResourcesResult represents the result of listing resources
type ListResourcesResult struct {
	PaginatedResult
	Resources []Resource `json:"resources"`
}

// Prompt represents a prompt template
type Prompt struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
}

// PromptArgument represents an argument for a prompt template
type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    *bool  `json:"required,omitempty"`
}

// GetPromptRequestParams represents parameters for getting a prompt
type GetPromptRequestParams struct {
	RequestParams
	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments,omitempty"`
}

// GetPromptResult represents the result of getting a prompt
type GetPromptResult struct {
	Result
	Description string          `json:"description,omitempty"`
	Messages    []PromptMessage `json:"messages"`
}

// ListPromptsResult represents the result of listing prompts
type ListPromptsResult struct {
	PaginatedResult
	Prompts []Prompt `json:"prompts"`
}

// CancelledNotificationParams represents parameters for cancellation notifications
type CancelledNotificationParams struct {
	NotificationParams
	RequestID RequestID `json:"requestId"`
	Reason    string    `json:"reason,omitempty"`
}
