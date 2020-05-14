package gorched

// TODO: transform to Vector and VectorF

// Position holds x and y coordinates in console pixels
type Position struct {
	x, y int
}

// Coords return x and y coordinates of this position in console pixels
func (p *Position) Coords() (int, int) {
	return p.x, p.y
}
