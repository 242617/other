package leetcode

func InvertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return root
	}
	if root.Left != nil {
		root.Left = InvertTree(root.Left)
	}
	if root.Right != nil {
		root.Right = InvertTree(root.Right)
	}
	root.Left, root.Right = root.Right, root.Left
	return root
}
