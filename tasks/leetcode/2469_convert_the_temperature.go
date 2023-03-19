package leetcode

import "math"

// Kelvin = Celsius + 273.15
// Fahrenheit = Celsius * 1.80 + 32.00
func ConvertTemperature(celsius float64) []float64 {
	return []float64{
		celsius + 273.15,
		math.Round((celsius*1.80+32.00)*10000) / 10000,
	}
}
