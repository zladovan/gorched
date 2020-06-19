package gorched

import (
	"fmt"

	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/gmath"
	"github.com/zladovan/gorched/physics"
)

// Label is text entity with one row of text.
type Label struct {
	// position is reference position of this label
	position gmath.Vector2i
	// format defines text formatting
	format Formatting
	// entity is used to draw text to the screen
	entity *tl.Text
	// text is text which will be drawn
	text string
}

// Alignment defines types how the text of the label will be shifted from Label's position
type Alignment = uint8

const (
	// Center alignment will draw test with Label position as the center
	Center Alignment = iota
	// Right alignment will draw text starting from Label position
	Right
	// Left alignment will draw text ending in Label position
	Left
)

// Formatting defines text formatting options for Label
type Formatting struct {
	// Color is foreground color
	Color tl.Attr
	// Background is background color
	Background tl.Attr
	// Align is text alignment
	Align Alignment
}

// NewLabel creates new label.
// Given position is reference position and final entity could be moved from this position according used fromat alignment.
func NewLabel(position gmath.Vector2i, text string, format Formatting) *Label {
	l := &Label{
		position: position,
		format:   format,
		text:     text,
	}
	l.refresh()
	return l
}

// Draw will draw this label to the screen
func (l *Label) Draw(s *tl.Screen) {
	l.entity.Draw(s)
}

// Tick does nothing now
func (l *Label) Tick(e tl.Event) {}

// SetText will change text of this label
func (l *Label) SetText(text string) {
	l.text = text
	l.refresh()
}

// SetPosition will change position of this label.
// Given position is reference position and final entity could be moved from this position according used fromat alignment.
func (l *Label) SetPosition(p gmath.Vector2i) {
	l.position = p
	l.refresh()
}

// Position returns reference position of this label
func (l *Label) Position() gmath.Vector2i {
	return l.position
}

// refresh recreate wrapped termloop.Text entity
func (l *Label) refresh() {
	len := len([]rune(l.text))
	x := l.position.X
	y := l.position.Y
	switch l.format.Align {
	case Left:
		x -= len
	case Center:
		x -= len / 2
	}
	l.entity = tl.NewText(x, y, l.text, l.format.Color|tl.AttrBold, l.format.Background)
}

// ZIndex return z-index of the label
// It should be higher than z-index of tank and trees but lower z-index of explosion.
func (l *Label) ZIndex() int {
	return 2001
}

// TempLabel is Label which is hidden after TTL seconds if it's not updated with one of Show methods.
// If you want to make it visible right after the creation set RemainingTTL to non zero value.
// Otherwise it will be shown after first call of one of Show methods.
type TempLabel struct {
	// it extends from Label
	*Label
	// TTL is how many seconds will be label visible when shown
	TTL float64
	// RemainingTTL is how many seconds remains to be hidden
	RemainingTTL float64
	// Remove if is true label will be removed from world after TTL seconds
	Remove bool
}

// Show makes label again visible for TTL seconds.
func (l *TempLabel) Show() {
	l.RemainingTTL = l.TTL
}

// ShowText sets some text to the label and show it for TTL seconds.
func (l *TempLabel) ShowText(s string) {
	l.SetText(s)
	l.Show()
}

// ShowNumber sets some number as to the label and show it for TTL seconds.
// See ShowText().
func (l *TempLabel) ShowNumber(i int) {
	l.ShowText(fmt.Sprintf("%d", i))
}

// Draw draws label if it is not out of ttl
func (l *TempLabel) Draw(s *tl.Screen) {
	if l.IsVisible() {
		l.Label.Draw(s)
		l.RemainingTTL -= s.TimeDelta()
	} else if l.Remove {
		s.Level().RemoveEntity(l)
	}
}

// IsVisible returns true if label is not yet ouf of time to be drawn
func (l *TempLabel) IsVisible() bool {
	return l.RemainingTTL > 0
}

// FlyingLabel is text entity which will fly up for two seconds and then it removes itself from world.
type FlyingLabel struct {
	*TempLabel
	body *physics.Body
}

// NewFlyingLabel creates FlyingLabel on given position with given text and fromatting.
// To set some text use one of the Show methods.
func NewFlyingLabel(position gmath.Vector2i, text string, format Formatting) *FlyingLabel {
	return &FlyingLabel{
		TempLabel: &TempLabel{
			Label:        NewLabel(position, text, format),
			TTL:          2,
			RemainingTTL: 2,
			Remove:       true,
		},
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
	l.TempLabel.SetPosition(*l.body.Position.As2I())

	// draw original label
	l.TempLabel.Draw(s)
}

// Body returns physical body of this label
func (l *FlyingLabel) Body() *physics.Body {
	return l.body
}

// ZIndex return z-index of the flying label
// It should be higher than z-index of standard label.
func (l *FlyingLabel) ZIndex() int {
	return l.TempLabel.ZIndex() + 1
}
