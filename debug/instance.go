// Package debug provides access to debugging support.
// By default it is inactive.
// You need to activate it by attaching it to game engine object (termloop.Game) with function debug.Attach.
// After attaching you can add debug messages to debug with debug.Log / debug.Logf functions.
package debug

import tl "github.com/JoelOtter/termloop"

// debug has only single instace
var instance = debug{}

// Attach activates debug.
// It starts to collect debug info and it starts rendering it's view on the screen.
func Attach(engine *tl.Game) {
	instance.Attach(engine)
}

// Log adds simple message to debug logs
func Log(s string) {
	instance.Log(s)
}

// Logf adds message with formating support like int fmt.Sprintf to debug logs
func Logf(s string, items ...interface{}) {
	instance.Logf(s, items...)
}
