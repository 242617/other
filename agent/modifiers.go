package agent

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ollama/ollama/api"

	"github.com/242617/other/tools"
)

type Modifier = func(*Agent)

func withDefaultClient() Modifier {
	return WithClient(api.NewClient(&url.URL{Scheme: "http", Host: "localhost:11434"}, &http.Client{}))
}
func WithClientFromAddress(scheme, address string) Modifier {
	return WithClient(api.NewClient(&url.URL{Scheme: scheme, Host: address}, &http.Client{}))
}
func WithClient(client *api.Client) Modifier {
	return func(a *Agent) { a.client = client }
}

func withDefaultModel() Modifier      { return WithModel("llama3.2:3b") }
func WithModel(model string) Modifier { return func(a *Agent) { a.model = model } }

func WithTools(tools ...tools.Tool) Modifier {
	return func(a *Agent) { a.tools = append(a.tools, tools...) }
}

func withDefaultSystem() Modifier {
	return WithSystem(
		strings.Join([]string{
			"You are a helpful assistant that answers questions.",
			"You will be given a series of messages from the user in order to generate a response.",
			"Answer in natural language.",
			"Be super straight. Do not be verbose.",
			"You should use 24-hour format for time.",
		}, "\n"),
	)
}
func WithSystem(system string) Modifier { return func(a *Agent) { a.system = system } }

func withDefaultOptions() Modifier {
	return WithOptions(map[string]any{
		"temperature":   0.0,
		"repeat_last_n": 2,
	})
}

func WithOptions(options map[string]any) Modifier {
	return func(a *Agent) { a.options = options }
}
