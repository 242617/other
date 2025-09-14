package agent

type inmemoryStorage struct {
	items []string
}

func (s *inmemoryStorage) Rpush(items ...string) error {
	s.items = append(s.items, items...)
	return nil
}

func (s *inmemoryStorage) Range() ([]string, error) {
	return s.items, nil
}
