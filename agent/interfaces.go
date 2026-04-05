package agent

import "context"

type Assistant interface {
	CallText(ctx context.Context, text string) (string, error)
	Call(ctx context.Context, content ...Content) (*Content, error)
}

type HistoryStorage interface {
	Append(items ...string) error
	List() ([]string, error)
}

type Provider interface {
	Call(ctx context.Context, model string, tools Tools, content []Content, storage HistoryStorage, onMessage MessageCallback) (*Content, error)
	EncodeSystemMessage(message Message) (string, error)
}

// Tool interface is for tool implementation
// Contains methods for describing and calling tool
type Tool interface {
	Name() string
	Info() ToolInfo
	Call(ctx context.Context, args string) string
}
