package terrain

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/gmath"
	"github.com/zladovan/gorched/physics"
)

// Column is entity representing by 1 console pixel wide rectangle with height for it's corresponding part in terrain.
type Column struct {
	*tl.Entity
	body       *physics.Body
	terrain    *Terrain
	canvas     *tl.Canvas
	bodyLocker *physics.TimeLocker
}

// NewColumn creates new Column entity for given terrain.
// Position defined by x and y should be from terrain line.
func NewColumn(terrain *Terrain, x, y int, canvas *tl.Canvas) *Column {
	body := &physics.Body{
		Position: gmath.Vector2f{X: float64(x), Y: float64(y + len((*canvas)[0]))},
		Mass:     5,
		Locked:   true,
	}
	return &Column{
		Entity:     tl.NewEntityFromCanvas(x, y, *canvas),
		body:       body,
		terrain:    terrain,
		canvas:     canvas,
		bodyLocker: &physics.TimeLocker{BodyToRelock: body, RemainingSeconds: 0.5},
	}
}

// Draw draws this column and process possible cuttings
func (t *Column) Draw(s *tl.Screen) {
	// update body locker
	t.bodyLocker.Update(s.TimeDelta())

	// update position of entity based on body position if not locked
	if !t.body.Locked {
		t.Entity.SetPosition(int(t.body.Position.X), int(t.body.Position.Y)-len((*t.canvas)[0]))
	}

	// draw entity
	t.Entity.Draw(s)
}

// Position returns top-left position of collider
func (t *Column) Position() (int, int) {
	return t.Entity.Position()
}

// Size returns size of collider
func (t *Column) Size() (int, int) {
	return t.Entity.Size()
}

// Body returns physical body for falling processing
func (t *Column) Body() *physics.Body {
	return t.body
}

// BottomLine returns single point at column's body position for collision with the ground when falling
func (t *Column) BottomLine() (int, int) {
	return 0, 0
}

// MakeHole will create hole in terrain which this column is part of with center at cx and cy coordinates with given radius r.
func (t *Column) MakeHole(cx, cy, r int) {
	t.terrain.MakeHole(cx, cy, r)
}
