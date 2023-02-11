package leetcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/242617/other/tasks/leetcode"
)

func TestMergeTwoLists(t *testing.T) {
	for _, data := range []struct {
		List1, List2 []int
		Result       []int
	}{
		{
			List1:  []int{1, 2, 4},
			List2:  []int{1, 3, 4, 5},
			Result: []int{1, 1, 2, 3, 4, 4, 5},
		},
		{
			List1:  []int{},
			List2:  []int{},
			Result: []int{},
		},
		{
			List1:  []int{},
			List2:  []int{0},
			Result: []int{0},
		},
		{
			List1:  []int{2},
			List2:  []int{1},
			Result: []int{1, 2},
		},
		{
			List1:  []int{1},
			List2:  []int{2},
			Result: []int{1, 2},
		},
	} {
		assert.Equal(t, data.Result,
			leetcode.ListNodeToIntSlice(
				leetcode.MergeTwoLists(
					leetcode.IntSliceToListNode(data.List1), leetcode.IntSliceToListNode(data.List2),
				),
			),
			data,
		)
	}
}
