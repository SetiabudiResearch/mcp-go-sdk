package server

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/protocol"
)

// handleListTools processes tools/list requests
func (s *Session) handleListTools(req *protocol.JSONRPCRequest) (*protocol.JSONRPCResponse, error) {
	s.server.mu.RLock()
	tools := make([]protocol.Tool, 0, len(s.server.tools))
	for name, tool := range s.server.tools {
		tools = append(tools, protocol.Tool{
			Name:        name,
			Description: tool.Description,
		})
	}
	s.server.mu.RUnlock()

	result := protocol.ListToolsResult{
		Tools: tools,
	}

	return &protocol.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}, nil
}

// handleCallTool processes tools/call requests
func (s *Session) handleCallTool(req *protocol.JSONRPCRequest) (*protocol.JSONRPCResponse, error) {
	var params protocol.CallToolRequestParams
	if err := json.Unmarshal(req.Params.(json.RawMessage), &params); err != nil {
		return nil, fmt.Errorf("invalid tool call params: %w", err)
	}
	log.Printf("Received tool call request: %+v", params)

	s.server.mu.RLock()
	tool, exists := s.server.tools[params.Name]
	s.server.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("tool not found: %s", params.Name)
	}

	// Convert arguments to reflect.Values
	handlerType := reflect.TypeOf(tool.Handler)
	args := make([]reflect.Value, handlerType.NumIn())
	for i := 0; i < handlerType.NumIn(); i++ {
		paramType := handlerType.In(i)
		paramValue := reflect.New(paramType).Interface()

		// Get argument value from params
		argName := fmt.Sprintf("arg%d", i)
		log.Printf("argName: %s", argName)
		if params.Arguments != nil {
			if argValue, ok := params.Arguments[argName]; ok {
				// Directly assign the argument value
				paramValue = argValue
			} else {
				return nil, fmt.Errorf("missing argument: %s", argName)
			}
		} else {
			return nil, fmt.Errorf("arguments map is nil")
		}
		args[i] = reflect.ValueOf(paramValue)
	}

	// Call the handler
	log.Printf("arguments: %+v", args)
	results := reflect.ValueOf(tool.Handler).Call(args)

	// Process results
	var content []interface{}
	var isError bool

	if len(results) == 2 { // Function returns (value, error)
		if !results[1].IsNil() { // Error occurred
			err := results[1].Interface().(error)
			content = []interface{}{protocol.NewTextContent(err.Error())}
			isError = true
		} else {
			content = []interface{}{results[0].String()}
		}
	} else { // Function returns single value
		content = []interface{}{results[0].String()}
	}

	result := protocol.CallToolResult{
		Content: content,
		IsError: isError,
	}

	return &protocol.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}, nil
}

// handleListResources processes resources/list requests
func (s *Session) handleListResources(req *protocol.JSONRPCRequest) (*protocol.JSONRPCResponse, error) {
	s.server.mu.RLock()
	resources := make([]protocol.Resource, 0, len(s.server.resources))
	for _, resource := range s.server.resources {
		resources = append(resources, protocol.Resource{
			URI:         resource.Pattern,
			Name:        resource.Pattern,
			Description: resource.Description,
		})
	}
	s.server.mu.RUnlock()

	result := protocol.ListResourcesResult{
		Resources: resources,
	}

	return &protocol.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}, nil
}

// handleReadResource processes resources/read requests
func (s *Session) handleReadResource(req *protocol.JSONRPCRequest) (*protocol.JSONRPCResponse, error) {
	var params protocol.ReadResourceRequestParams
	if err := json.Unmarshal(req.Params.(json.RawMessage), &params); err != nil {
		return nil, fmt.Errorf("invalid resource read params: %w", err)
	}

	// Find matching resource and extract parameters
	resource, resourceParams, err := s.server.matchResource(params.URI.String())
	if err != nil {
		return nil, fmt.Errorf("failed to match resource: %w", err)
	}

	// Read the resource
	contents, err := s.server.readResource(resource, resourceParams)
	if err != nil {
		return nil, fmt.Errorf("failed to read resource: %w", err)
	}

	result := protocol.ReadResourceResult{
		Contents: contents,
	}

	return &protocol.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}, nil
}

// handleListPrompts processes prompts/list requests
func (s *Session) handleListPrompts(req *protocol.JSONRPCRequest) (*protocol.JSONRPCResponse, error) {
	s.server.mu.RLock()
	prompts := make([]protocol.Prompt, 0, len(s.server.prompts))
	for name, prompt := range s.server.prompts {
		prompts = append(prompts, protocol.Prompt{
			Name:        name,
			Description: prompt.Description,
		})
	}
	s.server.mu.RUnlock()

	result := protocol.ListPromptsResult{
		Prompts: prompts,
	}

	return &protocol.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}, nil
}

// handleGetPrompt processes prompts/get requests
func (s *Session) handleGetPrompt(req *protocol.JSONRPCRequest) (*protocol.JSONRPCResponse, error) {
	var params protocol.GetPromptRequestParams
	if err := json.Unmarshal(req.Params.(json.RawMessage), &params); err != nil {
		return nil, fmt.Errorf("invalid prompt get params: %w", err)
	}

	s.server.mu.RLock()
	prompt, exists := s.server.prompts[params.Name]
	s.server.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("prompt not found: %s", params.Name)
	}

	// Render the prompt
	messages, err := s.server.renderPrompt(prompt, params.Arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to render prompt: %w", err)
	}

	result := protocol.GetPromptResult{
		Description: prompt.Description,
		Messages:    messages,
	}

	return &protocol.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}, nil
}
