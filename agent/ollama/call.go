package ollama

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/ollama/ollama/api"
	"github.com/pkg/errors"

	"github.com/242617/other/agent"
)

func (p *Ollama) Call(ctx context.Context, model string, tools agent.Tools, text string, storage agent.HistoryStorage, onMessage agent.MessageCallback) (string, error) {
	if p.encode == nil {
		return "", errors.New("empty encode")
	}
	if p.decode == nil {
		return "", errors.New("empty decode")
	}

	resCh := make(chan string, 1)

	h, err := p.getHistory(storage)
	if err != nil {
		return "", errors.Wrap(err, "list to history")
	}

	defer func(startFrom int) {
		if err := p.appendToHistory(storage, h[startFrom:]...); err != nil {
			slog.Error("append to history", "err", err, "count", h[startFrom:])
			panic(err) // TODO: Get rid of panic
		}
	}(len(h))

	var fn responseFunc

	call := func(messages ...api.Message) error {
		for _, message := range messages {
			onMessage(toMessage(message))
		}
		h = append(h, messages...) // Add messages to model
		req := api.ChatRequest{
			Model:    model,
			Messages: h,
			Stream:   new(bool),
			Tools:    toTools(tools.Info()),
			// Think:    &api.ThinkValue{Value: "low"},
		}
		if err := p.client.Chat(ctx, &req, fn); err != nil {
			slog.Error("completions", "err", err, "req", req)
			return errors.Wrap(err, "session agent client chat")
		}
		return nil
	}

	fn = func(res api.ChatResponse) error {
		msg := res.Message
		onMessage(toMessage(msg))
		h = append(h, msg) // Add message from model

		if len(msg.ToolCalls) == 0 { // TODO: Check "finish_reason"
			resCh <- msg.Content
			return nil
		}

		var messages []api.Message
		for _, toolCall := range msg.ToolCalls {
			for _, tool := range tools {
				if toolCall.Function.Name == tool.Name() {

					b, err := json.Marshal(toolCall.Function.Arguments)
					if err != nil {
						slog.Error("json marshal", "err", err, "arguments", toolCall.Function.Arguments)
						return errors.Wrap(err, "json marshal")
					}

					messages = append(messages,
						api.Message{
							Role:     "tool",
							Content:  tool.Call(ctx, string(b)),
							ToolName: tool.Name(),
						},
					)

					break

				}
			}
		}

		if err := call(messages...); err != nil {
			return errors.Wrap(err, "call")
		}

		return nil
	}

	message := api.Message{Role: "user", Content: text}
	if err := call(message); err != nil {
		slog.Error("call", "err", err, "message", message)
		return "", errors.Wrap(err, "call")
	}

	select {
	case msg := <-resCh:
		return msg, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
