package ui

import (
	"strings"

	"github.com/zladovan/gorched/draw"
	"github.com/zladovan/gorched/gmath"
)

// Text is simple component showing some text.
// It supports multiple lines.
// It's size will be derived from length of the longest line and the lines count.
// By default it will use colors defined by the container where it will be added.
type Text struct {
	*BaseComponent
	// lines contains text as array of lines
	lines []string
	// Colors defines foreground and background color for the text
	Colors Colors
}

// NewTextFromLines creates new Text component from given array of lines
func NewTextFromLines(lines []string) *Text {
	return &Text{BaseComponent: &BaseComponent{}, lines: lines}
}

// NewText creates new Text component from given text.
// Text can contain multiple lines.
func NewText(text string) *Text {
	return NewTextFromLines(strings.Split(text, "\n"))
}

// Dimensions returns length of longest line as X and lines count as Y
func (t *Text) Dimensions() gmath.Vector2i {
	w := 0
	for _, line := range t.lines {
		l := len([]rune(line))
		if l > w {
			w = l
		}
	}
	return gmath.Vector2i{X: w, Y: len(t.lines)}
}

// Refresh redraws this Text component to it's canvas
func (t *Text) Refresh() {
	d := t.Dimensions()
	p := draw.BlankPrinter(d.X, d.Y)
	p.Fg = t.Colors.Fg
	p.Bg = t.Colors.Bg
	p.Fill(' ')
	p.WriteLines(0, 0, t.lines)
	t.SetCanvas(p.Canvas)
}
