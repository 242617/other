package lms

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	openai "github.com/sashabaranov/go-openai"

	"github.com/242617/other/agent"
)

type (
	encodeFunc = func(message openai.ChatCompletionMessage) (string, error)
	decodeFunc = func(item string) (openai.ChatCompletionMessage, error)
)

func New(modifiers ...Modifier) *LMS {
	var lms LMS

	for _, modifier := range append([]Modifier{withDefaultClientTimeout()}, modifiers...) {
		modifier(&lms)
	}

	config := openai.DefaultConfig(lms.token)
	config.BaseURL = lms.address
	config.HTTPClient = &http.Client{Timeout: lms.timeout}

	return &LMS{
		encode: encode,
		decode: decode,
		client: openai.NewClientWithConfig(config),
		name:   "lms assistant",
	}
}

type LMS struct {
	token   string
	address string
	timeout time.Duration

	encode encodeFunc
	decode decodeFunc
	client *openai.Client
	name   string
}

func (p *LMS) EncodeSystemMessage(msg agent.Message) (string, error) {
	return p.encode(openai.ChatCompletionMessage{
		Role:    msg.Role,
		Content: msg.Text,
	})
}

func (p *LMS) loadHistory(storage agent.HistoryStorage) ([]openai.ChatCompletionMessage, error) {
	items, err := storage.List()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list history items")
	}

	history := make([]openai.ChatCompletionMessage, 0, len(items))
	for i, item := range items {
		msg, err := p.decode(item)
		if err != nil {
			slog.Warn("failed to decode history item", "index", i, "error", err)
			continue // Skip invalid items
		}
		history = append(history, msg)
	}

	return history, nil
}

func (p *LMS) saveHistory(storage agent.HistoryStorage, messages ...openai.ChatCompletionMessage) error {
	if len(messages) == 0 {
		return nil
	}

	items := make([]string, 0, len(messages))
	for _, msg := range messages {
		item, err := p.encode(msg)
		if err != nil {
			return errors.Wrap(err, "failed to encode message")
		}
		items = append(items, item)
	}

	return storage.Append(items...)
}

func (p *LMS) processToolCalls(
	ctx context.Context,
	toolCalls []openai.ToolCall,
	tools agent.Tools,
) []openai.ChatCompletionMessage {
	if len(toolCalls) == 0 {
		return nil
	}

	toolMessages := make([]openai.ChatCompletionMessage, 0, len(toolCalls))

	for _, toolCall := range toolCalls {
		var matchedTool agent.Tool
		for _, tool := range tools {
			if tool.Name() == toolCall.Function.Name {
				matchedTool = tool
				break
			}
		}

		if matchedTool == nil {
			slog.Warn("tool not found", "tool_name", toolCall.Function.Name)
			toolMessages = append(toolMessages, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    fmt.Sprintf("Error: Tool '%s' not found", toolCall.Function.Name),
				ToolCallID: toolCall.ID,
			})
			continue
		}

		// Execute tool
		result := matchedTool.Call(ctx, toolCall.Function.Arguments)

		// Add tool response message
		toolMessages = append(toolMessages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    result,
			ToolCallID: toolCall.ID,
		})
	}

	return toolMessages
}

func toAgentMessage(msg openai.ChatCompletionMessage) agent.Message {
	var extra string

	if len(msg.ToolCalls) > 0 {
		tools := make([]string, len(msg.ToolCalls))
		for i, toolCall := range msg.ToolCalls {
			tools[i] = fmt.Sprintf("%s(%s)", toolCall.Function.Name, toolCall.Function.Arguments)
		}
		extra += fmt.Sprintf("tools: [%s]", strings.Join(tools, ", "))
	}

	if msg.ToolCallID != "" {
		if extra != "" {
			extra += "; "
		}
		extra += fmt.Sprintf("tool_call_id: %s", msg.ToolCallID)
	}

	return agent.Message{
		Role:  msg.Role,
		Text:  msg.Content,
		Extra: extra,
	}
}

func convertTools(toolInfos []agent.ToolInfo) []openai.Tool {
	if len(toolInfos) == 0 {
		return nil
	}

	tools := make([]openai.Tool, len(toolInfos))
	for i, toolInfo := range toolInfos {
		tools[i] = openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        toolInfo.Function.Name,
				Description: toolInfo.Function.Description,
				Parameters:  convertToolParameters(toolInfo.Function.Parameters),
			},
		}
	}

	return tools
}

func convertToolParameters(params agent.ToolInfoFunctionParameters) map[string]interface{} {
	result := map[string]interface{}{
		"type": params.Type,
	}

	if len(params.Properties) > 0 {
		properties := make(map[string]interface{})
		for name, prop := range params.Properties {
			propMap := map[string]interface{}{
				"type": prop.Type,
			}
			if prop.Description != "" {
				propMap["description"] = prop.Description
			}
			if prop.Items != nil {
				propMap["items"] = prop.Items
			}
			properties[name] = propMap
		}
		result["properties"] = properties
	}

	if len(params.Required) > 0 {
		result["required"] = params.Required
	}

	return result
}
