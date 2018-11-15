package search

import (
	"fmt"
	"strings"
	"testing"
)

func TestListShift(t *testing.T) {
	item := newItem(5)
	a, b, c, d := 1, 2, 3, 5
	item.value = &a
	item.child.value = &b
	item.child.child.value = &c
	item.shift(&d)
	item.child.shift(nil)

	if (item.value == nil || *item.value != 5) ||
		(item.child.value != nil) ||
		(item.child.child.value == nil || *item.child.child.value != 1) ||
		(item.child.child.child.value == nil || *item.child.child.child.value != 2) ||
		(item.child.child.child.child.value == nil || *item.child.child.child.child.value != 3) {
		t.Fail()
	}
}

func trace(l *item) string {
	res := []string{}
	str := fmt.Sprintf("[%d:", l.number)
	if l.value != nil {
		str += fmt.Sprintf("%d", *l.value)
	} else {
		str += "_"
	}
	str += "]"
	if l.child != nil {
		str += trace(l.child)
	}
	res = append(res, str)
	return strings.Join(res, ", ")
}
