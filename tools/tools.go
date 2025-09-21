package tools

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/242617/other/agent"
)

type CallFunc = func(ctx context.Context, raw string) string

func CreateToolInfo(name, description string, args any) (agent.ToolInfo, error) {
	t := reflect.TypeOf(args)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return agent.ToolInfo{}, fmt.Errorf("args must be a struct, not %s", t.Kind())
	}

	properties := make(map[string]agent.ToolInfoFunctionParametersProperty)
	var required []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		jsonName := strings.Split(jsonTag, ",")[0]
		if jsonName == "" {
			continue
		}
		description := field.Tag.Get("description")

		var paramType string
		switch field.Type.Kind() {
		case reflect.String:
			paramType = "string"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			paramType = "number"
		case reflect.Bool:
			paramType = "boolean"
		case reflect.Slice:
			paramType = "array"
		default:
			return agent.ToolInfo{}, fmt.Errorf("unexpected field type %q for field %q", field.Type.Kind(), field.Name)
		}

		properties[jsonName] = agent.ToolInfoFunctionParametersProperty{
			Type:        paramType,
			Description: description,
		}

		if !strings.Contains(jsonTag, "omitempty") {
			required = append(required, jsonName)
		}
	}

	return agent.ToolInfo{
		Type: "function",
		Function: agent.ToolInfoFunction{
			Name:        name,
			Description: description,
			Parameters: agent.ToolInfoFunctionParameters{
				Type:       "object",
				Properties: properties,
				Required:   required,
			},
		},
	}, nil
}
