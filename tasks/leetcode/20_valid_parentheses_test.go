package leetcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/242617/other/tasks/leetcode"
)

func TestIsValid(t *testing.T) {
	for _, data := range []struct {
		String string
		Result bool
	}{
		{
			String: "()",
			Result: true,
		},
		{
			String: "([])",
			Result: true,
		},
		{
			String: "(]",
			Result: false,
		},
		{
			String: "((",
			Result: false,
		},
		{
			String: "(){}}{",
			Result: false,
		},
	} {
		assert.Equal(t, data.Result, leetcode.IsValid(data.String), data)
	}
}
