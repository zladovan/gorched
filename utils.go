package gorched

import (
	"strings"

	tl "github.com/JoelOtter/termloop"
)

// TODO: cleanup, split and move to more meaningfully named files / packages

// Printer allows to write some characters to termloop.Canvas
type Printer struct {
	// canvas is target where all writing will be done
	canvas *tl.Canvas
	//
	fg tl.Attr
	bg tl.Attr
}

// Write given string on position defined by x and y coordinates.
// If string would end out of canvas it will be stripped.
func (p *Printer) Write(x, y int, s string) {
	i := 0
	for _, c := range s {
		if x+i >= len((*p.canvas)) {
			break
		}
		(*p.canvas)[x+i][y].Fg = p.fg
		(*p.canvas)[x+i][y].Bg = p.bg
		(*p.canvas)[x+i][y].Ch = c
		i++
	}
}

// WriteCenterX writes given string on line y with horizontal alignment.
func (p *Printer) WriteCenterX(y int, s string) {
	p.Write(len(*p.canvas)/2-len(s)/2, y, s)
}

// WriteLines writes all given lines.
// Each line will start on x coordinate.
// First line will be on row y.
// Each next line will be on separated row under y.
func (p *Printer) WriteLines(x, y int, lines []string) {
	for i, l := range lines {
		p.Write(x, y+i, l)
	}
}

// NewMessage creates entity from given string with multiple lines.
// Each line will be printed with given fg (foreground) and bg (background) colors.
func NewMessage(msg string, fg tl.Attr, bg tl.Attr) *tl.Entity {
	lines := strings.Split(strings.TrimSpace(msg), "\n")
	h := len(lines)
	w := len([]rune(lines[0]))
	canvas := tl.NewCanvas(w, h)
	p := &Printer{canvas: &canvas, fg: fg, bg: bg}
	p.WriteLines(0, 0, lines)
	return tl.NewEntityFromCanvas(0, 0, canvas)
}

// MoveToScreenCenter changes position of given entity e to make it look to be in the center of give screen s.
// Center point of entity will be same as center point of the screen.
func MoveToScreenCenter(e *tl.Entity, s *tl.Screen) {
	w, h := e.Size()
	sw, sh := s.Size()
	e.SetPosition(sw/2-w/2, sh/2-h/2)
}

// DrawPallete draws all available colors as column with 6 colors per line.
// Given shift will change starting color and will result in printing less than 256 colors.
// Debug only !
func DrawPallette(s *tl.Screen, shift int) {
	for c := 0; c < 256-shift; c++ {
		x := c % 6
		y := c / 6
		s.RenderCell(x, y, &tl.Cell{Fg: tl.Attr(c + shift), Ch: '█'})
	}
}
