package debug

import (
	"fmt"
	"time"

	tl "github.com/JoelOtter/termloop"
)

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
