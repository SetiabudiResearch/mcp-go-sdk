# MCP Go SDK

This SDK provides a Go implementation of the Model Context Protocol (MCP), enabling bidirectional communication between clients and servers for tool execution, resource access, and prompt handling.

## Features

- **Transport Layer**
  - Multiple transport options (stdio, SSE, WebSocket)
  - Bidirectional communication
  - Configurable endpoints and settings

- **Server Implementation**
  - Tool registration and execution
  - Resource pattern matching and access
  - Prompt template rendering
  - Session management
  - Reflection-based handler invocation

- **Core Protocol Types**
  - JSON-RPC message handling
  - MCP-specific types (tools, resources, prompts)
  - Capabilities and initialization

## Installation

```bash
go get github.com/SetiabudiResearch/mcp-go-sdk
```

## Quick Start

Here's a simple example of creating an MCP server:

```go
package main

import (
    "context"
    "log"
    "github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/server"
    "github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/transport"
)

func main() {
    // Create a new server
    srv := server.NewServer("Example Server")

    // Add a tool
    srv.AddTool("greet", func(name string) string {
        return "Hello, " + name + "!"
    }, "Greet a person")

    // Create a session
    session := server.NewSession(context.Background(), srv)

    // Create and start a transport
    t := transport.NewStdioTransport(session)
    if err := t.Start(); err != nil {
        log.Fatal(err)
    }
}
```

## Usage Guide

### Creating a Server

```go
// Create a new server with a name
srv := server.NewServer("My Server")

// Optionally configure server capabilities
srv.WithCapabilities(protocol.ServerCapabilities{
    SupportsAsync: true,
})
```

### Adding Tools

```go
// Add a synchronous tool
srv.AddTool("myTool", func(arg1 string, arg2 int) (string, error) {
    return fmt.Sprintf("Processed %s with %d", arg1, arg2), nil
}, "Tool description")

// Add an asynchronous tool
srv.AddAsyncTool("longRunningTool", func(params string) error {
    // Long-running operation
    return nil
}, "Async tool description")
```

### Adding Resources

```go
// Add a resource with pattern matching
srv.AddResource("files/{path}", func(path string) ([]byte, error) {
    return ioutil.ReadFile(path)
}, "Access files")

// Resource patterns support multiple parameters
srv.AddResource("api/{version}/{endpoint}", func(version, endpoint string) (interface{}, error) {
    return callAPI(version, endpoint)
}, "API access")
```

### Adding Prompts

```go
// Add a simple text prompt
srv.AddPrompt("confirm", func(action string) string {
    return fmt.Sprintf("Are you sure you want to %s?", action)
}, "Confirmation prompt")

// Add a multi-message prompt
srv.AddPrompt("chat", func(context string) []protocol.PromptMessage {
    return []protocol.PromptMessage{
        {
            Role: protocol.RoleAssistant,
            Content: protocol.TextContent{
                Type: "text",
                Text: "How can I help you with " + context + "?",
            },
        },
    }
}, "Chat prompt")
```

### Transport Configuration

```go
// Create a session
session := server.NewSession(context.Background(), srv)

// Stdio transport (for CLI applications)
t := transport.NewStdioTransport(session)

// WebSocket transport (for web applications)
t := transport.NewWebSocketTransport(session, transport.WithAddress(":8080"))

// Server-Sent Events transport (for web browsers)
t := transport.NewSSETransport(session, transport.WithAddress(":8080"))

// Start the transport
if err := t.Start(); err != nil {
    log.Fatal(err)
}
```

## Example Applications

See the [examples](./examples) directory for complete example applications:

- [File Server](./examples/fileserver): A complete file server implementation using MCP
- [Calculator](./examples/calculator): A simple calculator service
- [Chat Bot](./examples/chatbot): An example chat bot using prompts

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 
