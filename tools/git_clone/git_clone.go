package git_clone

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/242617/other/agent"
)

func New() *GitClone { return &GitClone{} }

type GitClone struct{}

func (GitClone) Name() string { return "git_clone" }

func (t GitClone) Info() agent.ToolInfo {
	return agent.ToolInfo{
		Type: "function",
		Function: agent.ToolInfoFunction{
			Name:        t.Name(),
			Description: "Clone a Git repository into a specified directory.",
			Parameters: agent.ToolInfoFunctionParameters{
				Type: "object",
				Properties: map[string]agent.ToolInfoFunctionParametersProperty{
					"repository_name": {
						Type:        "string",
						Description: `Repository address to clone. Please use the full name of the repository, including the domain, group, and name, for example "https://github.com/go-task/task.git".`,
					},
					"directory": {
						Type:        "string",
						Description: "Directory to clone repository to.",
					},
				},
				Required: []string{"repository_name", "directory"},
			},
		},
	}
}

type Args struct {
	RepositoryName string `json:"repository_name"`
	Directory      string `json:"directory"`
}

func (t *GitClone) Call(_ context.Context, raw string) string {
	var args Args
	if err := json.Unmarshal([]byte(raw), &args); err != nil {
		return fmt.Sprintf("Cannot unmarshal arguments due to error: %q", err.Error())
	}
	slog.Debug("call", "args", args)

	cmd := exec.Command("git", "clone", args.RepositoryName, args.Directory)
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("Cannot clone %q into %q due to error: %q", args.RepositoryName, args.Directory, err.Error())
	}
	return fmt.Sprintf("✅ Repository %q cloned into %q successfully", args.RepositoryName, args.Directory)
}
