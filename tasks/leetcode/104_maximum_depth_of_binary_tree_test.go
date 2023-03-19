package leetcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/242617/other/tasks/leetcode"
)

func TestMaxDepth(t *testing.T) {
	for _, data := range []struct {
		Root   *leetcode.TreeNode
		Result int
	}{
		{
			Root: &leetcode.TreeNode{Val: 3,
				Left: &leetcode.TreeNode{Val: 9},
				Right: &leetcode.TreeNode{Val: 20,
					Left:  &leetcode.TreeNode{Val: 15},
					Right: &leetcode.TreeNode{Val: 7},
				},
			},
			Result: 3,
		},
		{
			Root: &leetcode.TreeNode{Val: 1,
				Right: &leetcode.TreeNode{Val: 2},
			},
			Result: 2,
		},
		{
			Root:   &leetcode.TreeNode{Val: 1},
			Result: 1,
		},
		{
			Root:   nil,
			Result: 0,
		},
	} {
		assert.Equal(t, data.Result, leetcode.MaxDepth(data.Root), data)
	}
}
