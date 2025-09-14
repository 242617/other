package open_router

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"

	"github.com/242617/other/agent"
)

func (p *OpenRouter) Call(ctx context.Context, model string, tools agent.Tools, text string, storage agent.HistoryStorage, onMessage agent.MessageCallback) (string, error) {
	if p.encode == nil {
		return "", errors.New("empty encode")
	}
	if p.decode == nil {
		return "", errors.New("empty decode")
	}
	if p.token == "" {
		return "", errors.New("empty token")
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

	call := func(messages ...message) error {
		for _, message := range messages {
			onMessage(message.ToMessage())
		}
		h = append(h, messages...) // Add messages to model
		req := request{
			Model:    model,
			Messages: h,
			Tools:    tools.Info(),
		}
		if err := p.completions(ctx, req, fn, tools); err != nil {
			slog.Error("completions", "err", err, "req", req)
			return errors.Wrap(err, "completions")
		}
		return nil
	}

	fn = func(res response) error {
		if len(res.Choices) != 1 {
			slog.Error("unexpected choices length", "choices length", len(res.Choices), "res", res)
			panic("unexpected choices length")
		}

		msg := res.Choices[0].Message
		onMessage(msg.ToMessage())
		h = append(h, msg) // Add message from model

		if len(msg.ToolCalls) == 0 { // TODO: Check "finish_reason"
			resCh <- msg.Content
			return nil
		}

		var messages []message
		for _, toolCall := range msg.ToolCalls {
			for _, tool := range tools {
				if toolCall.Function.Name == tool.Name() {

					result, err := tool.Call(ctx, toolCall.Function.Arguments)
					if err != nil {
						slog.Error("tool call", "err", err, "arguments", toolCall.Function.Arguments)
						return errors.Wrap(err, "tool call")
					}

					messages = append(messages,
						message{
							Role:       "tool",
							ToolCallID: toolCall.ID,
							Content:    result,
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

	message := message{Role: "user", Content: text}
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
