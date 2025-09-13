package tools

import (
	"context"

	"github.com/ollama/ollama/api"
)

type Tool interface {
	Name() string
	Tool() api.Tool
	Call(ctx context.Context, args map[string]any) (string, error)
}
