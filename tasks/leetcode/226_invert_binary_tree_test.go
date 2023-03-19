package leetcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/242617/other/tasks/leetcode"
)

func TestInvertTree(t *testing.T) {
	for _, data := range []struct {
		Root   *leetcode.TreeNode
		Result *leetcode.TreeNode
	}{
		{
			Root: &leetcode.TreeNode{Val: 4,
				Left: &leetcode.TreeNode{Val: 2,
					Left:  &leetcode.TreeNode{Val: 1},
					Right: &leetcode.TreeNode{Val: 3},
				},
				Right: &leetcode.TreeNode{Val: 7,
					Left:  &leetcode.TreeNode{Val: 6},
					Right: &leetcode.TreeNode{Val: 9},
				},
			},
			Result: &leetcode.TreeNode{Val: 4,
				Left: &leetcode.TreeNode{Val: 7,
					Left:  &leetcode.TreeNode{Val: 9},
					Right: &leetcode.TreeNode{Val: 6},
				},
				Right: &leetcode.TreeNode{Val: 2,
					Left:  &leetcode.TreeNode{Val: 3},
					Right: &leetcode.TreeNode{Val: 1},
				},
			},
		},
		{
			Root: &leetcode.TreeNode{Val: 2,
				Left:  &leetcode.TreeNode{Val: 1},
				Right: &leetcode.TreeNode{Val: 3},
			},
			Result: &leetcode.TreeNode{Val: 2,
				Left:  &leetcode.TreeNode{Val: 3},
				Right: &leetcode.TreeNode{Val: 1},
			},
		},
		{
			Root: &leetcode.TreeNode{Val: 2,
				Right: &leetcode.TreeNode{Val: 3},
			},
			Result: &leetcode.TreeNode{Val: 2,
				Left: &leetcode.TreeNode{Val: 3},
			},
		},
		{
			Root:   nil,
			Result: nil,
		},
	} {
		assert.Equal(t, data.Result, leetcode.InvertTree(data.Root), data)
	}
}
