package ollama

import (
	"encoding/json"

	"github.com/ollama/ollama/api"
	"github.com/pkg/errors"

	"github.com/242617/other/agent"
)

func encode(message api.Message) (string, error) {
	b, err := json.Marshal(message)
	if err != nil {
		return "", errors.Wrap(err, "json marshal")
	}
	return string(b), nil
}

func decode(item string) (api.Message, error) {
	var msg api.Message
	if err := json.Unmarshal([]byte(item), &msg); err != nil {
		return api.Message{}, errors.Wrap(err, "json unmarshal")
	}
	return msg, nil
}

func (p *Ollama) appendToHistory(storage agent.HistoryStorage, messages ...api.Message) error {
	items := make([]string, len(messages))
	for i, message := range messages {
		item, err := p.encode(message)
		if err != nil {
			return errors.Wrap(err, "encode message")
		}
		items[i] = item
	}
	if err := storage.Rpush(items...); err != nil {
		return errors.Wrap(err, "storage rpush")
	}
	return nil
}

func (p *Ollama) getHistory(storage agent.HistoryStorage) ([]api.Message, error) {
	list, err := storage.Range()
	if err != nil {
		return nil, errors.Wrap(err, "storage range")
	}
	history := make([]api.Message, len(list))
	for i, item := range list {
		msg, err := p.decode(item)
		if err != nil {
			return nil, errors.Wrap(err, "decode message")
		}
		history[i] = msg
	}
	return history, nil
}
