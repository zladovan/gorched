package draw

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/gmath"
)

// RadialGradient can be used for for resolving colors for gradient effect.
// You need to specify two vectors A and B where A is the middle of the circle and B is point on circle line.
// Then you need to set a Color for vector A which will be modified along path to vector B.
// Distance between A and B will be split to StepCount equal slices. Color in each next slice is calculated as previous slice color + Step.
type RadialGradient struct {
	A, B            gmath.Vector2i
	ColorA          tl.Attr
	Step, StepCount int
}

// Color calculates color for given vector.
// See ShadesGradient for details.
func (g *RadialGradient) Color(c gmath.Vector2i) tl.Attr {
	steps := float64(g.StepCount)
	maxDistance := g.A.DistanceF(&g.B)
	distance := g.A.DistanceF(&c)

	// we want to scale steps only if there is not enough steps for maximum distance
	if maxDistance > steps {
		distance = distance / maxDistance * steps
	}
	// ensure we will not go over maximum step
	if distance > steps-1 {
		distance = steps - 1
	}

	return tl.Attr(int(g.ColorA) + g.Step*int(distance))
}