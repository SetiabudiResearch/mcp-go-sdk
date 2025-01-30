// Package mcp provides a Go implementation of the Machine Control Protocol (MCP).
//
// The MCP protocol enables bidirectional communication between clients and servers
// for tool execution, resource access, and prompt handling. This implementation
// provides a complete framework for building MCP servers and clients.
//
// The package is organized into several sub-packages:
//
//   - protocol: Core protocol types and message definitions
//   - server: Server implementation with session management
//   - transport: Transport layer implementations (stdio, SSE, WebSocket)
//
// Basic usage example:
//
//	func main() {
//	    // Create a new server
//	    srv := server.NewServer("Example Server")
//
//	    // Add a tool
//	    srv.AddTool("greet", func(name string) string {
//	        return "Hello, " + name + "!"
//	    }, "Greet a person")
//
//	    // Create a session
//	    session := server.NewSession(context.Background(), srv)
//
//	    // Create and start a transport
//	    t := transport.NewStdioTransport(session)
//	    if err := t.Start(); err != nil {
//	        log.Fatal(err)
//	    }
//	}
package mcp
