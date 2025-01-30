// Package server provides the core server implementation for the MCP protocol.
//
// The server package handles:
//   - Server lifecycle management
//   - Session management
//   - Tool registration and execution
//   - Resource pattern matching and access
//   - Prompt template rendering
//   - Request and notification handling
//
// Server Creation and Configuration:
//
//	// Create a new server
//	srv := server.NewServer("My Server")
//
//	// Configure server capabilities
//	srv.WithCapabilities(protocol.ServerCapabilities{
//	    SupportsAsync: true,
//	})
//
// Tool Registration:
//
//	// Add a synchronous tool
//	srv.AddTool("myTool", func(arg1 string, arg2 int) (string, error) {
//	    return fmt.Sprintf("Processed %s with %d", arg1, arg2), nil
//	}, "Tool description")
//
//	// Add an asynchronous tool
//	srv.AddAsyncTool("longRunningTool", func(params string) error {
//	    // Long-running operation
//	    return nil
//	}, "Async tool description")
//
// Resource Registration:
//
//	// Add a resource with pattern matching
//	srv.AddResource("files/{path}", func(path string) ([]byte, error) {
//	    return ioutil.ReadFile(path)
//	}, "Access files")
//
// Prompt Registration:
//
//	// Add a prompt template
//	srv.AddPrompt("confirm", func(action string) string {
//	    return fmt.Sprintf("Are you sure you want to %s?", action)
//	}, "Confirmation prompt")
//
// Session Management:
//
//	// Create a new session
//	session := server.NewSession(context.Background(), srv)
//
//	// Handle requests through the session
//	response, err := session.HandleRequest(request)
//
// The server package uses reflection to dynamically invoke handlers and convert
// parameters, making it easy to register any Go function as a tool, resource,
// or prompt handler.
package server
