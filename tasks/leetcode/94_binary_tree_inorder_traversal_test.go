package leetcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/242617/other/tasks/leetcode"
)

func TestInorderTraversal(t *testing.T) {
	for _, data := range []struct {
		Root   *leetcode.TreeNode
		Output []int
	}{
		{
			Root: &leetcode.TreeNode{Val: 1,
				Right: &leetcode.TreeNode{Val: 2,
					Left: &leetcode.TreeNode{Val: 3},
				},
			},
			Output: []int{1, 3, 2},
		},
	} {
		assert.Equal(t, data.Output, leetcode.InorderTraversal(data.Root), data)
	}
}
