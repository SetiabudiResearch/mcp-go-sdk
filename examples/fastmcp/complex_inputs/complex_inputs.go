package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/fastmcp"
	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/protocol"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// SearchParams represents search parameters
type SearchParams struct {
	Query  string   `json:"query"`
	Fields []string `json:"fields"`
	Limit  int      `json:"limit"`
}

func main() {
	// Create a new FastMCP app
	app := fastmcp.New("Complex Inputs Demo")

	// Add a tool that processes a User struct
	app.Tool("createUser", func(user User) (*User, error) {
		// Validate user data
		if user.Name == "" {
			return nil, fmt.Errorf("name is required")
		}
		if user.Email == "" {
			return nil, fmt.Errorf("email is required")
		}

		// Set creation time
		user.CreatedAt = time.Now()
		return &user, nil
	}, "Create a new user")

	// Add a tool that handles arrays and maps
	app.Tool("processArray", func(items []string, metadata map[string]interface{}) interface{} {
		result := make(map[string]interface{})
		result["count"] = len(items)
		result["items"] = items
		result["metadata"] = metadata
		return result
	}, "Process an array of items with metadata")

	// Add a tool that handles complex search parameters
	app.Tool("search", func(params SearchParams) (string, error) {
		if params.Query == "" {
			return "", fmt.Errorf("query is required")
		}
		if params.Limit <= 0 {
			params.Limit = 10 // Default limit
		}

		return fmt.Sprintf(
			"Searching for %q in fields %v with limit %d",
			params.Query,
			params.Fields,
			params.Limit,
		), nil
	}, "Search with complex parameters")

	// Add a tool that handles binary data
	app.Tool("processBinary", func(data []byte) map[string]interface{} {
		return map[string]interface{}{
			"size":     len(data),
			"preview":  fmt.Sprintf("%x", data[:min(len(data), 10)]),
			"is_empty": len(data) == 0,
		}
	}, "Process binary data")

	// Add a resource that returns structured data
	app.Resource("config", func() map[string]interface{} {
		return map[string]interface{}{
			"version":     "1.0.0",
			"debug_mode":  true,
			"max_items":   1000,
			"api_timeout": "30s",
			"features": map[string]bool{
				"feature1": true,
				"feature2": false,
			},
		}
	}, "Get server configuration")

	// Add a prompt that uses structured data
	app.Prompt("dataPrompt", func(data interface{}) []protocol.PromptMessage {
		// Convert data to JSON string for display
		jsonData, _ := json.MarshalIndent(data, "", "  ")
		return []protocol.PromptMessage{
			{
				Role: protocol.RoleAssistant,
				Content: protocol.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Received data:\n%s\n\nHow would you like to process this data?", jsonData),
				},
			},
		}
	}, "Prompt for data processing")

	// Run the server with stdio transport
	log.Println("Starting Complex Inputs Demo...")
	if err := app.RunStdio(); err != nil {
		log.Fatal(err)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
