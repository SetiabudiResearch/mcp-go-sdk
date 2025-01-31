package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/server"
	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/transport"
)

func main() {
	// Parse command line flags
	transportType := flag.String("transport", "stdio", "Transport type (stdio, sse, or websocket)")
	addr := flag.String("addr", ":8080", "Address to listen on for HTTP transports")
	flag.Parse()

	log.Printf("Starting MCP server with %s transport", *transportType)

	// Create a new server
	srv := server.NewServer("MCP Example Server")
	log.Printf("Created server")

	// Register tools
	srv.AddTool("reverseText", func(text string) (string, error) {
		log.Printf("Received request: %+v", text)

		log.Printf("Received text: %s", text)

		runes := []rune(text)
		log.Printf("Original text: %s", string(runes))
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		log.Printf("Reversed text: %s", string(runes))
		return string(runes), nil
	}, "Reverses the input text")
	log.Printf("Registered tool: reverseText")

	// Create a session
	session := server.NewSession(context.Background(), srv)
	log.Printf("Created session")

	// Create the transport
	var t transport.Transport
	switch *transportType {
	case "stdio":
		t = transport.NewStdioTransport(session)
	case "sse":
		t = transport.NewSSETransport(session, transport.WithAddress(*addr))
	case "websocket":
		t = transport.NewWebSocketTransport(session, transport.WithAddress(*addr))
	default:
		log.Fatalf("Unknown transport type: %s", *transportType)
	}
	log.Printf("Created transport")

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the transport in a goroutine
	errChan := make(chan error, 1)
	go func() {
		log.Printf("Starting transport...")
		if err := t.Start(); err != nil {
			errChan <- fmt.Errorf("transport error: %w", err)
		}
	}()
	log.Printf("Server ready to accept connections")

	// Wait for shutdown signal or error
	select {
	case <-sigChan:
		log.Println("Shutting down...")
		if httpT, ok := t.(transport.HTTPTransport); ok {
			if err := httpT.Stop(context.Background()); err != nil {
				log.Printf("Error stopping transport: %v", err)
			}
		}
	case err := <-errChan:
		log.Printf("Error: %v", err)
	}
}

// Example requests:
// 1. Initialize the server:
// {"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0.0"}}}
//
// 2. Call the reverseText tool (note: use arg0 for the first parameter):
// {"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "reverseText", "arguments": {"text": "Hello, World!"}}}
