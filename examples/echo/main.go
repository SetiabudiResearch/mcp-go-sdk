package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/fastmcp"
	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/protocol"
)

func main() {
	// Create a new FastMCP app
	app := fastmcp.New("Echo Server")

	// Add an echo tool
	app.Tool("echo", func(message string) string {
		return message
	}, "Echo back the input message")

	// Add a reverse tool
	app.Tool("reverse", func(message string) string {
		runes := []rune(message)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	}, "Reverse the input message")

	// Add an uppercase tool
	app.Tool("uppercase", func(message string) string {
		return strings.ToUpper(message)
	}, "Convert message to uppercase")

	// Add a counter resource
	var counter int
	app.Resource("counter", func() string {
		counter++
		return fmt.Sprintf("Counter: %d", counter)
	}, "A simple counter that increments on each access")

	// Add a confirmation prompt
	app.Prompt("confirm", func(action string) []protocol.PromptMessage {
		return []protocol.PromptMessage{
			{
				Role: protocol.RoleAssistant,
				Content: protocol.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Are you sure you want to %s?", action),
				},
			},
			{
				Role: protocol.RoleUser,
				Content: protocol.TextContent{
					Type: "text",
					Text: "Please respond with 'yes' or 'no'.",
				},
			},
		}
	}, "A confirmation prompt")

	// Run the server with stdio transport
	log.Println("Starting Echo Server...")
	if err := app.RunStdio(); err != nil {
		log.Fatal(err)
	}
}
