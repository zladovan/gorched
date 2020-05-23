package gorched

import (
	"fmt"

	tl "github.com/JoelOtter/termloop"
)

// TODO: find more descriptive name

// Label represents text entity which is hidden after one second if it's not updated.
type Label struct {
	// it extends text entity
	*tl.Text
	// position of center of this label
	position Position
	// how many seconds will be label visible
	ttl float64
}

// NewLabel creates new label with center in position defined by given x and y coordinates.
// Label has not text yet and it's hidden. You need to call Show().
func NewLabel(x, y int, color tl.Attr) *Label {
	return &Label{
		Text:     tl.NewText(x, y, "", color|tl.AttrBold, tl.ColorDefault),
		position: Position{x, y},
	}
}

// Show sets some text to the label and show it for one second.
func (l *Label) Show(s string) {
	l.ttl = 1
	l.Text.SetText(s)
	l.Text.SetPosition(l.position.x-len(s)/2, l.position.y)
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
