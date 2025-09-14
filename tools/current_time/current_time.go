package current_time

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/242617/other/agent"
)

func New() *CurrentTime { return &CurrentTime{} }

type CurrentTime struct{}

func (CurrentTime) Name() string { return "get_current_time" }
func (t CurrentTime) Info() agent.ToolInfo {
	return agent.ToolInfo{
		Type: "function",
		Function: agent.ToolInfoFunction{
			Name:        t.Name(),
			Description: "Get the current time in UTC for one particular time zone.",
			Parameters: agent.ToolInfoFunctionParameters{
				Type: "object",
				Properties: map[string]agent.ToolInfoFunctionParametersProperty{
					"tz": {
						Type: "string",
						Description: strings.Join([]string{
							"One timezone to get the current time in.",
							`For example, "Europe/Moscow".`,
						}, "\n"),
					},
				},
				Required: []string{"tz"},
			},
		},
	}
}

type Args struct {
	TimeZone string `json:"tz"`
}

func (t *CurrentTime) Call(_ context.Context, raw string) (string, error) {
	var args Args
	if err := json.Unmarshal([]byte(raw), &args); err != nil {
		return "", errors.Wrap(err, "json unmarshal")
	}

	loc, err := time.LoadLocation(args.TimeZone)
	if err != nil {
		return "", errors.Wrap(err, "time load location")
	}

	return time.Now().In(loc).Format(time.DateTime), nil
}
