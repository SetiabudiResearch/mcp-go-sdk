package server

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

// resourcePattern represents a parsed resource pattern
type resourcePattern struct {
	pattern     string
	regex       *regexp.Regexp
	paramNames  []string
	paramTypes  []reflect.Type
	handlerType reflect.Type
}

// parseResourcePattern parses a resource pattern into a regex and parameter info
func parseResourcePattern(pattern string, handler interface{}) (*resourcePattern, error) {
	// Validate handler
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		return nil, fmt.Errorf("handler must be a function")
	}

	// Extract parameter names and build regex
	var paramNames []string
	var paramTypes []reflect.Type
	regexStr := pattern

	// Find all {param} in pattern
	paramRegex := regexp.MustCompile(`\{([^}]+)\}`)
	matches := paramRegex.FindAllStringSubmatch(pattern, -1)

	for i, match := range matches {
		paramName := match[1]
		paramNames = append(paramNames, paramName)

		// Get parameter type from handler
		if i >= handlerType.NumIn() {
			return nil, fmt.Errorf("not enough parameters in handler for pattern %s", pattern)
		}
		paramTypes = append(paramTypes, handlerType.In(i))

		// Replace {param} with regex capture group
		regexStr = strings.Replace(regexStr, match[0], `([^/]+)`, 1)
	}

	// Compile the regex
	regex, err := regexp.Compile("^" + regexStr + "$")
	if err != nil {
		return nil, fmt.Errorf("invalid pattern %s: %w", pattern, err)
	}

	return &resourcePattern{
		pattern:     pattern,
		regex:       regex,
		paramNames:  paramNames,
		paramTypes:  paramTypes,
		handlerType: handlerType,
	}, nil
}

// matchResource finds a matching resource and extracts parameters
func (s *Server) matchResource(uri string) (Resource, map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	parsedURI, err := url.Parse(uri)
	if err != nil {
		return Resource{}, nil, fmt.Errorf("invalid URI: %w", err)
	}

	// Try to match each resource pattern
	for pattern, resource := range s.resources {
		resourcePattern, err := parseResourcePattern(pattern, resource.Handler)
		if err != nil {
			continue
		}

		matches := resourcePattern.regex.FindStringSubmatch(parsedURI.Path)
		if matches == nil {
			continue
		}

		// Extract parameters
		params := make(map[string]interface{})
		for i, name := range resourcePattern.paramNames {
			// Convert parameter value to the correct type
			paramValue := reflect.New(resourcePattern.paramTypes[i]).Interface()
			if err := convertValue(matches[i+1], paramValue); err != nil {
				return Resource{}, nil, fmt.Errorf("invalid parameter %s: %w", name, err)
			}
			params[name] = reflect.ValueOf(paramValue).Elem().Interface()
		}

		return resource, params, nil
	}

	return Resource{}, nil, fmt.Errorf("no matching resource found for %s", uri)
}

// readResource reads data from a resource using its handler
func (s *Server) readResource(resource Resource, params map[string]interface{}) ([]interface{}, error) {
	// Convert parameters to reflect.Values
	handlerType := reflect.TypeOf(resource.Handler)
	args := make([]reflect.Value, handlerType.NumIn())

	for i := 0; i < handlerType.NumIn(); i++ {
		paramType := handlerType.In(i)
		paramName := fmt.Sprintf("param%d", i)

		if paramValue, ok := params[paramName]; ok {
			args[i] = reflect.ValueOf(paramValue)
		} else {
			args[i] = reflect.Zero(paramType)
		}
	}

	// Call the handler
	results := reflect.ValueOf(resource.Handler).Call(args)

	// Process results
	var contents []interface{}

	if len(results) == 2 { // Handler returns (value, error)
		if !results[1].IsNil() {
			return nil, results[1].Interface().(error)
		}
		contents = []interface{}{results[0].Interface()}
	} else { // Handler returns single value
		contents = []interface{}{results[0].Interface()}
	}

	return contents, nil
}

// convertValue converts a string value to the target type
func convertValue(value string, target interface{}) error {
	v := reflect.ValueOf(target).Elem()

	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i int64
		if _, err := fmt.Sscanf(value, "%d", &i); err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var u uint64
		if _, err := fmt.Sscanf(value, "%d", &u); err != nil {
			return err
		}
		v.SetUint(u)
	case reflect.Float32, reflect.Float64:
		var f float64
		if _, err := fmt.Sscanf(value, "%f", &f); err != nil {
			return err
		}
		v.SetFloat(f)
	case reflect.Bool:
		var b bool
		if _, err := fmt.Sscanf(value, "%t", &b); err != nil {
			return err
		}
		v.SetBool(b)
	default:
		return fmt.Errorf("unsupported type: %v", v.Type())
	}

	return nil
}
