package leetcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/242617/other/tasks/leetcode"
)

func TestIsAlienSorted(t *testing.T) {
	for _, data := range []struct {
		Words  []string
		Order  string
		Result bool
	}{
		{
			Words:  []string{"hello", "leetcode"},
			Order:  "hlabcdefgijkmnopqrstuvwxyz",
			Result: true,
		},
		{
			Words:  []string{"word", "world", "row"},
			Order:  "worldabcefghijkmnpqstuvxyz",
			Result: false,
		},
		{
			Words:  []string{"apple", "app"},
			Order:  "abcdefghijklmnopqrstuvwxyz",
			Result: false,
		},
		{
			Words:  []string{"kuvp", "q"},
			Order:  "ngxlkthsjuoqcpavbfdermiywz",
			Result: true,
		},
		{
			Words:  []string{"apap", "app"},
			Order:  "abcdefghijklmnopqrstuvwxyz",
			Result: true,
		},
	} {
		assert.Equal(t, data.Result, leetcode.IsAlienSorted(data.Words, data.Order), data)
	}
}
