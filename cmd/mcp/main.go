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

	// Create a new server
	srv := server.NewServer("MCP Example Server")

	// Create a session
	session := server.NewSession(context.Background(), srv)

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

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the transport in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := t.Start(); err != nil {
			errChan <- fmt.Errorf("transport error: %w", err)
		}
	}()

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
