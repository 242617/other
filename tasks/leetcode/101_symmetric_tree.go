package leetcode

func IsSymmetric(root *TreeNode) bool {
	a, b := make([]*int, 0, 1000), make([]*int, 0, 1000)

	var preOrder func(*TreeNode)
	preOrder = func(n *TreeNode) {
		if n == nil {
			a = append(a, nil)
			return
		}
		a = append(a, &n.Val)
		preOrder(n.Left)
		preOrder(n.Right)
	}
	var postOrder func(*TreeNode)
	postOrder = func(n *TreeNode) {
		if n == nil {
			b = append(b, nil)
			return
		}
		postOrder(n.Left)
		postOrder(n.Right)
		b = append(b, &n.Val)
	}

	preOrder(root)
	postOrder(root)

	for i, j := 0, len(b)-1; i < len(a) && j >= 0; i, j = i+1, j-1 {
		if a[i] != b[j] {
			if (a[i] == nil && b[j] != nil) || (b[j] == nil && a[i] != nil) || *a[i] != *b[j] {
				return false
			}
		}
	}

	return true
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}
