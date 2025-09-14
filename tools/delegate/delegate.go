package delegate

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/242617/other/agent"
)

func New(role string, assistant agent.Assistant) *Delegate {
	return &Delegate{role: role, assistant: assistant}
}

type Delegate struct {
	role      string
	assistant agent.Assistant
}

func (t Delegate) Name() string { return fmt.Sprintf("delegate_%s", t.role) }
func (t Delegate) Info() agent.ToolInfo {
	return agent.ToolInfo{
		Type: "function",
		Function: agent.ToolInfoFunction{
			Name: t.Name(),
			Description: strings.Join([]string{
				"Delegate a task to assistant. Describe task thoroughly.",
				fmt.Sprintf("Assistant's role is %q.", t.role),
			}, "\n"),
			Parameters: agent.ToolInfoFunctionParameters{
				Type: "object",
				Properties: map[string]agent.ToolInfoFunctionParametersProperty{
					"task": {
						Type:        "string",
						Description: "Task description to perform.",
					},
				},
				Required: []string{"task"},
			},
		},
	}
}

type Args struct {
	Task string `json:"task"`
}

func (t *Delegate) Call(ctx context.Context, raw string) (string, error) {
	var args Args
	if err := json.Unmarshal([]byte(raw), &args); err != nil {
		return "", errors.Wrap(err, "json unmarshal")
	}
	return t.assistant.Call(ctx, args.Task)
}
