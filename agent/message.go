package agent

import (
	"encoding/base64"
	"fmt"
	"os"
)

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
// For ContentTypeText: Content is the text string
// For ContentTypeImage: Content is base64-encoded image data
type Content struct {
	Type     ContentType `json:"type"`
	Content  string      `json:"content"`
	MimeType string      `json:"mime_type,omitempty"` // Required for images (e.g., "image/png", "image/jpeg")
}

// NewTextContent creates a text content
func NewTextContent(text string) Content {
	return Content{
		Type:    ContentTypeText,
		Content: text,
	}
}

// NewImageContent creates an image content from base64-encoded data
func NewImageContent(base64Data, mimeType string) Content {
	return Content{
		Type:     ContentTypeImage,
		Content:  base64Data,
		MimeType: mimeType,
	}
}

// NewImageContentFromFile loads an image file and creates a Content with base64-encoded data
func NewImageContentFromFile(filePath, mimeType string) (Content, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Content{}, fmt.Errorf("read image file: %w", err)
	}

	if mimeType == "" {
		mimeType = "image/png" // Default
	}

	return Content{
		Type:     ContentTypeImage,
		Content:  base64.StdEncoding.EncodeToString(data),
		MimeType: mimeType,
	}, nil
}
