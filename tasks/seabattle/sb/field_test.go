package sb

import "testing"

func Test_Field(t *testing.T) {
	_, err := NewField(-1, 0)
	assumeError(t, err, ErrIncorrectFieldSize)

	f, err := NewField(5, 5)
	assumeNoError(t, err)

	_, err = f.cell(NewPoint(-1, 0))
	assumeError(t, err, ErrOutOfBounds)

	c, err := f.cell(NewPoint(1, 1))
	assumeNoError(t, err)
	if c.ship != nil {
		t.Fatalf(`unexpected ship: %v`, c.ship)
	}

	err = f.add(NewShip(NewPoint(1, 1), NewPoint(2, 2)))
	assumeNoError(t, err)
	if c.ship == nil {
		t.Fatalf(`expecting ship here: %v`, c.ship)
	}
}
