package agent

import "fmt"

type MessageCallback = func(msg Message)

type Message struct {
	Role  string
	Text  string
	Extra string
}

func (m Message) String() string { return fmt.Sprintf("> [%s]: %s", m.Role, m.Text) }
