package gmath

import "math"

// Vector2i represents 2-dimensional vector with integer components
type Vector2i struct {
	X, Y int
}

// Distance returns distance of this vector v to given vector as integer
func (v *Vector2i) Distance(u *Vector2i) int {
	return int(math.Round(v.DistanceF(u)))
}

// DistanceF returns distance of this vector v to given vector u
func (v *Vector2i) DistanceF(u *Vector2i) float64 {
	return v.As2F().Distance(u.As2F())
}

// Translate returns new vector which is moved by given x and y from this vector.
func (v *Vector2i) Translate(x, y int) *Vector2i {
	return &Vector2i{X: v.X + x, Y: v.Y + y}
}

// RotateAround returns new vector created by rotating this vector around given vector u by angle in radians given by rad.
func (v *Vector2i) RotateAround(u *Vector2i, rad float64) *Vector2i {
	return v.As2F().RotateAround(u.As2F(), rad).As2I()
}

// As2F Converts this vector to Vector2f with float64 coordinates
func (v *Vector2i) As2F() *Vector2f {
	return &Vector2f{X: float64(v.X), Y: float64(v.Y)}
}

// Vector2f represents 2-dimensional vector with float64 components
type Vector2f struct {
	X, Y float64
}

// RotateAround returns new vector created by rotating this vector around given vector u by angle in radians given by rad.
func (v *Vector2f) RotateAround(u *Vector2f, rad float64) *Vector2f {

	ms := math.Sin(rad)
	mc := math.Cos(rad)

	// translate to origin
	px := v.X - u.X
	py := v.Y - u.Y

	// rotate
	x := px*mc - py*ms
	y := px*ms + py*mc

	// translate back from origin
	x += u.X
	y += u.Y

	return &Vector2f{X: x, Y: y}
}

// Distance returns distance of this vector v to given vector u
func (v *Vector2f) Distance(u *Vector2f) float64 {
	return math.Sqrt(math.Pow(v.X-u.X, 2) + math.Pow(v.Y-u.Y, 2))
}

// As2I Converts this vector to Vector2i with int coordinates
func (v *Vector2f) As2I() *Vector2i {
	return &Vector2i{X: int(v.X), Y: int(v.Y)}
}

// TODO: add more operations
