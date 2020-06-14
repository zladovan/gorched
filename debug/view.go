package debug

import (
	tl "github.com/JoelOtter/termloop"
)

// debugView is entity for showing all debug related info
type debugView struct {
	d            *debug
	fpsCaption   *tl.Text
	fpsText      *tl.FpsText
	hidden       bool
	showPallette bool
}

// newDebugView creates new entity for showing all debug related info
func newDebugView(d *debug) *debugView {
	return &debugView{
		d:          d,
		fpsCaption: tl.NewText(0, 0, "Fps: ", tl.ColorBlack, tl.ColorDefault),
		fpsText:    tl.NewFpsText(5, 0, tl.ColorBlack, tl.ColorDefault, 1),
	}
}

// Draw is drawing logs to the screen
func (v *debugView) Draw(s *tl.Screen) {
	if v.hidden {
		return
	}

	// starting x position of debug text lines
	x := 0
	if v.showPallette {
		DrawPallette(s, -1)
		x += 6
	}

	// show line with fps
	v.fpsCaption.SetPosition(x, 0)
	v.fpsCaption.Draw(s)
	v.fpsText.SetPosition(x+5, 0)
	v.fpsText.Draw(s)

	// show last 10 logs with new log always at the bottom
	for i, l := range v.d.lastNLogs(10) {
		tl.NewText(x, i+1, l, tl.ColorBlack, tl.ColorDefault).Draw(s)
	}
}

// Tick updates debug view on every tick
func (v *debugView) Tick(e tl.Event) {
	// toggle debug view
	switch e.Key {
	case tl.KeyCtrlD:
		v.hidden = !v.hidden
	}

	// toggle pallette
	switch e.Ch {
	case 'p':
		v.showPallette = !v.showPallette
	}

	// log mouse click position
	switch e.Type {
	case tl.EventMouse:
		if e.Key == tl.MouseLeft {
			Logf("Click on cell x=%d y=%d", e.MouseX, e.MouseY)
			if v.showPallette && e.MouseX < 6 && e.MouseY <= 256/6 {
				Logf("Color int=%d", e.MouseX+e.MouseY*6-1)
			}
		}
	}
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
