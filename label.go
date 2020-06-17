package gorched

import (
	"fmt"

	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/gmath"
	"github.com/zladovan/gorched/physics"
)

// TODO: find more descriptive name

// Label represents text entity which is hidden after one second if it's not updated.
type Label struct {
	// it extends text entity
	*tl.Text
	// position of center of this label
	position gmath.Vector2i
	// how many seconds will be label visible
	maxttl float64
	// remaining seconds for label to be visible
	ttl float64
}

// NewLabel creates new label with center in position defined by given x and y coordinates.
// Label has not text yet and it's hidden. You need to call Show().
func NewLabel(x, y int, color tl.Attr) *Label {
	return &Label{
		Text:     tl.NewText(x, y, "", color|tl.AttrBold, tl.ColorDefault),
		position: gmath.Vector2i{X: x, Y: y},
		maxttl:   1,
	}
}

// Show sets some text to the label and show it for one second.
func (l *Label) Show(s string) {
	l.ttl = l.maxttl
	l.Text.SetText(s)
	l.Text.SetPosition(l.position.X-len(s)/2, l.position.Y)
}

// ShowNumber sets some number as text.
// See Show().
func (l *Label) ShowNumber(i int) {
	l.Show(fmt.Sprintf("%d", i))
}

// Draw draws label if it is not out of ttl
func (l *Label) Draw(s *tl.Screen) {
	if l.ttl > 0 {
		l.Text.Draw(s)
		l.ttl -= s.TimeDelta()
	}
}

// FlyingLabel is text entity which will fly up for two seconds and then it removes itself from world.
// After create you need to call one of Show methods to show some text.
type FlyingLabel struct {
	*Label
	body *physics.Body
}

// NewFlyingLabel creates FlyingLabel on given position with given color.
// To set some text use one of the Show methods.
func NewFlyingLabel(position gmath.Vector2i, color tl.Attr) *FlyingLabel {
	l := NewLabel(position.X, position.Y, color)
	l.maxttl = 2
	return &FlyingLabel{
		Label: l,
		body: &physics.Body{
			Position: *position.As2F(),
			Mass:     0.5,
			Velocity: gmath.Vector2f{X: 0, Y: -8},
		},
	}
}

// Draw draws label if it is not out of ttl
func (l *FlyingLabel) Draw(s *tl.Screen) {
	// update label y coordinate based on physical body
	lx, _ := l.Label.Text.Position()
	l.Label.Text.SetPosition(lx, l.body.Position.As2I().Y)

	// draw original label
	l.Label.Draw(s)

	// after ttl remove from level
	if l.ttl <= 0 {
		s.Level().RemoveEntity(l)
	}
}

// Body returns physical body of this label
func (l *FlyingLabel) Body() *physics.Body {
	return l.body
}

// ZIndex return z-index of the label
// It should be higher than z-index of tank and trees but lower z-index of explosion.
func (l *FlyingLabel) ZIndex() int {
	return 2001
}
