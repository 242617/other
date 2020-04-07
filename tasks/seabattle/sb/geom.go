package sb

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidStringPointRepresentation   = errors.New("invalid string point representation")
	ErrIncorrectStringPointRepresentation = errors.New("incorrect string point representation")
)

const Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewPoint(x, y int) Point {
	return Point{X: x, Y: y}
}

type Point struct {
	X, Y int
}

func ParsePoint(str string) (*Point, error) {
	l := len([]rune(str))
	if l < 2 || l > 3 {
		return nil, ErrInvalidStringPointRepresentation
	}

	rawDigit, rawLetter := str[0:1], str[1:]

	digit, err := strconv.Atoi(rawDigit)
	if err != nil {
		return nil, err
	}
	digit--

	index := strings.Index(Alphabet, rawLetter)
	if index == -1 || digit > len(Alphabet) {
		return nil, ErrIncorrectStringPointRepresentation
	}

	return &Point{index, digit}, nil
}
