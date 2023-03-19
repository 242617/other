package leetcode

func ReverseList(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}

	current := head
	var previous *ListNode
	for {
		if current == nil {
			break
		}
		next := current.Next
		current.Next = previous
		previous = current
		current = next
	}

	return previous
}
