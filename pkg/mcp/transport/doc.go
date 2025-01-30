// Package transport provides different transport layer implementations for the MCP protocol.
//
// The transport package includes:
//   - Stdio transport for command-line applications
//   - WebSocket transport for web applications
//   - Server-Sent Events (SSE) transport for web browsers
//
// Each transport implements the Transport interface:
//
//	type Transport interface {
//	    Start() error
//	    Stop() error
//	}
//
// Stdio Transport:
//
//	// Create a stdio transport
//	t := transport.NewStdioTransport(session)
//
//	// Start the transport
//	if err := t.Start(); err != nil {
//	    log.Fatal(err)
//	}
//
// WebSocket Transport:
//
//	// Create a WebSocket transport with options
//	t := transport.NewWebSocketTransport(session,
//	    transport.WithAddress(":8080"),
//	    transport.WithPath("/ws"),
//	)
//
//	// Start the transport
//	if err := t.Start(); err != nil {
//	    log.Fatal(err)
//	}
//
// SSE Transport:
//
//	// Create an SSE transport with options
//	t := transport.NewSSETransport(session,
//	    transport.WithAddress(":8080"),
//	    transport.WithPath("/events"),
//	)
//
//	// Start the transport
//	if err := t.Start(); err != nil {
//	    log.Fatal(err)
//	}
//
// Transport Options:
//
// Each transport type supports configuration through options:
//
//	type TransportOption func(t Transport)
//
//	// Common options
//	WithAddress(addr string)      // Set the listening address
//	WithPath(path string)         // Set the endpoint path
//	WithTLSConfig(config *tls.Config) // Configure TLS
//
// The transport package handles all the low-level communication details,
// allowing the server to focus on business logic.
package transport
