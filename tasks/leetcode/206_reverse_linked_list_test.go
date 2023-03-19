package leetcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/242617/other/tasks/leetcode"
)

func TestReverseList(t *testing.T) {
	for _, data := range []struct {
		Input  []int
		Output []int
	}{
		{
			Input:  []int{1, 2, 3, 4, 5},
			Output: []int{5, 4, 3, 2, 1},
		},
		{
			Input:  []int{1, 2},
			Output: []int{2, 1},
		},
		{
			Input:  []int{},
			Output: []int{},
		},
	} {
		assert.Equal(t, data.Output,
			leetcode.ListNodeToIntSlice(
				leetcode.ReverseList(
					leetcode.IntSliceToListNode(data.Input),
				),
			),
			data,
		)
	}
}
