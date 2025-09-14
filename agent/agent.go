package agent

import (
	"context"
	"strings"

	"github.com/pkg/errors"
)

func New(modifiers ...Modifier) (*Agent, error) {
	var a Agent

	defaultModifiers := []Modifier{
		withDefaultModel(),
		withDefaultOptions(),
		withDefaultSystem(),
		withDefaultOnMessageFunc(),
	}
	for _, modifier := range append(defaultModifiers, modifiers...) {
		modifier(&a)
	}

	if a.provider == nil {
		return nil, errors.New("empty provider")
	}
	if a.model == "" {
		return nil, errors.New("empty model")
	}

	return &a, nil
}

type Agent struct {
	provider      Provider
	model         string
	tools         Tools
	options       map[string]any
	system        string
	onMessageFunc MessageCallback
}

// Session starts a new session
func (a *Agent) Session(ctx context.Context) (*Session, error) {
	storage := &inmemoryStorage{}

	encoded, err := a.provider.EncodeSystemMessage(Message{Role: "system", Text: a.system})
	if err != nil {
		return nil, errors.Wrap(err, "provider encode system message")
	}
	if err := storage.Rpush(encoded); err != nil {
		return nil, errors.Wrap(err, "storage rpush")
	}

	s := Session{
		agent:   a,
		storage: storage,
	}
	return &s, nil
}

type Session struct {
	agent   *Agent
	storage HistoryStorage
}

func (s *Session) Call(ctx context.Context, text string) (string, error) {
	res, err := s.agent.provider.Call(ctx, s.agent.model, s.agent.tools, text, s.storage, s.agent.onMessageFunc)
	if err != nil {
		return "", errors.Wrap(err, "provider call")
	}
	return res, nil
}

func (s *Session) History() string {
	list, err := s.storage.Range()
	if err != nil {
		panic(err)
	}
	return strings.Join(list, "\n\n")
}
