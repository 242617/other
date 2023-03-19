package leetcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/242617/other/tasks/leetcode"
)

func TestConvertTemperature(t *testing.T) {
	for _, data := range []struct {
		Celsius float64
		Output  []float64
	}{
		{
			Celsius: 36.50,
			Output:  []float64{309.65000, 97.70000},
		},
		{
			Celsius: 122.11,
			Output:  []float64{395.26000, 251.79800},
		},
	} {
		assert.Equal(t, data.Output, leetcode.ConvertTemperature(data.Celsius), data)
	}
}
