package sb

import (
	"bytes"
	"fmt"
	"math"
)

type Ship struct {
	coords map[Point]bool
}

func NewShip(begin, end Point) *Ship {
	minX := int(math.Min(float64(begin.X), float64(end.X)))
	maxX := int(math.Max(float64(begin.X), float64(end.X)))
	minY := int(math.Min(float64(begin.Y), float64(end.Y)))
	maxY := int(math.Max(float64(begin.Y), float64(end.Y)))
	s := Ship{
		coords: map[Point]bool{},
	}
	var n int
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			s.coords[NewPoint(x, y)] = false
			n++
		}
	}
	return &s
}

func (s *Ship) Knock(p Point) bool {
	_, ok := s.coords[p]
	if !ok {
		return false
	}
	s.coords[p] = true
	return true
}

func (s *Ship) Coordinates() []Point {
	res := make([]Point, len(s.coords))
	var n int
	for point := range s.coords {
		res[n] = point
		n++
	}
	return res
}

func (s *Ship) Knocked() bool {
	if s.Destroyed() {
		return false
	}
	return s.wounds() > 0
}

func (s *Ship) Destroyed() bool {
	return len(s.Coordinates()) == s.wounds()
}

func (s *Ship) wounds() int {
	var wounds int
	for _, hit := range s.coords {
		if hit {
			wounds++
		}
	}
	return wounds
}

func (s *Ship) String() string {
	buf := bytes.NewBufferString("\n")
	buf.WriteString(fmt.Sprintf(`%v`, s.Coordinates()))
	buf.WriteString("\n")
	return buf.String()
}
