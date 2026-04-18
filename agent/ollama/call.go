package ollama

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log/slog"

	"github.com/ollama/ollama/api"
	"github.com/pkg/errors"

	"github.com/242617/other/agent"
)

func (p *Ollama) Call(
	ctx context.Context,
	model string,
	tools agent.Tools,
	content []agent.Content,
	storage agent.HistoryStorage,
	onMessage agent.MessageCallback,
) (*agent.Content, error) {
	if p.encode == nil {
		return nil, errors.New("empty encode")
	}
	if p.decode == nil {
		return nil, errors.New("empty decode")
	}
	if onMessage == nil {
		onMessage = func(agent.Message) {}
	}

	resCh := make(chan *agent.Content, 1)

	h, err := p.loadHistory(storage)
	if err != nil {
		return nil, errors.Wrap(err, "load history")
	}

	defer func(startFrom int) {
		if err := p.appendToHistory(storage, h[startFrom:]...); err != nil {
			slog.Error("append to history", "err", err,
				"count", h[startFrom:],
			)
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
			slog.Error("chat", "err", err, "req", req)
			return errors.Wrap(err, "session agent client chat")
		}
		return nil
	}

	fn = func(res api.ChatResponse) error {
		msg := res.Message
		onMessage(toMessage(msg))
		h = append(h, msg) // Add message from model

		if len(msg.ToolCalls) == 0 { // TODO: Check "finish_reason"
			resCh <- &agent.Content{
				Type:    agent.ContentTypeText,
				Content: msg.Content,
			}
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

	// Convert agent.Content to Ollama message with images
	message := api.Message{Role: "user"}
	for _, c := range content {
		switch c.Type {
		case agent.ContentTypeText:
			message.Content = c.Content
		case agent.ContentTypeImage:
			// Ollama expects raw binary images - decode base64
			imageData, err := base64.StdEncoding.DecodeString(c.Content)
			if err != nil {
				return nil, errors.Wrap(err, "decode base64 image")
			}
			message.Images = append(message.Images, imageData)
		}
	}

	if err := call(message); err != nil {
		slog.Error("call", "err", err, "message", message)
		return nil, errors.Wrap(err, "call")
	}

	select {
	case msg := <-resCh:
		return msg, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
