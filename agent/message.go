package agent

import "fmt"

type MessageCallback = func(msg Message)

type Message struct {
	Role  string
	Text  string
	Extra string
}

func (m Message) String() string { return fmt.Sprintf("> [%s]: %s", m.Role, m.Text) }

// Content types for multimodal messages
type ContentType string

const (
	ContentTypeText  ContentType = "text"
	ContentTypeImage ContentType = "image"
)

// Content represents a piece of content in a multimodal message
type Content struct {
	Type     ContentType `json:"type"`
	Content  string      `json:"content"`
	MimeType string      `json:"mime_type,omitempty"` // Optional for images (e.g., "image/png", "image/jpeg")
}
