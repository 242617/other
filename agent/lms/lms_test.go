package lms_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/242617/other/agent"
	"github.com/242617/other/agent/lms"
)

// mockHistoryStorage implements agent.HistoryStorage for testing
type mockHistoryStorage struct {
	items []string
}

func (m *mockHistoryStorage) Append(items ...string) error {
	m.items = append(m.items, items...)
	return nil
}

func (m *mockHistoryStorage) List() ([]string, error) {
	return m.items, nil
}

// mockTool implements agent.Tool for testing
type mockTool struct {
	name   string
	info   agent.ToolInfo
	result string
}

func (m *mockTool) Name() string {
	return m.name
}

func (m *mockTool) Info() agent.ToolInfo {
	return m.info
}

func (m *mockTool) Call(ctx context.Context, args string) string {
	return m.result
}

// createMockTool creates a simple mock tool for testing
func createMockTool(name, description, result string) *mockTool {
	return &mockTool{
		name: name,
		info: agent.ToolInfo{
			Type: "function",
			Function: agent.ToolInfoFunction{
				Name:        name,
				Description: description,
				Parameters: agent.ToolInfoFunctionParameters{
					Type: "object",
					Properties: map[string]agent.ToolInfoFunctionParametersProperty{
						"input": {
							Type:        "string",
							Description: "Input parameter",
						},
					},
					Required: []string{"input"},
				},
			},
		},
		result: result,
	}
}

// getTextFromMessage extracts text from OpenAI ChatCompletionMessage
// It handles both Content field and MultiContent field
func getTextFromMessage(msg openai.ChatCompletionMessage) string {
	if msg.Content != "" {
		return msg.Content
	}
	// If Content is empty, check MultiContent
	var textParts []string
	for _, part := range msg.MultiContent {
		if part.Type == openai.ChatMessagePartTypeText {
			textParts = append(textParts, part.Text)
		}
	}
	return strings.Join(textParts, " ")
}

func TestNew(t *testing.T) {
	provider := lms.New(lms.WithAddress("http://127.0.0.1:1234"))
	assert.NotNil(t, provider)
}

func TestEncodeSystemMessage(t *testing.T) {
	provider := lms.New(lms.WithAddress("http://127.0.0.1:1234"))

	agentMsg := agent.Message{
		Role: openai.ChatMessageRoleSystem,
		Text: "You are a helpful assistant",
	}

	encoded, err := provider.EncodeSystemMessage(agentMsg)
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)
}

// Integration test with mock HTTP server
func TestCallWithMockServer(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Accept both /v1/chat/completions and /chat/completions paths
		if r.URL.Path != "/v1/chat/completions" && r.URL.Path != "/chat/completions" {
			http.NotFound(w, r)
			return
		}

		// Parse request
		var req openai.ChatCompletionRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		// Create mock response
		resp := openai.ChatCompletionResponse{
			ID:      "chatcmpl-test",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   req.Model,
			Choices: []openai.ChatCompletionChoice{
				{
					Index: 0,
					Message: openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleAssistant,
						Content: "Mock response to: " + getTextFromMessage(req.Messages[len(req.Messages)-1]),
					},
					FinishReason: "stop",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create provider with mock server URL
	// Note: go-openai client automatically adds /v1 prefix to BaseURL
	provider := lms.New(lms.WithAddress(server.URL))

	// Test basic call
	ctx := context.Background()
	storage := &mockHistoryStorage{}
	tools := agent.Tools{}

	var receivedMessages []agent.Message
	onMessage := func(msg agent.Message) {
		receivedMessages = append(receivedMessages, msg)
	}

	response, err := provider.Call(ctx, "test-model", tools, []agent.Content{{Type: agent.ContentTypeText, Content: "Hello, world!"}}, storage, onMessage)
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, "Mock response to: Hello, world!", response.Content)

	// Check that messages were recorded
	assert.Len(t, receivedMessages, 2) // User message + assistant response
	assert.Equal(t, openai.ChatMessageRoleUser, receivedMessages[0].Role)
	assert.Equal(t, "Hello, world!", receivedMessages[0].Text)
	assert.Equal(t, openai.ChatMessageRoleAssistant, receivedMessages[1].Role)

	// Check that history was saved
	items, err := storage.List()
	require.NoError(t, err)
	assert.Len(t, items, 2) // User message + assistant response
}

// Test with tool calling
func TestCallWithToolCalling(t *testing.T) {
	// Create mock server that supports tool calling
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Accept both /v1/chat/completions and /chat/completions paths
		if r.URL.Path != "/v1/chat/completions" && r.URL.Path != "/chat/completions" {
			http.NotFound(w, r)
			return
		}

		// Parse request
		var req openai.ChatCompletionRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		// Check if this is the initial request or tool response
		lastMessage := req.Messages[len(req.Messages)-1]

		var resp openai.ChatCompletionResponse

		if lastMessage.Role == "user" {
			// Initial request - return with tool call
			resp = openai.ChatCompletionResponse{
				ID:      "chatcmpl-test",
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   req.Model,
				Choices: []openai.ChatCompletionChoice{
					{
						Index: 0,
						Message: openai.ChatCompletionMessage{
							Role:    openai.ChatMessageRoleAssistant,
							Content: "I'll use the calculator tool",
							ToolCalls: []openai.ToolCall{
								{
									ID:   "call_1",
									Type: "function",
									Function: openai.FunctionCall{
										Name:      "calculator",
										Arguments: `{"input":"2+2"}`,
									},
								},
							},
						},
						FinishReason: "tool_calls",
					},
				},
			}
		} else {
			// Tool response - return final answer
			resp = openai.ChatCompletionResponse{
				ID:      "chatcmpl-test-2",
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   req.Model,
				Choices: []openai.ChatCompletionChoice{
					{
						Index: 0,
						Message: openai.ChatCompletionMessage{
							Role:    openai.ChatMessageRoleAssistant,
							Content: "The result is 4",
						},
						FinishReason: "stop",
					},
				},
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create provider with mock server URL
	// Note: go-openai client automatically adds /v1 prefix to BaseURL
	provider := lms.New(lms.WithAddress(server.URL))

	// Create mock calculator tool
	calcTool := createMockTool("calculator", "A calculator tool", "4")
	tools := agent.Tools{calcTool}

	ctx := context.Background()
	storage := &mockHistoryStorage{}

	var receivedMessages []agent.Message
	onMessage := func(msg agent.Message) {
		receivedMessages = append(receivedMessages, msg)
	}

	response, err := provider.Call(ctx, "test-model", tools, []agent.Content{{Type: agent.ContentTypeText, Content: "What is 2+2?"}}, storage, onMessage)
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, "The result is 4", response.Content)

	// Check that we got messages including tool calls
	assert.GreaterOrEqual(t, len(receivedMessages), 3) // User + Assistant + Tool + Assistant

	// Find tool-related messages
	foundToolCall := false
	foundToolResult := false
	for _, msg := range receivedMessages {
		if strings.Contains(msg.Extra, "calculator") {
			foundToolCall = true
		}
		if msg.Role == "tool" {
			foundToolResult = true
		}
	}

	assert.True(t, foundToolCall, "Should have found tool call in messages")
	assert.True(t, foundToolResult, "Should have found tool result in messages")
}

// Test error scenarios
func TestCallWithConnectionError(t *testing.T) {
	// Create provider with invalid URL
	provider := lms.New(lms.WithAddress("http://invalid-host:9999"))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	storage := &mockHistoryStorage{}
	tools := agent.Tools{}

	_, err := provider.Call(ctx, "test-model", tools, []agent.Content{{Type: agent.ContentTypeText, Content: "Hello"}}, storage, nil)
	assert.Error(t, err)
}

// Test context cancellation
func TestCallWithContextCancellation(t *testing.T) {
	// Create mock server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Delay longer than context timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create provider with mock server URL
	// Note: go-openai client automatically adds /v1 prefix to BaseURL
	provider := lms.New(lms.WithAddress(server.URL))

	// Short timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	storage := &mockHistoryStorage{}
	tools := agent.Tools{}

	_, err := provider.Call(ctx, "test-model", tools, []agent.Content{{Type: agent.ContentTypeText, Content: "Hello"}}, storage, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

// Test with empty tools
func TestCallWithEmptyTools(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := openai.ChatCompletionResponse{
			ID:      "chatcmpl-test",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   "test-model",
			Choices: []openai.ChatCompletionChoice{
				{
					Index: 0,
					Message: openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleAssistant,
						Content: "Hello without tools",
					},
					FinishReason: "stop",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create provider with mock server URL
	// Note: go-openai client automatically adds /v1 prefix to BaseURL
	provider := lms.New(lms.WithAddress(server.URL))

	ctx := context.Background()
	storage := &mockHistoryStorage{}
	tools := agent.Tools{} // Empty tools

	response, err := provider.Call(ctx, "test-model", tools, []agent.Content{{Type: agent.ContentTypeText, Content: "Hello"}}, storage, nil)
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, "Hello without tools", response.Content)
}
