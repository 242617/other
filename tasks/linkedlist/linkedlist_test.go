package linkedlist_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/242617/other/tasks/linkedlist"
)

func TestBasic(t *testing.T) {
	first := linkedlist.NewNode("first")
	second := linkedlist.NewNode("second")
	third := linkedlist.NewNode("third")

	list := new(linkedlist.LinkedList[string])
	assert.Equal(t, 0, list.Length(), "list must have length zero after creating")

	list.Append(&first)
	assert.Equal(t, 1, list.Length(), "list must have length one with one node")
	assert.Equal(t, &first, list.Head(), `first element must be "first"`)

	first.Add(&second)
	second.Add(&third)
	assert.Equal(t, 3, list.Length(), "list must have length three with three nodes")

	assert.Equal(t, &second, first.Next(), `"first" must have "second" as next`)
	assert.Equal(t, &third, second.Next(), `"second" must have "third" as next`)

	assert.Equal(t, 3, list.Length(), "length after triple adds must be three")
	assert.Equal(t, "first->second->third", list.String(), "unexpected list values")

	assert.Equal(t, &second, list.Get(1), `second node must be "second"`)
	assert.Equal(t, &third, list.Get(2), `third node must be "third"`)
	assert.Nil(t, list.Get(3), `fourth node must be nil`)

	list.Remove(1)
	assert.Equal(t, 2, list.Length(), "list length must be two after removing node")
	assert.Equal(t, &first, list.Head(), `head must be "first" after removing second node`)
	assert.Equal(t, &third, first.Next(), `next for first must be "third"`)

	list.Remove(1)
	assert.Equal(t, 1, list.Length(), "list length must be one after removing the last node")
	assert.Equal(t, &first, list.Head(), `head must be "first" after removing second node`)
	assert.Nil(t, first.Next(), `next for first must be nil`)
}

func TestReverse(t *testing.T) {
	first := linkedlist.NewNode("first")
	second := linkedlist.NewNode("second")
	third := linkedlist.NewNode("third")

	list := new(linkedlist.LinkedList[string])
	reversed := list.Reverse()
	assert.Equal(t, 0, reversed.Length(), "reversed list must have zero length after creating")

	list.Append(&first)
	list.Append(&second)
	list.Append(&third)

	reversed = list.Reverse()
	assert.Equal(t, 3, reversed.Length(), "length after reverse must be three")
	assert.Equal(t, "third->second->first", reversed.String(), "unexpected reversed list values")
	reversed = reversed.Reverse()
	assert.Equal(t, "first->second->third", reversed.String(), "unexpected doubly reversed list values")
}
