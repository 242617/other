package leetcode_test

import (
	"testing"

	"github.com/242617/other/tasks/leetcode"
	"github.com/stretchr/testify/assert"
)

func TestRotate(t *testing.T) {
	for _, data := range []struct {
		Input  [][]int
		Output [][]int
	}{
		{
			Input:  [][]int{{1}},
			Output: [][]int{{1}},
		},
		{
			Input:  [][]int{{1, 2}, {3, 4}},
			Output: [][]int{{3, 1}, {4, 2}},
		},
		{
			Input:  [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			Output: [][]int{{7, 4, 1}, {8, 5, 2}, {9, 6, 3}},
		},
		{
			Input:  [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16}},
			Output: [][]int{{13, 9, 5, 1}, {14, 10, 6, 2}, {15, 11, 07, 3}, {16, 12, 8, 4}},
		},
		{
			Input:  [][]int{{5, 1, 9, 11}, {2, 4, 8, 10}, {13, 3, 6, 7}, {15, 14, 12, 16}},
			Output: [][]int{{15, 13, 2, 5}, {14, 3, 4, 1}, {12, 6, 8, 9}, {16, 7, 10, 11}},
		},
		{
			Input:  [][]int{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 10}, {11, 12, 13, 14, 15}, {16, 17, 18, 19, 20}, {21, 22, 23, 24, 25}},
			Output: [][]int{{21, 16, 11, 6, 1}, {22, 17, 12, 7, 2}, {23, 18, 13, 8, 3}, {24, 19, 14, 9, 4}, {25, 20, 15, 10, 5}},
		},
	} {
		leetcode.Rotate(data.Input)
		assert.Equal(t, data.Output, data.Input, data)
	}
}
