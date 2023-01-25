package linkedlist

import (
	"fmt"
	"strings"
)

func NewNode[T any](value T) ListNode[T] { return ListNode[T]{Value: value} }

type ListNode[T any] struct {
	next  *ListNode[T]
	Value T
}

func (n *ListNode[T]) Next() *ListNode[T]    { return n.next }
func (n *ListNode[T]) Add(node *ListNode[T]) { n.next = node }

func (n *ListNode[T]) String() string {
	next := "<nil>"
	if n.next != nil {
		next = fmt.Sprintf("%v", n.next.Value)
	}
	return fmt.Sprintf("%v [ %s ]", n.next.Value, next)
}

func New[T any](head *ListNode[T]) *LinkedList[T] { return &LinkedList[T]{head: head} }

type LinkedList[T any] struct {
	head *ListNode[T]
}

func (l *LinkedList[T]) Head() *ListNode[T] { return l.head }
func (l *LinkedList[T]) Length() int {
	var length int
	for current := l.head; current != nil; current = current.next {
		length++
	}
	return length
}

func (l *LinkedList[T]) String() string {
	values := make([]string, 0, l.Length())
	for current := l.head; current != nil; current = current.next {
		values = append(values, fmt.Sprintf("%v", current.Value))
	}
	return strings.Join(values, "->")
}

func (l *LinkedList[T]) Get(n int) *ListNode[T] {
	if n > l.Length() {
		return nil
	}
	var node *ListNode[T]
	var i int
	for current := l.head; current != nil; current = current.next {
		if i == n {
			node = current
			break
		}
		i++
	}
	return node
}

func (l *LinkedList[T]) Remove(n int) {
	if n > l.Length() {
		return
	}

	var previous *ListNode[T]
	var i int
	for current := l.head; current != nil; current = current.next {
		if i == n {
			previous.next = current.next
			break
		}
		i++
		previous = current
	}
}

func (l *LinkedList[T]) Append(node *ListNode[T]) {
	if l.head == nil {
		l.head = node
		return
	}
	for current := l.head; current != nil; current = current.next {
		if current.next == nil {
			current.next = node
			break
		}
	}
}

func (l *LinkedList[T]) Reverse() *LinkedList[T] {
	if l.head == nil {
		return &LinkedList[T]{}
	}

	var res LinkedList[T]

	var previous *ListNode[T]
	for current := l.head; current != nil; {
		next := current.next
		current.next = previous
		if next == nil {
			res.head = current
			break
		}
		previous = current
		current = next
	}

	return &res
}
