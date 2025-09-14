package agent

import "strings"

type Modifier = func(*Agent)

func WithProvider(provider Provider) Modifier {
	return func(a *Agent) { a.provider = provider }
}

func withDefaultModel() Modifier      { return WithModel("llama3.2:3b") }
func WithModel(model string) Modifier { return func(a *Agent) { a.model = model } }

func WithTools(tools ...Tool) Modifier {
	return func(a *Agent) { a.tools = append(a.tools, tools...) }
}

func withDefaultOptions() Modifier {
	return WithOptions(map[string]any{
		"temperature":   0.0,
		"repeat_last_n": 2,
	})
}
func WithOptions(options map[string]any) Modifier {
	return func(a *Agent) { a.options = options }
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

func withDefaultOnMessageFunc() Modifier { return WithOnMessageFunc(func(Message) {}) }
func WithOnMessageFunc(fn MessageCallback) Modifier {
	return func(a *Agent) { a.onMessageFunc = fn }
}
