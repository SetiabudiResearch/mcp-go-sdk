package main

import (
	"fmt"
	"log"
	"math"

	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp"
)

func main() {
	// Create a new MCP server
	server := mcp.NewServer("Calculator")

	// Add basic arithmetic tools
	server.AddTool("add", func(a, b float64) float64 {
		return a + b
	}, "Add two numbers")

	server.AddTool("subtract", func(a, b float64) float64 {
		return a - b
	}, "Subtract b from a")

	server.AddTool("multiply", func(a, b float64) float64 {
		return a * b
	}, "Multiply two numbers")

	server.AddTool("divide", func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	}, "Divide a by b")

	// Add advanced math tools
	server.AddTool("power", func(base, exp float64) float64 {
		return math.Pow(base, exp)
	}, "Raise base to the power of exp")

	server.AddTool("sqrt", func(x float64) (float64, error) {
		if x < 0 {
			return 0, fmt.Errorf("cannot calculate square root of negative number")
		}
		return math.Sqrt(x), nil
	}, "Calculate the square root of x")

	// Add a history resource
	var history []string
	server.AddResource("history", func() []string {
		return history
	}, "Get calculation history")

	// Add a calculation prompt
	server.AddPrompt("calculate", func(expression string) string {
		return fmt.Sprintf("Please calculate: %s", expression)
	}, "Format a calculation request")

	// Start the server
	log.Println("Starting calculator server...")
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
