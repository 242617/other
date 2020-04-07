package sm

var paths = map[State]map[Action]path{
	StateInit: {
		ActionCreateMatrix: {
			target: StateSetup,
			change: true,
		},
	},

	StateSetup: {
		ActionShip: {
			target: StateReady,
			change: true,
		},
	},

	StateReady: {
		ActionShot: {
			change: false,
		},
		ActionState: {
			change: false,
		},
		ActionClear: {
			target: StateInit,
			change: true,
		},
	},
}
