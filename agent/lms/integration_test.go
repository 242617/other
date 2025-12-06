package lms_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/242617/other/agent"
	"github.com/242617/other/agent/inmemory_storage"
	"github.com/242617/other/agent/lms"
)

// TestLiveIntegration tests against a real LM Studio instance
// This test requires LM Studio to be running at http://127.0.0.1:1234 with openai/gpt-oss-20b model loaded
func TestLiveIntegration(t *testing.T) {
	// Skip if running in CI or without explicit flag
	if testing.Short() {
		t.Skip("Skipping live integration test")
	}

	// Create provider for live LM Studio instance
	provider := lms.New(lms.WithHost("http://127.0.0.1:1234"))
	storage := inmemory_storage.New(10, 100) // 10 pinned system messages, 100 shifting conversation messages
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test basic chat completion
	t.Run("basic_chat", func(t *testing.T) {
		var messages []agent.Message
		onMessage := func(msg agent.Message) {
			messages = append(messages, msg)
			t.Logf("Message [%s]: %s", msg.Role, msg.Text)
		}

		response, err := provider.Call(
			ctx,
			"openai/gpt-oss-20b",
			agent.Tools{},
			"Hello! Please respond with exactly 'Hello from LM Studio!' and nothing else.",
			storage,
			onMessage,
		)

		require.NoError(t, err)
		assert.NotEmpty(t, response)
		assert.Contains(t, response, "Hello")
		t.Logf("Final response: %s", response)

		// Verify we got both user and assistant messages
		assert.GreaterOrEqual(t, len(messages), 2)
	})

	// Test with a simple math tool
	t.Run("with_tool_calling", func(t *testing.T) {
		// Create a simple calculator tool
		calcTool := &mockTool{
			name: "calculator",
			info: agent.ToolInfo{
				Type: "function",
				Function: agent.ToolInfoFunction{
					Name:        "calculator",
					Description: "Performs basic arithmetic calculations",
					Parameters: agent.ToolInfoFunctionParameters{
						Type: "object",
						Properties: map[string]agent.ToolInfoFunctionParametersProperty{
							"expression": {
								Type:        "string",
								Description: "Mathematical expression to evaluate (e.g., '2+2', '10*5')",
							},
						},
						Required: []string{"expression"},
					},
				},
			},
			result: "42", // Mock result
		}

		tools := agent.Tools{calcTool}
		var messages []agent.Message
		onMessage := func(msg agent.Message) {
			messages = append(messages, msg)
			t.Logf("Message [%s]: %s | Extra: %s", msg.Role, msg.Text, msg.Extra)
		}

		response, err := provider.Call(
			ctx,
			"openai/gpt-oss-20b",
			tools,
			"Please use the calculator tool to compute 2+2. Just call the tool with expression '2+2'.",
			storage,
			onMessage,
		)

		require.NoError(t, err)
		assert.NotEmpty(t, response)
		t.Logf("Final response: %s", response)

		// Check if any message contains tool call information
		foundToolCall := false
		foundToolResult := false
		for _, msg := range messages {
			if msg.Extra != "" {
				t.Logf("Message with extra info: %s", msg.Extra)
				if msg.Role == "assistant" && len(msg.Extra) > 0 {
					foundToolCall = true
				}
				if msg.Role == "tool" {
					foundToolResult = true
				}
			}
		}

		t.Logf("Tool call found: %v, Tool result found: %v", foundToolCall, foundToolResult)
	})
}

// TestLiveSimpleChat is a minimal test to check basic connectivity
func TestLiveSimpleChat(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping live test")
	}

	provider := lms.New(lms.WithHost("http://127.0.0.1:1234"))
	storage := &mockHistoryStorage{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := provider.Call(
		ctx,
		"openai/gpt-oss-20b",
		agent.Tools{},
		"Say 'Hello world' and nothing else.",
		storage,
		func(agent.Message) {},
	)

	if err != nil {
		t.Logf("Live test failed (this is expected if LM Studio is not running): %v", err)
		t.Skip("LM Studio not available")
		return
	}

	require.NoError(t, err)
	assert.NotEmpty(t, response)
	t.Logf("LM Studio response: %s", response)
}
