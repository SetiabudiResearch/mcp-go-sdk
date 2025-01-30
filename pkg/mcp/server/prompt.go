package server

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"

	"github.com/SetiabudiResearch/mcp-go-sdk/pkg/mcp/protocol"
)

// promptTemplate represents a parsed prompt template
type promptTemplate struct {
	handler     interface{}
	description string
	arguments   []protocol.PromptArgument
	template    *template.Template
}

// parsePromptTemplate parses a prompt handler into a template
func parsePromptTemplate(name string, handler interface{}, description string) (*promptTemplate, error) {
	// Validate handler
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		return nil, fmt.Errorf("handler must be a function")
	}

	// Extract argument information
	var arguments []protocol.PromptArgument
	for i := 0; i < handlerType.NumIn(); i++ {
		paramType := handlerType.In(i)
		required := true

		arguments = append(arguments, protocol.PromptArgument{
			Name:        fmt.Sprintf("arg%d", i),
			Description: fmt.Sprintf("Argument of type %v", paramType),
			Required:    &required,
		})
	}

	return &promptTemplate{
		handler:     handler,
		description: description,
		arguments:   arguments,
	}, nil
}

// renderPrompt renders a prompt with the given arguments
func (s *Server) renderPrompt(prompt Prompt, args map[string]string) ([]protocol.PromptMessage, error) {
	// Convert arguments to reflect.Values
	handlerType := reflect.TypeOf(prompt.Handler)
	handlerArgs := make([]reflect.Value, handlerType.NumIn())

	for i := 0; i < handlerType.NumIn(); i++ {
		paramType := handlerType.In(i)
		argName := fmt.Sprintf("arg%d", i)

		if argValue, ok := args[argName]; ok {
			// Create a new value of the parameter type
			paramValue := reflect.New(paramType).Interface()
			if err := convertValue(argValue, paramValue); err != nil {
				return nil, fmt.Errorf("invalid argument %s: %w", argName, err)
			}
			handlerArgs[i] = reflect.ValueOf(paramValue).Elem()
		} else {
			handlerArgs[i] = reflect.Zero(paramType)
		}
	}

	// Call the handler
	results := reflect.ValueOf(prompt.Handler).Call(handlerArgs)

	// Process results
	var messages []protocol.PromptMessage

	switch result := results[0].Interface().(type) {
	case string:
		// Single message template
		messages = []protocol.PromptMessage{
			{
				Role: protocol.RoleAssistant,
				Content: protocol.TextContent{
					Type: "text",
					Text: result,
				},
			},
		}
	case []protocol.PromptMessage:
		// Multiple messages
		messages = result
	default:
		return nil, fmt.Errorf("invalid prompt handler return type: %T", result)
	}

	return messages, nil
}

// renderTemplate renders a text template with the given arguments
func renderTemplate(tmpl string, args map[string]interface{}) (string, error) {
	t, err := template.New("prompt").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("invalid template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, args); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	return buf.String(), nil
}
