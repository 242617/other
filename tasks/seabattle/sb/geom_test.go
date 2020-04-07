package sb

import "testing"

func Test_Parse(t *testing.T) {
	if _, err := ParsePoint(""); err != ErrInvalidStringPointRepresentation {
		t.Fatal("no error on invalid string point representation")
	}

	if _, err := ParsePoint("29Ъ"); err != ErrIncorrectStringPointRepresentation {
		t.Fatal("no error on incorrect string point representation")
	}

	for k, v := range map[string]struct{ x, y int }{
		"1A": {0, 0},
		"3D": {3, 2},
	} {

		p, err := ParsePoint(k)
		if err != nil {
			t.Fatalf("cannot parse point: %s", err)
		}
		if p.X != v.x || p.Y != v.y {
			t.Fatalf("invalid point: want { %d, %d }, got: { %d, %d }", v.x, v.y, p.X, p.Y)
		}

	}
}
