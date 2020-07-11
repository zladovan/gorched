package entities

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/draw"
	"github.com/zladovan/gorched/gmath"
	"github.com/zladovan/gorched/physics"
)

// Tomb is entity representing tomb stone shown on position where tank was killed
type Tomb struct {
	body   *physics.Body
	canvas *tl.Canvas
}

// NewTomb creates new Tomb entity on given position with given color.
// It should have same color as tank instead of which it was added to world.
func NewTomb(position gmath.Vector2i, color tl.Attr) *Tomb {
	return &Tomb{
		body: &physics.Body{
			Position: *position.As2F(),
			Mass:     3,
		},
		canvas: createTombCanvas(color),
	}
}

// createTombCanvas creates canvas with tomb sprite
func createTombCanvas(color tl.Attr) *tl.Canvas {
	p := draw.BlankPrinter(3, 2).WithFg(color)
	p.WriteLines(0, 0, []string{
		"▄█▄",
		" █",
	})
	return p.Canvas
}

// Draw draws tomb stone
func (t *Tomb) Draw(s *tl.Screen) {
	offsetx := -1
	offsety := -2
	for i := 0; i < len(*t.canvas); i++ {
		for j := 0; j < len((*t.canvas)[0]); j++ {
			s.RenderCell(int(t.body.Position.X)+i+offsetx, int(t.body.Position.Y)+j+offsety, &(*t.canvas)[i][j])
		}
	}
}

// Tick does nothing now
func (t *Tomb) Tick(e tl.Event) {}

// ZIndex return z-index of the tomb.
// It should be lower than z-index of tank.
func (t *Tomb) ZIndex() int {
	return 1999
}

// Body returns physical body of the tomb used for falling simulation
func (t *Tomb) Body() *physics.Body {
	return t.body
}

// BottomLine returns line x coordinates for collision with the ground when falling
func (t *Tomb) BottomLine() (int, int) {
	return 0, 0
}
