package lms

import (
	"encoding/json"

	"github.com/pkg/errors"
	openai "github.com/sashabaranov/go-openai"
)

func encode(message openai.ChatCompletionMessage) (string, error) {
	data, err := json.Marshal(message)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal message")
	}
	return string(data), nil
}

func decode(item string) (openai.ChatCompletionMessage, error) {
	var message openai.ChatCompletionMessage
	if err := json.Unmarshal([]byte(item), &message); err != nil {
		return openai.ChatCompletionMessage{}, errors.Wrap(err, "failed to unmarshal message")
	}
	return message, nil
}
