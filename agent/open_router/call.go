package open_router

import (
	"context"
	"encoding/base64"
	"log/slog"
	"os"

	"github.com/pkg/errors"

	"github.com/242617/other/agent"
)

func (p *OpenRouter) Call(ctx context.Context, model string, tools agent.Tools, content []agent.Content, storage agent.HistoryStorage, onMessage agent.MessageCallback) (*agent.Content, error) {
	if p.encode == nil {
		return nil, errors.New("empty encode")
	}
	if p.decode == nil {
		return nil, errors.New("empty decode")
	}
	if p.token == "" {
		return nil, errors.New("empty token")
	}

	resCh := make(chan *agent.Content, 1)

	h, err := p.getHistory(storage)
	if err != nil {
		return nil, errors.Wrap(err, "list to history")
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
			resCh <- &agent.Content{
				Type:    agent.ContentTypeText,
				Content: msg.Content,
			}
			return nil
		}

		var messages []message
		for _, toolCall := range msg.ToolCalls {
			for _, tool := range tools {
				if toolCall.Function.Name == tool.Name() {

					messages = append(messages,
						message{
							Role:       "tool",
							ToolCallID: toolCall.ID,
							Content:    tool.Call(ctx, toolCall.Function.Arguments),
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

	// OpenRouter supports multimodal content similar to OpenAI API
	message := message{Role: "user"}
	for _, c := range content {
		switch c.Type {
		case agent.ContentTypeText:
			message.Content = c.Content
		case agent.ContentTypeImage:
			// Load image file and encode to base64
			imageData, err := os.ReadFile(c.Content)
			if err != nil {
				return nil, errors.Wrap(err, "read image file")
			}
			base64Data := base64.StdEncoding.EncodeToString(imageData)
			mimeType := c.MimeType
			if mimeType == "" {
				mimeType = "image/png" // Default to PNG
			}
			// OpenRouter expects image URLs or base64 data URIs
			message.Content += " [image:" + "data:" + mimeType + ";base64," + base64Data + "]"
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
