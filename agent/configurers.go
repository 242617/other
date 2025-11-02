package agent

import (
	"strings"

	"github.com/242617/other/agent/inmemory_storage"
)

type Configurer interface{ Modifier | Option }

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
			`<instructions>`,
			`- ALWAYS follow <self_reflection> and <answering_rules>.`,
			`- You are a helpful assistant.`,
			`- Do not ask user to do your task.`,
			`- Always follow chain of thoughts <chain_of_thought>.`,

			`<self_reflection>`,
			`1. Spend time thinking of a rubric, from a role POV, until you are confident.`,
			`2. Keep going until solved.`,
			`3. Split each task for subtasks. Modify files iteratively, one or two changes at a time.`,
			`4. After each task check your actions by comparing them with the scheduled plan. What was scheduled and what was completed?`,
			`5. Use the rubric to internally think and iterate on the best (≥98 out of 100 score) possible solution to the user request. IF your response is not hitting the top marks acrosall categories in the rubric, you need to start again",`,
			`</self_reflection>`,

			`<answering_rules>`,
			`1. USE the language of USER message.`,
			`2. Act as a role assigned.`,
			`3. Answer the question in a natural, human-like manner.`,
			`4. If not requested by the user, no actionable items are needed by default.`,
			`5. Respect user's current time and date, his location.`,
			`6. Answer briefly and concisely, 4-5 sentences at most.`,
			`7. Think hard.`,
			`</answering_rules>`,

			`<chain_of_thought>`,
			`- Understand request (and timing of user action) → PLAN explore (read-only) → propose collaborative plan with options/risks/tests → ask if it matches → output/act`,
			`- When all steps succeed and are confirmed, mark successul task completion with ✅.`,
			`</chain_of_thought>`,

			`</instructions>`,
		}, "\n"),
	)
}
func WithSystem(system ...string) Modifier {
	return func(a *Agent) { a.system = strings.Join(system, "\n") }
}

func withDefaultOnMessageFunc() Modifier { return WithOnMessageFunc(func(Message) {}) }
func WithOnMessageFunc(fn MessageCallback) Modifier {
	return func(a *Agent) { a.messageCallback = fn }
}

type Option = func(*Agent) error

func withDefaultStorage() Modifier { return WithStorage(inmemory_storage.New(1, 100)) }
func WithStorage(storage HistoryStorage) Modifier {
	return func(a *Agent) { a.storage = storage }
}
