package gmath

import "math"

// Min returns smaller one from given values x and y.
// It's integer version of math.Min.
func Min(x, y int) int {
	return int(math.Min(float64(x), float64(y)))
}

// Max returns bigger one from given values x and y.
// It's integer version of math.Max
func Max(x, y int) int {
	return int(math.Max(float64(x), float64(y)))
}

// Clamp returns given value if it's in range given by min and max.
// If value < min then min will be returned.
// If value > max then max will be returned.
func Clamp(min, max, value float64) float64 {
	return math.Min(math.Max(min, value), max)
}
