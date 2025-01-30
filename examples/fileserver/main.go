package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/protocol"
	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/server"
	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/transport"
)

// FileServer implements a simple file server with MCP
type FileServer struct {
	rootDir string
	srv     *server.Server
}

// NewFileServer creates a new file server
func NewFileServer(rootDir string) (*FileServer, error) {
	// Create the server
	srv := server.NewServer("File Server")
	fs := &FileServer{
		rootDir: rootDir,
		srv:     srv,
	}

	// Add tools
	if err := fs.registerTools(); err != nil {
		return nil, err
	}

	// Add resources
	if err := fs.registerResources(); err != nil {
		return nil, err
	}

	// Add prompts
	if err := fs.registerPrompts(); err != nil {
		return nil, err
	}

	return fs, nil
}

// registerTools registers file operation tools
func (fs *FileServer) registerTools() error {
	// Write file tool
	if err := fs.srv.AddTool("writeFile", fs.writeFile, "Write content to a file"); err != nil {
		return err
	}

	// Delete file tool
	if err := fs.srv.AddTool("deleteFile", fs.deleteFile, "Delete a file"); err != nil {
		return err
	}

	// Create directory tool
	if err := fs.srv.AddTool("createDir", fs.createDir, "Create a directory"); err != nil {
		return err
	}

	return nil
}

// registerResources registers file system resources
func (fs *FileServer) registerResources() error {
	// File content resource
	if err := fs.srv.AddResource("file://{path}", fs.readFile, "Read file contents"); err != nil {
		return err
	}

	// Directory listing resource
	if err := fs.srv.AddResource("dir://{path}", fs.listDir, "List directory contents"); err != nil {
		return err
	}

	return nil
}

// registerPrompts registers file operation prompts
func (fs *FileServer) registerPrompts() error {
	// File operation prompt
	if err := fs.srv.AddPrompt("fileOp", fs.fileOpPrompt, "File operation prompt"); err != nil {
		return err
	}

	return nil
}

// Tool implementations

func (fs *FileServer) writeFile(path string, content string) error {
	fullPath := filepath.Join(fs.rootDir, path)
	return ioutil.WriteFile(fullPath, []byte(content), 0644)
}

func (fs *FileServer) deleteFile(path string) error {
	fullPath := filepath.Join(fs.rootDir, path)
	return os.Remove(fullPath)
}

func (fs *FileServer) createDir(path string) error {
	fullPath := filepath.Join(fs.rootDir, path)
	return os.MkdirAll(fullPath, 0755)
}

// Resource implementations

func (fs *FileServer) readFile(path string) (string, error) {
	fullPath := filepath.Join(fs.rootDir, path)
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (fs *FileServer) listDir(path string) ([]string, error) {
	fullPath := filepath.Join(fs.rootDir, path)
	entries, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name())
	}
	return files, nil
}

// Prompt implementations

func (fs *FileServer) fileOpPrompt(path string, operation string) []protocol.PromptMessage {
	return []protocol.PromptMessage{
		{
			Role: protocol.RoleAssistant,
			Content: protocol.TextContent{
				Type: "text",
				Text: fmt.Sprintf("You are about to %s the file at %s. Are you sure?", operation, path),
			},
		},
		{
			Role: protocol.RoleUser,
			Content: protocol.TextContent{
				Type: "text",
				Text: "Please confirm with 'yes' or 'no'.",
			},
		},
	}
}

func main() {
	// Parse command line flags
	rootDir := flag.String("root", ".", "Root directory to serve")
	transportType := flag.String("transport", "stdio", "Transport type (stdio, sse, or websocket)")
	addr := flag.String("addr", ":8080", "Address to listen on for HTTP transports")
	flag.Parse()

	// Create the file server
	fs, err := NewFileServer(*rootDir)
	if err != nil {
		log.Fatalf("Failed to create file server: %v", err)
	}

	// Create a session
	session := server.NewSession(context.Background(), fs.srv)

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

	// Start the transport
	log.Printf("Starting file server on %s transport", *transportType)
	if err := t.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
