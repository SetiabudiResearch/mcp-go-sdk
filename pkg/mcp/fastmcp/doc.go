// Package fastmcp provides a simplified interface for creating MCP servers.
//
// FastMCP reduces boilerplate code and provides a fluent interface for configuring
// servers. It wraps the core MCP server implementation with a more ergonomic API
// while maintaining full protocol compatibility.
//
// Basic usage:
//
//	app := fastmcp.New("My App")
//
//	// Add tools
//	app.Tool("greet", func(name string) string {
//	    return "Hello, " + name + "!"
//	}, "Greet a person")
//
//	// Add async tools
//	app.AsyncTool("longProcess", func(params string) error {
//	    // Long-running operation
//	    return nil
//	}, "Process data asynchronously")
//
//	// Add resources
//	app.Resource("files/{path}", func(path string) ([]byte, error) {
//	    return ioutil.ReadFile(path)
//	}, "Access files")
//
//	// Add prompts
//	app.Prompt("confirm", func(action string) []protocol.PromptMessage {
//	    return []protocol.PromptMessage{
//	        {
//	            Role: protocol.RoleAssistant,
//	            Content: protocol.TextContent{
//	                Type: "text",
//	                Text: fmt.Sprintf("Are you sure you want to %s?", action),
//	            },
//	        },
//	    }
//	}, "Confirmation prompt")
//
//	// Run with stdio transport
//	if err := app.RunStdio(); err != nil {
//	    log.Fatal(err)
//	}
//
// Server Capabilities:
//
// FastMCP automatically configures default server capabilities:
//   - Tool list change notifications
//   - Resource subscriptions and list change notifications
//   - Prompt list change notifications
//   - Logging support
//
// You can customize capabilities when creating a new FastMCP instance:
//
//	app := fastmcp.New("My App", server.WithCapabilities(protocol.ServerCapabilities{
//	    Tools: &protocol.ToolsCapability{
//	        ListChanged: boolPtr(true),
//	    },
//	    Resources: &protocol.ResourcesCapability{
//	        Subscribe:   boolPtr(true),
//	        ListChanged: boolPtr(true),
//	    },
//	}))
//
// Transport Options:
//
//	// Run with stdio (for CLI apps)
//	app.RunStdio()
//
//	// Run with WebSocket (for web apps)
//	app.RunWebSocket(":8080")
//
//	// Run with SSE (for web browsers)
//	app.RunSSE(":8080")
//
// The FastMCP API is designed to be chainable:
//
//	fastmcp.New("My App").
//	    Tool("tool1", handler1, "desc1").
//	    Tool("tool2", handler2, "desc2").
//	    Resource("res1", handler3, "desc3").
//	    RunStdio()
//
// Handler Function Types:
//
// FastMCP supports various handler function signatures:
//
//	// Simple tools
//	func(string) string
//	func(int) (int, error)
//
//	// Complex tools with structs
//	func(User) (*User, error)
//	func(SearchParams) ([]Result, error)
//
//	// Resources with parameters
//	func(path string) ([]byte, error)
//	func(id string, version int) (interface{}, error)
//
//	// Prompts with multiple messages
//	func(context string) []protocol.PromptMessage
//	func(action string) protocol.TextContent
//
// For more control, you can access the underlying server:
//
//	app := fastmcp.New("My App")
//	server := app.Server()
//
// FastMCP is ideal for building MCP applications quickly while maintaining
// access to the full power of the protocol when needed.
package fastmcp
