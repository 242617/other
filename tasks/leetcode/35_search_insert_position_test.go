package leetcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/242617/other/tasks/leetcode"
)

func TestSearchInsert(t *testing.T) {
	for _, data := range []struct {
		Nums   []int
		Target int
		Result int
	}{
		{
			Nums:   []int{1, 3, 5, 6},
			Target: 5,
			Result: 2,
		},
		{
			Nums:   []int{1, 3, 5, 6},
			Target: 2,
			Result: 1,
		},
		{
			Nums:   []int{1, 3, 5, 6},
			Target: 7,
			Result: 4,
		},
		{
			Nums:   []int{1, 3, 5, 6},
			Target: 0,
			Result: 0,
		},
		{
			Nums:   []int{2, 5},
			Target: 1,
			Result: 0,
		},
		{
			Nums:   []int{-1, 3, 5, 6},
			Target: 0,
			Result: 1,
		},
		{
			Nums:   []int{-6, -5, -3, -1},
			Target: -2,
			Result: 3,
		},
	} {
		assert.Equal(t, data.Result, leetcode.SearchInsert(data.Nums, data.Target), data)
	}
}
