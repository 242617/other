package sb

import "testing"

func Test_GameFlow(t *testing.T) {

	_, err := Create(-1)
	assumeError(t, err, ErrIncorrectFieldSize)

	_, err = Create(0)
	assumeError(t, err, ErrIncorrectFieldSize)

	game, err := Create(6)
	assumeError(t, game.Setup([]*Ship{NewShip(NewPoint(0, 0), NewPoint(6, 0))}), ErrOutOfBounds)
	assumeError(t, game.Setup([]*Ship{NewShip(NewPoint(5, 0), NewPoint(6, 1))}), ErrOutOfBounds)

	ships := []*Ship{
		NewShip(NewPoint(0, 0), NewPoint(0, 0)),
		NewShip(NewPoint(1, 0), NewPoint(1, 1)),
		NewShip(NewPoint(1, 3), NewPoint(2, 5)),
	}
	assumeNoError(t, game.Setup(ships))
	assumeState(t, game.State(), 0, 0, len(ships), 0)
	assumeShot(t, game.Shot, NewPoint(6, 0), nil, ErrOutOfBounds)

	assumeShot(t, game.Shot, NewPoint(0, 0), &info{Destroy: true, Knock: true}, nil)
	assumeState(t, game.State(), 1, 0, len(ships), 1)

	assumeShot(t, game.Shot, NewPoint(0, 0), nil, ErrRepeatedShot)
	assumeState(t, game.State(), 1, 0, len(ships), 1)

	assumeShot(t, game.Shot, NewPoint(0, 1), nil, nil)
	assumeState(t, game.State(), 1, 0, len(ships), 2)

	assumeShot(t, game.Shot, NewPoint(1, 0), &info{Knock: true}, nil)
	assumeState(t, game.State(), 1, 1, len(ships), 3)
	assumeShot(t, game.Shot, NewPoint(2, 0), nil, nil)
	assumeState(t, game.State(), 1, 1, len(ships), 4)
	assumeShot(t, game.Shot, NewPoint(1, 1), &info{Destroy: true, Knock: true}, nil)
	assumeState(t, game.State(), 2, 0, len(ships), 5)

	assumeShot(t, game.Shot, NewPoint(0, 1), nil, ErrRepeatedShot)
	assumeState(t, game.State(), 2, 0, len(ships), 5)
	assumeShot(t, game.Shot, NewPoint(2, 1), nil, nil)
	assumeState(t, game.State(), 2, 0, len(ships), 6)

	shots := game.State().ShotCount
	for _, point := range []Point{
		NewPoint(1, 3),
		NewPoint(2, 3),
		NewPoint(1, 4),
		NewPoint(2, 4),
		NewPoint(1, 5),
	} {
		assumeShot(t, game.Shot, point, &info{Knock: true}, nil)
		shots++
		assumeState(t, game.State(), 2, 1, len(ships), shots)
	}
	assumeShot(t, game.Shot, NewPoint(2, 5), &info{Destroy: true, Knock: true, End: true}, nil)
	shots++
	assumeState(t, game.State(), 3, 0, len(ships), shots)

	assumeShot(t, game.Shot, NewPoint(0, 0), nil, ErrGameEnded)
	assumeState(t, game.State(), 3, 0, len(ships), shots)

}

func assumeShot(t *testing.T, shot func(Point) (*info, error), point Point, wantInfo *info, wantErr error) {
	gotInfo, gotErr := shot(point)
	if gotErr != wantErr {
		t.Fatalf(`unexpected error: want "%v", got: "%v"`, gotErr, wantErr)
	}
	if gotInfo != nil && wantInfo != nil && gotInfo.Destroy != wantInfo.Destroy {
		t.Fatalf(`unexpected destroy info: want "%t", got: "%t"`, gotInfo.Destroy, wantInfo.Destroy)
	}
	if gotInfo != nil && wantInfo != nil && gotInfo.Knock != wantInfo.Knock {
		t.Fatalf(`unexpected knock info: want "%t", got: "%t"`, gotInfo.Knock, wantInfo.Knock)
	}
	if gotInfo != nil && wantInfo != nil && gotInfo.End != wantInfo.End {
		t.Fatalf(`unexpected end info: want "%t", got: "%t"`, gotInfo.End, wantInfo.End)
	}
}

func assumeState(t *testing.T, state *state, destroyed, knocked, shipCount, shotCount int) {
	if state.Destroyed != destroyed {
		t.Fatalf(`incorrect calculating destroyed ships: want: "%d", got: "%d"`, destroyed, state.Destroyed)
	}
	if state.Knocked != knocked {
		t.Fatalf(`incorrect calculating knocked ships: want: "%d", got: "%d"`, knocked, state.Knocked)
	}
	if state.ShipCount != shipCount {
		t.Fatalf(`incorrect calculating ships: want: "%d", got: "%d"`, shipCount, state.ShipCount)
	}
	if state.ShotCount != shotCount {
		t.Fatalf(`incorrect calculating shots: want: "%d", got: "%d"`, shotCount, state.ShotCount)
	}
}

func assumeNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf(`unexpected not nil error: "%v"`, err)
	}
}

func assumeError(t *testing.T, got, want error) {
	if got != want {
		t.Fatalf(`unexpected error: want: "%v", got: "%v"`, want, got)
	}
}
