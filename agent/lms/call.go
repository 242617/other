package lms

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"
	openai "github.com/sashabaranov/go-openai"

	"github.com/242617/other/agent"
)

func (p *LMS) Call(
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

	resChan := make(chan *agent.Content, 1)

	h, err := p.loadHistory(storage)
	if err != nil {
		return nil, errors.Wrap(err, "load history")
	}

	defer func(startFrom int) {
		if err := p.saveHistory(storage, h[startFrom:]...); err != nil {
			slog.Error("failed to save history", "err", err,
				"count", h[startFrom:],
			)
			panic(err) // TODO: Get rid of panic
		}
	}(len(h))

	var call func(messages ...openai.ChatCompletionMessage) error
	call = func(messages ...openai.ChatCompletionMessage) error {
		for _, msg := range messages {
			onMessage(toAgentMessage(msg))
		}

		h = append(h, messages...)

		req := openai.ChatCompletionRequest{
			Model:    model,
			Messages: h,
			Tools:    convertTools(tools.Info()),
		}

		resp, err := p.client.CreateChatCompletion(ctx, req)
		if err != nil {
			slog.Error("create chat completion", "err", err, "req", req)
			return errors.Wrap(err, "create chat completion")
		}

		if len(resp.Choices) == 0 {
			return errors.New("no choices in response")
		}

		responseMsg := resp.Choices[0].Message
		onMessage(toAgentMessage(responseMsg))
		h = append(h, responseMsg)

		// Check if tools were called
		if len(responseMsg.ToolCalls) == 0 {
			// No tools called, return final response
			resChan <- &agent.Content{
				Type:    agent.ContentTypeText,
				Content: responseMsg.Content,
			}
			return nil
		}

		// Process tool calls
		toolMessages := p.processToolCalls(ctx, responseMsg.ToolCalls, tools)
		if len(toolMessages) == 0 {
			// No valid tool responses, return the assistant's message as-is
			resChan <- &agent.Content{
				Type:    agent.ContentTypeText,
				Content: responseMsg.Content,
			}
			return nil
		}

		// Continue conversation with tool responses
		return call(toolMessages...)
	}

	// Convert agent.Content to OpenAI multimodal message parts
	userMessageParts := make([]openai.ChatMessagePart, 0, len(content))
	for _, c := range content {
		switch c.Type {
		case agent.ContentTypeText:
			userMessageParts = append(userMessageParts, openai.ChatMessagePart{
				Type: openai.ChatMessagePartTypeText,
				Text: c.Content,
			})
		case agent.ContentTypeImage:
			// Content already contains base64-encoded image data
			mimeType := c.MimeType
			if mimeType == "" {
				mimeType = "image/png" // Default to PNG
			}
			userMessageParts = append(userMessageParts, openai.ChatMessagePart{
				Type: openai.ChatMessagePartTypeImageURL,
				ImageURL: &openai.ChatMessageImageURL{
					URL:    "data:" + mimeType + ";base64," + c.Content,
					Detail: openai.ImageURLDetailAuto,
				},
			})
		}
	}

	userMessage := openai.ChatCompletionMessage{
		Role:         openai.ChatMessageRoleUser,
		MultiContent: userMessageParts,
	}

	if err := call(userMessage); err != nil {
		slog.Error("conversation failed", "error", err)
		return nil, errors.Wrap(err, "conversation failed")
	}

	select {
	case response := <-resChan:
		return response, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
