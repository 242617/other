package agent

import (
	"context"
	"strings"

	"github.com/pkg/errors"
)

func New[T Configurer](items ...T) (*Agent, error) {
	var a Agent

	configurers := []any{
		withDefaultModel(),
		withDefaultOptions(),
		withDefaultSystem(),
		withDefaultStorage(),
		withDefaultOnMessageFunc(),
	}
	for _, item := range items {
		configurers = append(configurers, item)
	}

	for _, configurer := range configurers {
		switch configurer := any(configurer).(type) {
		case Modifier:
			configurer(&a)
		case Option:
			if err := configurer(&a); err != nil {
				return nil, errors.Wrap(err, "apply option")
			}
		default:
			return nil, errors.New("unexpected type")
		}
	}

	if a.provider == nil {
		return nil, errors.New("empty provider")
	}
	if a.model == "" {
		return nil, errors.New("empty model")
	}
	if a.storage == nil {
		return nil, errors.New("empty storage")
	}
	if a.messageCallback == nil {
		return nil, errors.New("empty message callback")
	}

	return &a, nil
}

type Agent struct {
	provider        Provider
	model           string
	tools           Tools
	options         map[string]any
	system          string
	storage         HistoryStorage
	messageCallback MessageCallback
}

// Session starts a new session
func (a *Agent) Session(ctx context.Context) (*Session, error) {
	encoded, err := a.provider.EncodeSystemMessage(Message{Role: "system", Text: a.system})
	if err != nil {
		return nil, errors.Wrap(err, "provider encode system message")
	}
	if err := a.storage.Append(encoded); err != nil {
		return nil, errors.Wrap(err, "storage rpush")
	}

	s := Session{agent: a}
	return &s, nil
}

type Session struct {
	agent *Agent
}

func (s *Session) Call(ctx context.Context, text string) (string, error) {
	res, err := s.agent.provider.Call(ctx, s.agent.model, s.agent.tools, text, s.agent.storage, s.agent.messageCallback)
	if err != nil {
		return "", errors.Wrap(err, "provider call")
	}
	return res, nil
}

func (s *Session) History() string {
	list, err := s.agent.storage.List()
	if err != nil {
		panic(err)
	}
	return strings.Join(list, "\n\n")
}
