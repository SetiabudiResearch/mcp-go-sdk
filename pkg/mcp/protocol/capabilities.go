package protocol

// ServerCapabilities defines the capabilities of an MCP server
type ServerCapabilities struct {
	Experimental map[string]map[string]interface{} `json:"experimental,omitempty"`
	Logging      *LoggingCapability                `json:"logging,omitempty"`
	Prompts      *PromptsCapability                `json:"prompts,omitempty"`
	Resources    *ResourcesCapability              `json:"resources,omitempty"`
	Tools        *ToolsCapability                  `json:"tools,omitempty"`
}

// PromptsCapability defines prompt-related capabilities
type PromptsCapability struct {
	ListChanged *bool `json:"listChanged,omitempty"`
}

// ResourcesCapability defines resource-related capabilities
type ResourcesCapability struct {
	Subscribe   *bool `json:"subscribe,omitempty"`
	ListChanged *bool `json:"listChanged,omitempty"`
}

// ToolsCapability defines tool-related capabilities
type ToolsCapability struct {
	ListChanged *bool `json:"listChanged,omitempty"`
}

// LoggingCapability defines logging-related capabilities
type LoggingCapability struct{}
