package inmemory_storage

import (
	"slices"
	"sync"
)

// New created inmemory storage for messages.
// total = pinned + shifting
func New(pinned, shifting int) *InmemoryStorage {
	return &InmemoryStorage{system: make([]string, 0, pinned), messages: make([]string, 0, shifting)}
}

type InmemoryStorage struct {
	system   []string
	messages []string
	mutex    sync.Mutex
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *InmemoryStorage) Append(messages ...string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(messages) == 0 {
		return nil
	}

	if systemCap := cap(s.system) - len(s.system); systemCap > 0 {
		n := min(systemCap, len(messages))
		s.system = append(s.system, messages[:n]...)
		messages = messages[n:]
	}

	for _, message := range messages {
		if len(s.messages) < cap(s.messages) {
			s.messages = append(s.messages, message)
		} else {
			s.messages = append(s.messages[1:], message)
		}
	}
	return nil
}

func (s *InmemoryStorage) List() ([]string, error) {
	system := slices.Clone(s.system)
	messages := slices.Clone(s.messages)
	return append(system, messages...), nil
}
