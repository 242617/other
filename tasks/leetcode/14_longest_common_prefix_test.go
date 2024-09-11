package leetcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/242617/other/tasks/leetcode"
)

func TestLongestCommonPrefix(t *testing.T) {
	for _, data := range []struct {
		Strings []string
		Result  string
	}{
		{
			Strings: []string{"aa", "aab", "ccc", "cccd", "eeee"},
			Result:  "",
		},
		{
			Strings: []string{"flower", "flow", "flight"},
			Result:  "fl",
		},
		{
			Strings: []string{"dog", "racecar", "car"},
			Result:  "",
		},
	} {
		assert.Equal(t, data.Result, leetcode.LongestCommonPrefix(data.Strings), data)
	}
}
