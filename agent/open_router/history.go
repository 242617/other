package open_router

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/242617/other/agent"
)

func encode(message message) (string, error) {
	b, err := json.Marshal(message)
	if err != nil {
		return "", errors.Wrap(err, "json marshal")
	}
	return string(b), nil
}

func decode(item string) (message, error) {
	var msg message
	if err := json.Unmarshal([]byte(item), &msg); err != nil {
		return message{}, errors.Wrap(err, "json unmarshal")
	}
	return msg, nil
}

func (p *OpenRouter) appendToHistory(storage agent.HistoryStorage, messages ...message) error {
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

func (p *OpenRouter) getHistory(storage agent.HistoryStorage) ([]message, error) {
	list, err := storage.Range()
	if err != nil {
		return nil, errors.Wrap(err, "storage range")
	}
	history := make([]message, len(list))
	for i, item := range list {
		msg, err := p.decode(item)
		if err != nil {
			return nil, errors.Wrap(err, "decode message")
		}
		history[i] = msg
	}
	return history, nil
}
