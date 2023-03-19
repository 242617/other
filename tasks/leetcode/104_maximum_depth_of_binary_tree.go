package leetcode

func MaxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	var max int
	var dive func(*TreeNode, int)
	dive = func(node *TreeNode, level int) {
		if level > max {
			max = level
		}
		if node.Left != nil {
			dive(node.Left, level+1)
		}
		if node.Right != nil {
			dive(node.Right, level+1)
		}
	}
	dive(root, 1)
	return max
}
