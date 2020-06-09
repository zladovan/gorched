package gorched

import (
	"fmt"
	"time"

	tl "github.com/JoelOtter/termloop"
)

// Debug provides access to debugging support.
// By default it is inactive.
// You need to activate it by attaching it to game engine object (termloop.Game) with function Attach.
// You can add debug messages to debug with Log / Logf functions.
var Debug = debug{}

// debug holds debug related information
type debug struct {
	active bool
	engine *tl.Game
	logs   []string
}

// Attach activates debug.
// It starts to collect debug info and it starts rendering it's view on the screen.
func (d *debug) Attach(engine *tl.Game) {
	d.engine = engine
	d.active = true
	d.engine.SetDebugOn(true)

	// add debug view to the screen
	d.engine.Screen().AddEntity(newDebugView(d))
}

// Log adds simple message to debug logs
func (d *debug) Log(s string) {
	if !d.active {
		return
	}
	d.logs = append(d.logs, fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05.000"), s))
	if d.engine != nil {
		d.engine.Log(s)
	}
}

// Logf adds message with formating support like int fmt.Sprintf to debug logs
func (d *debug) Logf(s string, items ...interface{}) {
	d.Log(fmt.Sprintf(s, items...))
}

// lastNLogs returns maximum of last n logs, it could be less if there are not more logs than n
func (d *debug) lastNLogs(n int) []string {
	if n <= 0 {
		return []string{}
	}
	logsCount := len(d.logs)
	if logsCount == 0 {
		return []string{}
	}
	startIndex := logsCount - n
	if startIndex < 0 {
		startIndex = 0
	}
	return d.logs[startIndex:]
}

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
			Debug.Logf("Click on cell x=%d y=%d", e.MouseX, e.MouseY)
			if v.showPallette && e.MouseX < 6 && e.MouseY <= 256/6 {
				Debug.Logf("Color int=%d", e.MouseX+e.MouseY*6-1)
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
