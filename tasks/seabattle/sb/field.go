package sb

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	ErrIncorrectFieldSize = errors.New("incorrect field size")
	ErrPointIsBusy        = errors.New("point is busy")
	ErrOutOfBounds        = errors.New("out of bounds")
)

type field [][]*cell

func NewField(width, height int) (*field, error) {
	if width < 1 || height < 1 {
		return nil, ErrIncorrectFieldSize
	}

	var f field = make([][]*cell, width)
	for x := 0; x < width; x++ {
		col := make([]*cell, height)
		for y := 0; y < height; y++ {
			col[y] = &cell{x: x, y: y}
		}
		f[x] = col
	}

	return &f, nil
}

func (f *field) check(p Point) error {
	if p.X < 0 || p.X > len(*f)-1 ||
		p.Y < 0 || p.Y > len((*f)[0])-1 {
		return ErrOutOfBounds
	}
	return nil
}

func (f *field) cell(p Point) (*cell, error) {
	if err := f.check(p); err != nil {
		return nil, err
	}
	return (*f)[p.X][p.Y], nil
}

func (f *field) add(s *Ship) error {
	for _, p := range s.Coordinates() {
		if err := f.check(p); err != nil {
			return err
		}
		if (*f)[p.X][p.Y].ship != nil {
			return ErrPointIsBusy
		}
	}
	for _, p := range s.Coordinates() {
		(*f)[p.X][p.Y].ship = s
	}
	return nil
}

func (f *field) String() string {
	buf := bytes.NewBufferString("\n")
	for x := 0; x < len(*f); x++ {
		for y := 0; y < len((*f)[0]); y++ {
			hit := (*f)[x][y].hit
			ship := (*f)[x][y].ship != nil
			item := btos(ship, hit)
			buf.WriteString(fmt.Sprintf(` %s `, item))
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

func btos(ship bool, hit bool) string {
	if ship {
		if hit {
			return "*"
		}
		return "+"
	}
	return "."
}
