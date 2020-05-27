package gmath

import "math"

// Vector2i represents 2-dimensional vector with integer components
type Vector2i struct {
	X, Y int
}

// Distance returns distance of this vector v to given vector u rounded to nearest integer (half away from zero)
func (v *Vector2i) Distance(u *Vector2i) int {
	return int(math.Round(v.DistanceF(u)))
}

// DistanceF returns distance of this vector v to given vector u
func (v *Vector2i) DistanceF(u *Vector2i) float64 {
	return math.Sqrt(math.Pow(float64(v.X-u.X), 2) + math.Pow(float64(v.Y-u.Y), 2))
}
