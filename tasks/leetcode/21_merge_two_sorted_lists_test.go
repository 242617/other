package leetcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/242617/other/tasks/leetcode"
)

func TestMergeTwoLists(t *testing.T) {
	for _, data := range []struct {
		List1, List2 []int
		Output       []int
	}{
		{
			List1:  []int{1, 2, 4},
			List2:  []int{1, 3, 4, 5},
			Output: []int{1, 1, 2, 3, 4, 4, 5},
		},
		{
			List1:  []int{},
			List2:  []int{},
			Output: []int{},
		},
		{
			List1:  []int{},
			List2:  []int{0},
			Output: []int{0},
		},
		{
			List1:  []int{2},
			List2:  []int{1},
			Output: []int{1, 2},
		},
		{
			List1:  []int{1},
			List2:  []int{2},
			Output: []int{1, 2},
		},
	} {
		assert.Equal(t, data.Output,
			leetcode.ListNodeToIntSlice(
				leetcode.MergeTwoLists(
					leetcode.IntSliceToListNode(data.List1), leetcode.IntSliceToListNode(data.List2),
				),
			),
			data,
		)
	}
}
