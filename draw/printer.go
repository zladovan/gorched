package draw

import (
	tl "github.com/JoelOtter/termloop"
)

// Printer allows to write some characters to termloop.Canvas
type Printer struct {
	// canvas is target where all writings will be done
	Canvas *tl.Canvas
	// fg holds foreground color which will be used for all writings
	Fg tl.Attr
	// bg holds background color which will be used for all writings
	Bg tl.Attr
}

// BlankPrinter creates new Printer for newly created termloop.Canvas with given width and height.
// This method can be used for creating new termloop.Canvas if you want directly start to print something on it.
func BlankPrinter(width, height int) *Printer {
	canvas := tl.NewCanvas(width, height)
	return &Printer{Canvas: &canvas}
}

// Write given string on position defined by x and y coordinates.
// If string would end out of canvas it will be stripped.
func (p *Printer) Write(x, y int, s string) {
	for i, c := range []rune(s) {
		p.WritePoint(x+i, y, c)
	}
}

// WritePoint prints rune c on canvas on position defined by x and y coordinates
func (p *Printer) WritePoint(x, y int, c rune) {
	p.drawCell(x, y, tl.Cell{
		Fg: p.Fg,
		Bg: p.Bg,
		Ch: c,
	})
}

// drawCell will set cell c on position given by x and y coordinates
// if given coordinates are out of range it will silently ignore it
// if background color of given cell is 0 it is considered as transaprent
// if foreground color of given cell is 0 foreground of underlying cell will be used
func (p *Printer) drawCell(x, y int, c tl.Cell) {
	// check if it's possible to draw on given coordinates
	if x < 0 || x > p.MaxX() {
		return
	}
	if y < 0 || y > p.MaxY() {
		return
	}

	// rune is always changed
	(*p.Canvas)[x][y].Ch = c.Ch

	// keep original bg if cell to be drawn has zero background
	if c.Bg != 0 {
		(*p.Canvas)[x][y].Bg = c.Bg
	}

	// keep original fg if cell to be drawn has zero foreground
	if c.Fg != 0 {
		(*p.Canvas)[x][y].Fg = c.Fg
	}
}

// WriteCenterX writes given string on line y with horizontal alignment.
func (p *Printer) WriteCenterX(y int, s string) {
	p.Write(p.CenterX()-len([]rune(s))/2, y, s)
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

// WriteHorizontal prints given text s horizontally.
// It will start on position given by x and y.
// Depending on up flag it will increase (up=false) or decrease (up=true) y position for each next rune in given text s.
func (p *Printer) WriteHorizontal(x, y int, s string, up bool) {
	dir := 1
	if up {
		dir = -1
	}
	for i, c := range []rune(s) {
		p.WritePoint(x, y+i*dir, c)
	}
}

// WriteHorizontalUp prints given text s horizontally.
// It will start on position given by x and y and it will decreasing y for each next rune in given text s.
func (p *Printer) WriteHorizontalUp(x, y int, s string) {
	p.WriteHorizontal(x, y, s, true)
}

// WriteHorizontalDown prints given text s horizontally.
// It will start on position given by x and y and it will increasing y for each next rune in given text s.
func (p *Printer) WriteHorizontalDown(x, y int, s string) {
	p.WriteHorizontal(x, y, s, false)
}

// Draw draws another canvas on given position x and y
func (p *Printer) Draw(x, y int, canvas *tl.Canvas) {
	for i := 0; i < len(*canvas); i++ {
		for j := 0; j < len((*canvas)[0]); j++ {
			p.drawCell(x+i, y+j, (*canvas)[i][j])
		}
	}
}

// Fill fills entire canvas with given rune c
func (p *Printer) Fill(c rune) {
	for x := 0; x < p.Width(); x++ {
		for y := 0; y < p.Height(); y++ {
			p.WritePoint(x, y, c)
		}
	}
}

// Width returns width of target canvas
func (p *Printer) Width() int {
	return len(*p.Canvas)
}

// Height returns height of target canvas
func (p *Printer) Height() int {
	return len((*p.Canvas)[0])
}

// CenterX returns x coordinate of the center point on target canvas
func (p *Printer) CenterX() int {
	return p.Width() / 2
}

// CenterY returns y coordinate of the center point on target canvas
func (p *Printer) CenterY() int {
	return p.Height() / 2
}

// MaxX returns highest x coordinate on target canvas
func (p *Printer) MaxX() int {
	return p.Width() - 1
}

// MaxY returns highest y coordinate on target canvas
func (p *Printer) MaxY() int {
	return p.Height() - 1
}

// WithFg creates new Printer for the same target canvas but with changed foreground color
func (p *Printer) WithFg(fg tl.Attr) *Printer {
	return p.WithColors(fg, p.Bg)
}

// WithBg creates new Printer for the same target canvas but with changed background color
func (p *Printer) WithBg(bg tl.Attr) *Printer {
	return p.WithColors(p.Fg, bg)
}

// WithDefaultBg cretes new Printer for the same target canvas but with default background color
func (p *Printer) WithDefaultBg() *Printer {
	return p.WithBg(tl.ColorDefault)
}

// WithColors cretes new Printer for the same target canvas but with changed foreground and background colors
func (p *Printer) WithColors(fg tl.Attr, bg tl.Attr) *Printer {
	return &Printer{
		Canvas: p.Canvas,
		Fg:     fg,
		Bg:     bg,
	}
}
