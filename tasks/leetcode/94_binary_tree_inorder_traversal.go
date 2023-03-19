package leetcode

func InorderTraversal(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	var dive func(root *TreeNode) []int
	dive = func(root *TreeNode) []int {
		res := []int{}
		if root.Left != nil {
			res = append(res, dive(root.Left)...)
		}
		res = append(res, root.Val)
		if root.Right != nil {
			res = append(res, dive(root.Right)...)
		}
		return res
	}
	return dive(root)
}
