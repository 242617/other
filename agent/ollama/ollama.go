package ollama

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ollama/ollama/api"

	"github.com/242617/other/agent"
)

type (
	encodeFunc  = func(message api.Message) (string, error)
	dencodeFunc = func(item string) (api.Message, error)

	responseFunc = func(res api.ChatResponse) error
)

func New(host string) *Ollama {
	return &Ollama{
		encode: encode,
		decode: decode,
		client: api.NewClient(&url.URL{Scheme: "http", Host: host}, &http.Client{}),
		name:   "ollama assistant",
	}
}

type Ollama struct {
	encode encodeFunc
	decode dencodeFunc
	client *api.Client
	name   string
}

func toMessage(m api.Message) agent.Message {
	var extra string
	if len(m.ToolCalls) > 0 {
		tools := make([]string, len(m.ToolCalls))
		for i, toolCall := range m.ToolCalls {
			tools[i] = toolCall.Function.Name + toolCall.Function.Arguments.String()
		}
		extra += fmt.Sprintf("tools: %q", strings.Join(tools, ", "))
	}

	if m.ToolName != "" {
		extra += fmt.Sprintf("tool name: %q", m.ToolName)
	}

	return agent.Message{
		Role:  m.Role,
		Text:  m.Content,
		Extra: extra,
	}
}

func toAPIProperties(properties map[string]agent.ToolInfoFunctionParametersProperty) map[string]api.ToolProperty {
	res := make(map[string]api.ToolProperty, len(properties))
	for k, v := range properties {
		res[k] = api.ToolProperty{
			Type:        api.PropertyType{v.Type},
			Items:       v.Items,
			Description: v.Description,
		}
	}
	return res
}

func toTool(tool agent.ToolInfo) api.Tool {
	return api.Tool{
		Type: tool.Type,
		Function: api.ToolFunction{
			Name:        tool.Function.Name,
			Description: tool.Function.Description,
			Parameters: api.ToolFunctionParameters{
				Type:       tool.Function.Parameters.Type,
				Properties: toAPIProperties(tool.Function.Parameters.Properties),
				Required:   tool.Function.Parameters.Required,
			},
		},
	}
}
func toTools(tools []agent.ToolInfo) []api.Tool {
	res := make([]api.Tool, len(tools))
	for i, tool := range tools {
		res[i] = toTool(tool)
	}
	return res
}

func (p *Ollama) EncodeSystemMessage(msg agent.Message) (string, error) {
	return p.encode(api.Message{Role: msg.Role, Content: msg.Text})
}
