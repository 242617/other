package document

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

/// Utils

const (
	trailing = "---\n"
	leading  = "\n---\n"
)

func unwrap(str string) string {
	str = strings.TrimPrefix(str, trailing)
	str = strings.TrimSuffix(str, leading)
	return str
}
func wrap(str string) string {
	return trailing + str + leading
}

func remove(doc string, path string) (string, error) {
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(doc), &root); err != nil {
		return "", fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	trimmedPath := strings.TrimPrefix(path, ".")
	if trimmedPath == "" {
		return "", fmt.Errorf("empty path")
	}

	keys := strings.Split(trimmedPath, ".")
	if err := removeNestedKey(&root, keys); err != nil {
		return "", err
	}

	out, err := yaml.Marshal(&root)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return string(out), nil
}

var ErrKeyNotFound = errors.New("key not found")

func removeNestedKey(node *yaml.Node, keys []string) error {
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return fmt.Errorf("empty document")
		}
		return removeNestedKey(node.Content[0], keys)
	}

	// if node.Kind != yaml.MappingNode {
	// 	return fmt.Errorf("expected mapping node")
	// }

	if len(keys) == 1 {
		// Последний ключ - удаляем его
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Value == keys[0] {
				node.Content = append(node.Content[:i], node.Content[i+2:]...)
				return nil
			}
		}
		return errors.Wrap(ErrKeyNotFound, fmt.Sprintf("key %q not found", keys[0]))
	}

	// Ищем следующий уровень вложенности
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == keys[0] {
			return removeNestedKey(node.Content[i+1], keys[1:])
		}
	}
	return errors.Wrap(ErrKeyNotFound, fmt.Sprintf("key %q not found", keys[0]))
}
