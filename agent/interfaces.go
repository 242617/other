package agent

import "context"

type Assistant interface {
	Call(ctx context.Context, text string) (string, error)
}

type HistoryStorage interface {
	Rpush(items ...string) error
	Range() ([]string, error)
}

type Provider interface {
	Call(ctx context.Context, model string, tools Tools, text string, storage HistoryStorage, onMessage MessageCallback) (string, error)
	EncodeSystemMessage(message Message) (string, error)
}

// Tool interface is for tool implementation
// Contains methods for describing and calling tool
type Tool interface {
	Name() string
	Info() ToolInfo
	Call(ctx context.Context, args string) string
}
