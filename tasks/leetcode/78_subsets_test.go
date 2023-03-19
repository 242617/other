package leetcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/242617/other/tasks/leetcode"
)

func TestSubsets(t *testing.T) {
	for _, data := range []struct {
		Input  []int
		Output [][]int
	}{
		{
			Input:  []int{1, 2, 3},
			Output: [][]int{{}, {1}, {2}, {1, 2}, {3}, {1, 3}, {2, 3}, {1, 2, 3}},
		},
		{
			Input:  []int{0},
			Output: [][]int{{}, {0}},
		},
		{
			Input:  []int{},
			Output: [][]int{},
		},
		{
			Input:  nil,
			Output: nil,
		},
	} {
		assert.Equal(t, data.Output, leetcode.Subsets(data.Input), data)
	}
}
