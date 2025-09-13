package current_time

import (
	"context"
	"fmt"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/pkg/errors"
)

func New() *CurrentTime { return &CurrentTime{} }

type CurrentTime struct{}

func (CurrentTime) Name() string { return "get_current_time" }
func (ct CurrentTime) Tool() api.Tool {
	toolFunction := api.ToolFunction{
		Name:        ct.Name(),
		Description: "Get the current time in UTC for one particular time zone.",
	}
	toolFunction.Parameters.Properties = map[string]api.ToolProperty{
		"tz": {
			Type:        []string{"string"},
			Description: "One timezone to get the current time in.",
		},
	}
	toolFunction.Parameters.Required = []string{"tz"}

	return api.Tool{
		Type:     "function",
		Function: toolFunction,
	}
}

func (ct *CurrentTime) Call(_ context.Context, args map[string]any) (string, error) {
	tz, ok := args["tz"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected type for tz: %T (%v)", args["tz"], args["tz"])
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return "", errors.Wrap(err, "time load location")
	}

	return time.Now().In(loc).Format(time.DateTime), nil
}
