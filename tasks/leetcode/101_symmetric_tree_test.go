package leetcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/242617/other/tasks/leetcode"
)

func TestIsSymmetric(t *testing.T) {
	for _, data := range []struct {
		Root   *leetcode.TreeNode
		Result bool
	}{
		{
			Root: &leetcode.TreeNode{Val: 1,
				Left: &leetcode.TreeNode{Val: 2,
					Left:  &leetcode.TreeNode{Val: 3},
					Right: &leetcode.TreeNode{Val: 4},
				},
				Right: &leetcode.TreeNode{Val: 2,
					Left:  &leetcode.TreeNode{Val: 4},
					Right: &leetcode.TreeNode{Val: 3},
				},
			},
			Result: true,
		},
		{
			Root: &leetcode.TreeNode{Val: 1,
				Left: &leetcode.TreeNode{Val: 2,
					Right: &leetcode.TreeNode{Val: 3},
				},
				Right: &leetcode.TreeNode{Val: 2,
					Right: &leetcode.TreeNode{Val: 3},
				},
			},
			Result: false,
		},
	} {
		assert.Equal(t, data.Result, leetcode.IsSymmetric(data.Root), data)
	}
}
