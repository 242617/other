package sb

import "testing"

func Test_Ship(t *testing.T) {
	s := NewShip(Point{3, 2}, Point{2, 1})
	c := []Point{{2, 1}, {2, 2}, {3, 1}, {3, 2}}
	if !equal(s.Coordinates(), c) {
		t.Fatalf(`unexpected ship coordinates: want "%v", got "%v"`, c, s.Coordinates())
	}

	ok := s.Knock(NewPoint(2, 1))
	if !ok {
		t.Fatal("expected successful knocking ship")
	}
	if !s.Knocked() {
		t.Fatal("expected ship to be knocked")
	}

	if !s.Knock(NewPoint(2, 2)) {
		t.Fatal("expected successful knocking ship")
	}
	if !s.Knock(NewPoint(3, 1)) {
		t.Fatal("expected successful knocking ship")
	}
	if !s.Knock(NewPoint(3, 2)) {
		t.Fatal("expected successful knocking ship")
	}
	if !s.Destroyed() {
		t.Fatal("expected ship to be knocked and destroyed")
	}
	if s.Knocked() {
		t.Fatal("expected ship to be not knocked after destroying")
	}
}

func Test_equal(t *testing.T) {
	var (
		p1     = NewPoint(0, 0)
		p2     = NewPoint(0, 1)
		p3     = NewPoint(1, 0)
		p4     = NewPoint(1, 1)
		t1, t2 = []Point{p1, p1, p2, p4}, []Point{p2, p2, p1, p4}
		t3, t4 = []Point{p1, p3}, []Point{p2, p3}
		t5, t6 = []Point{p1, p3, p4}, []Point{p2, p1, p3}
	)
	if !equal(t1, t2) {
		t.Fatalf(`expecting equal point slices: want: "%v", got: "%v"`, t1, t2)
	}
	if equal(t3, t4) || equal(t4, t3) {
		t.Fatalf(`expecting not equal point slices: want: "%v", got: "%v"`, t3, t4)
	}
	if equal(t5, t6) || equal(t6, t5) {
		t.Fatalf(`expecting not equal point slices: want: "%v", got: "%v"`, t5, t6)
	}
}

func equal(a, b []Point) bool {
	if len(b) != len(a) {
		return false
	}
	t := map[Point]bool{}
	for _, point := range a {
		t[point] = false
	}
	for _, point := range b {
		t[point] = true
	}
	for _, v := range t {
		if !v {
			return false
		}
	}
	return true
}
