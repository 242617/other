package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/pkg/errors"

	"github.com/242617/other/tools"
)

func New(modifiers ...Modifier) (*Agent, error) {
	var a Agent

	defaultModifiers := []Modifier{withDefaultClient(), withDefaultModel(), withDefaultSystem(), withDefaultOptions()}
	for _, modifier := range append(defaultModifiers, modifiers...) {
		modifier(&a)
	}

	if a.client == nil {
		return nil, errors.New("empty client")
	}
	if a.model == "" {
		return nil, errors.New("empty model")
	}

	return &a, nil
}

type Agent struct {
	client  *api.Client
	model   string
	options map[string]any
	system  string
	tools   []tools.Tool
}

func (a *Agent) Tools() []api.Tool {
	tools := make([]api.Tool, len(a.tools))
	for i, tool := range a.tools {
		tools[i] = tool.Tool()
	}
	return tools
}

func (a *Agent) Session(ctx context.Context) *Session {
	s := Session{
		agent:   a,
		history: []api.Message{{Role: "system", Content: a.system}},
	}
	return &s
}

type Session struct {
	agent   *Agent
	history []api.Message
}

func (s *Session) String() string {
	res := make([]string, len(s.history))
	for i, message := range s.history {
		res[i] = fmt.Sprintf("> [%s]: %s", message.Role, message.Content)
	}
	return strings.Join(res, "\n")
}

func (s *Session) Call(ctx context.Context, msg string) (string, error) {
	resCh := make(chan string, 1)

	message := api.Message{Role: "user", Content: msg}
	var fn func(res api.ChatResponse) error

	call := func(messages ...api.Message) error {
		s.history = append(s.history, messages...) // Add messages to model
		req := api.ChatRequest{
			Model:    s.agent.model,
			Stream:   new(bool),
			Messages: s.history,
			Tools:    s.agent.Tools(),
			Options:  s.agent.options,
		}
		if err := s.agent.client.Chat(ctx, &req, fn); err != nil {
			return errors.Wrap(err, "session agent client chat")
		}
		return nil
	}

	fn = func(res api.ChatResponse) error {
		s.history = append(s.history, res.Message) // Add message from model

		if len(res.Message.ToolCalls) == 0 {
			resCh <- res.Message.Content
			return nil
		}

		var messages []api.Message
		for _, toolCall := range res.Message.ToolCalls {
			for _, tool := range s.agent.tools {
				if toolCall.Function.Name == tool.Name() {
					result, err := tool.Call(ctx, toolCall.Function.Arguments)
					if err != nil {
						return errors.Wrap(err, "tool call")
					}

					messages = append(messages,
						api.Message{
							Role:     "tool",
							Content:  result,
							ToolName: tool.Name(),
						},
					)
				}
			}
		}

		if err := call(messages...); err != nil {
			return errors.Wrap(err, "call")
		}

		return nil
	}

	if err := call(message); err != nil {
		return "", errors.Wrap(err, "call")
	}

	select {
	case msg := <-resCh:
		return msg, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
