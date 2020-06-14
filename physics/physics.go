package physics

import "github.com/zladovan/gorched/gmath"

// Physics represents very simplistic physical model.
// In this model there are some bodies which can move and fall.
// Optionally body can land on the ground.
// There are no other collisions resolved here instead of landing on the ground.
// If you want to apply physics to your object implement HasBody interface.
// If you want to let your object land on the ground implement Lander interface too.
type Physics struct {
	// Gravity holds gravitational acceleration
	Gravity float64
	// Ground resolves nearest ground y coordinate for position given by x and y
	Ground func(x, y int) int
}

// Body is the main object in physical model.
type Body struct {
	// Position holds position of body
	Position gmath.Vector2f
	// Velocity holds velocity of body
	Velocity gmath.Vector2f
	// Mass holds mass of body
	Mass float64
	// Locked if is true no physics is applied to this body
	Locked bool
}

// HasBody describes some object which can provide physical Body.
type HasBody interface {
	Body() *Body
}

// Lander describes some object which is able to land on the ground.
// It should return x coordinates (start, end) relative to body position of the bottom line used for the collision with ground.
type Lander interface {
	BottomLine() (int, int)
}

// Apply will update body according to this physical model.
// Velocity will be updated by gravitational accelleration.
// Position will be updated by velocity.
// Position could be trimmed to the ground positions if given object e implements Lander interface too.
func (p *Physics) Apply(e HasBody, dt float64) {
	body := e.Body()
	y := int(body.Position.Y)

	// nothing to update if body is locked
	if body.Locked {
		return
	}

	// update velocity by gravity
	body.Velocity.Y += p.Gravity * dt * body.Mass

	// update position by velocity
	body.Position.X += body.Velocity.X * dt
	body.Position.Y += body.Velocity.Y * dt

	// process landing on the ground if possible
	if faller, ok := e.(Lander); ok {
		bx1, bx2 := faller.BottomLine()
		minx := int(body.Position.X) + bx1
		maxx := int(body.Position.X) + bx2

		// check collision with bottom line and the ground
		for i := minx; i <= maxx; i++ {
			groundy := p.Ground(i, y)
			if int(body.Position.Y) > groundy {
				body.Position.Y = float64(groundy)
				body.Velocity.Y = 0
				break
			}
		}
	}
}
