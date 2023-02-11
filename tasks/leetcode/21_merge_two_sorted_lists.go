package leetcode

import (
	"strconv"
	"strings"
)

func MergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	switch {
	case list1 == nil:
		return list2
	case list2 == nil:
		return list1
	}

	var res []int
	for {
		if list1 == nil {
			for current := list2; current != nil; current = current.Next {
				res = append(res, current.Val)
			}
			break
		}
		if list2 == nil {
			for current := list1; current != nil; current = current.Next {
				res = append(res, current.Val)
			}
			break
		}
		if list1.Val == list2.Val {
			res = append(res, list1.Val, list2.Val)
			list1, list2 = list1.Next, list2.Next
			continue
		}
		for list1 != nil && list2 != nil && list1.Val < list2.Val {
			res = append(res, list1.Val)
			list1 = list1.Next
		}
		for list1 != nil && list2 != nil && list2.Val < list1.Val {
			res = append(res, list2.Val)
			list2 = list2.Next
		}
	}

	var node *ListNode
	for i := len(res) - 1; i >= 0; i-- {
		n := &ListNode{
			Val:  res[i],
			Next: node,
		}
		node = n
	}
	return node
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func IntSliceToListNode(a []int) *ListNode {
	var node *ListNode
	for i := len(a) - 1; i >= 0; i-- {
		n := &ListNode{
			Val:  a[i],
			Next: node,
		}
		node = n
	}
	return node
}

func ListNodeToIntSlice(l *ListNode) []int {
	if l == nil {
		return []int{}
	}
	a := []int{l.Val}
	for current := l.Next; current != nil; current = current.Next {
		a = append(a, current.Val)
	}
	return a
}

func (n *ListNode) String() string {
	nodes := []string{}
	for current := n; current != nil; current = current.Next {
		nodes = append(nodes, strconv.Itoa(current.Val))
	}
	return strings.Join(nodes, "->")
}
