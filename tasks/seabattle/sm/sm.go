package sm

type State string

const (
	StateInit  State = "init"
	StateSetup State = "setup"
	StateReady State = "ready"
)

type Action string

const (
	ActionCreateMatrix Action = "create_matrix"
	ActionShip         Action = "ship"
	ActionShot         Action = "shot"
	ActionClear        Action = "clear"
	ActionState        Action = "state"
)

type path struct {
	target State
	change bool
}

type StateMachine struct {
	current State
}

func New() *StateMachine {
	return &StateMachine{current: StateInit}
}

func (sm *StateMachine) Current() State {
	return sm.current
}

func (sm *StateMachine) Action(action Action) bool {
	p, ok := paths[sm.current][action]
	if !ok {
		return false
	}
	if p.change && sm.current != p.target {
		sm.current = p.target
	}
	return true
}
