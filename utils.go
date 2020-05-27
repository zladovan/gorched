package gorched

import (
	"strings"

	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/draw"
)

// TODO: cleanup, split and move to more meaningfully named files / packages

// NewMessage creates entity from given string with multiple lines.
// Each line will be printed with given fg (foreground) and bg (background) colors.
func NewMessage(msg string, fg tl.Attr, bg tl.Attr) *tl.Entity {
	lines := strings.Split(strings.TrimSpace(msg), "\n")
	h := len(lines)
	w := len([]rune(lines[0]))
	p := draw.BlankPrinter(w, h).WithColors(fg, bg)
	p.WriteLines(0, 0, lines)
	return tl.NewEntityFromCanvas(0, 0, *p.Canvas)
}

// MoveToScreenCenter changes position of given entity e to make it look to be in the center of give screen s.
// Center point of entity will be same as center point of the screen.
func MoveToScreenCenter(e *tl.Entity, s *tl.Screen) {
	w, h := e.Size()
	sw, sh := s.Size()
	e.SetPosition(sw/2-w/2, sh/2-h/2)
}

// DrawPallette draws all available colors as column with 6 colors per line.
// Given shift will change starting color and will result in printing less than 256 colors.
// Debug only !
func DrawPallette(s *tl.Screen, shift int) {
	for c := 0; c < 256-shift; c++ {
		x := c % 6
		y := c / 6
		s.RenderCell(x, y, &tl.Cell{Fg: tl.Attr(c + shift), Ch: '█'})
	}
}
