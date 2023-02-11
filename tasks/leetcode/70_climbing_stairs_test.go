package leetcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/242617/other/tasks/leetcode"
)

func TestClimbStairs(t *testing.T) {
	for _, data := range []struct {
		Input  int
		Output int
	}{
		{
			Input:  2,
			Output: 2,
		},
		{
			Input:  3,
			Output: 3,
		},
		{
			Input:  4,
			Output: 5,
		},
		{
			Input:  44,
			Output: 1134903170,
		},
	} {
		assert.Equal(t, data.Output, leetcode.ClimbStairs(data.Input), data)
	}
}
