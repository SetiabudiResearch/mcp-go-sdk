# MCP Go SDK Examples

This directory contains example applications demonstrating various features of the MCP Go SDK.

## Echo Server (`/echo`)

A simple echo server that demonstrates basic MCP features:
- Tool registration and execution
- Resource pattern matching
- Prompt templates
- Basic string manipulation

```bash
cd echo
go run main.go
```

## Complex Inputs (`/complex`)

Demonstrates handling of complex data types:
- Structured data (structs, arrays, maps)
- Binary data
- Custom types
- Error handling
- JSON serialization

```bash
cd complex
go run main.go
```

## Desktop Integration (`/desktop`)

Shows how to integrate with desktop features:
- File system operations
- System information
- Environment variables
- Cross-platform compatibility
- User confirmation prompts

```bash
cd desktop
go run main.go
```

## Common Features

All examples demonstrate:
1. Using the FastMCP API for simplified server creation
2. Proper error handling
3. Resource pattern matching
4. Prompt template rendering
5. Tool registration and execution

## Running the Examples

Each example can be run using `go run main.go` in its respective directory. The examples use the stdio transport by default, making them easy to test from the command line.

## Testing with MCP Clients

You can test these examples using any MCP client that supports the stdio transport. For example:

```bash
# Using the MCP CLI tool
mcp connect go run main.go

# Using Claude Desktop
mcp install main.go
```

## Notes

- The examples use the FastMCP API for simplicity
- Error handling follows Go best practices
- Each example is self-contained
- The examples demonstrate different aspects of the MCP protocol 