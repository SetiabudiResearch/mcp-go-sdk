package mcp

import (
	"testing"
)

func TestNewServer(t *testing.T) {
	server := NewServer("test")
	if server.name != "test" {
		t.Errorf("expected server name to be 'test', got %s", server.name)
	}
}

func TestAddTool(t *testing.T) {
	server := NewServer("test")

	// Test valid tool
	err := server.AddTool("add", func(a, b int) int { return a + b }, "Add two numbers")
	if err != nil {
		t.Errorf("unexpected error adding valid tool: %v", err)
	}

	// Test nil handler
	err = server.AddTool("nil", nil, "Invalid tool")
	if err == nil {
		t.Error("expected error for nil handler, got nil")
	}

	// Test non-function handler
	err = server.AddTool("invalid", "not a function", "Invalid tool")
	if err == nil {
		t.Error("expected error for non-function handler, got nil")
	}
}

func TestAddAsyncTool(t *testing.T) {
	server := NewServer("test")

	// Test valid async tool
	err := server.AddAsyncTool("fetch", func(url string) (string, error) { return "", nil }, "Fetch URL")
	if err != nil {
		t.Errorf("unexpected error adding valid async tool: %v", err)
	}

	// Test invalid return type
	err = server.AddAsyncTool("invalid", func(x int) (int, string) { return 0, "" }, "Invalid tool")
	if err == nil {
		t.Error("expected error for invalid return type, got nil")
	}
}

func TestAddResource(t *testing.T) {
	server := NewServer("test")

	// Test valid resource
	err := server.AddResource("data://{id}", func(id string) string { return id }, "Get data")
	if err != nil {
		t.Errorf("unexpected error adding valid resource: %v", err)
	}

	// Test nil handler
	err = server.AddResource("nil://test", nil, "Invalid resource")
	if err == nil {
		t.Error("expected error for nil handler, got nil")
	}
}

func TestAddPrompt(t *testing.T) {
	server := NewServer("test")

	// Test valid prompt
	err := server.AddPrompt("greet", func(name string) string { return "Hello, " + name }, "Greeting prompt")
	if err != nil {
		t.Errorf("unexpected error adding valid prompt: %v", err)
	}

	// Test nil handler
	err = server.AddPrompt("nil", nil, "Invalid prompt")
	if err == nil {
		t.Error("expected error for nil handler, got nil")
	}
}

func TestWithDependencies(t *testing.T) {
	deps := []string{"github.com/example/pkg1", "github.com/example/pkg2"}
	server := NewServer("test", WithDependencies(deps))

	if len(server.dependencies) != len(deps) {
		t.Errorf("expected %d dependencies, got %d", len(deps), len(server.dependencies))
	}

	for i, dep := range deps {
		if server.dependencies[i] != dep {
			t.Errorf("expected dependency %s, got %s", dep, server.dependencies[i])
		}
	}
}
