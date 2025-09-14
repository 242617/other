package open_router

import (
	"fmt"
	"strings"

	"github.com/242617/other/agent"
)

type (
	encodeFunc  = func(message message) (string, error)
	dencodeFunc = func(item string) (message, error)

	responseFunc = func(res response) error
)

func New(token string) *OpenRouter {
	return &OpenRouter{
		encode: encode,
		decode: decode,
		debug:  false,
		token:  token,
		name:   "open router assistant",
	}
}

type OpenRouter struct {
	encode encodeFunc
	decode dencodeFunc
	debug  bool
	token  string
	name   string
}

type message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCallID string     `json:"tool_call_id"`
	ToolCalls  []toolCall `json:"tool_calls"`
}

func (m message) ToMessage() agent.Message {
	var extra string
	if len(m.ToolCalls) > 0 {
		tools := make([]string, len(m.ToolCalls))
		for i, toolCall := range m.ToolCalls {
			tools[i] = toolCall.Function.Name + toolCall.Function.Arguments
		}
		extra += fmt.Sprintf("tools: %q", strings.Join(tools, ", "))
	}

	if m.ToolCallID != "" {
		extra += fmt.Sprintf("tool name: %q", m.ToolCallID)
	}

	return agent.Message{
		Role:  m.Role,
		Text:  m.Content,
		Extra: extra,
	}
}

type toolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type request struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Tools    []tool    `json:"tools"`
}

type tool = agent.ToolInfo

type response struct {
	ID      string `json:"id"`
	Choices []struct {
		Message      message `json:"message"`
		FinishReason string  `json:"finish_reason"` // "tool_calls"
	} `json:"choices"`
	Usage Usage `json:"usage"`
}

func (p *OpenRouter) EncodeSystemMessage(msg agent.Message) (string, error) {
	return p.encode(message{Role: msg.Role, Content: msg.Text})
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
