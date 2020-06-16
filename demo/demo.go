// Package demo contains support for playing scripted demos.
//
// There is an entity Demo which can be added to the screen and it will play demo Script.
// Demo Script consist of ordered commands. Each command can perform some in game action.
//
// There is support for parsing demo script from string representation.
// Script string should contain one command per line.
// It can contain empty lines and commented lines starting with `#`.
// Alternatively there can be multiple commands per line separated by `;` (semicolon).
//
// Supported commands:
//
//   - `wait s` - will wait for `s` seconds
//   - `hideMessageBox` - hides any visible message box
//	 - `setAngle a` - continuosly change cannon angle to `a` for active player
//   - `shoot p` - load power to `p` and shoot with active player
//   - `waitForFinishTurn` - will wait for all explosions and bullets are gone and turn is shifted to next player
//	 - `nextRound` - switches game to the next round
//   - `exit` - exits game
package demo

import (
	"os"

	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched"
)

// Demo is entity which when added to the game will play given Script.
// Use NewDemo to create new demo with some script.
// Use LoadFromFile to create new demo with script loaded from file.
type Demo struct {
	script              Script
	currentCommandIndex int
	memory              *Memory
}

// NewDemo creates new demo with given script.
func NewDemo(script Script) *Demo {
	return &Demo{script: script, memory: NewMemory()}
}

// LoadFromFile loads demo from file with demo script commands
func LoadFromFile(path string) (*Demo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	script, err := ParseScript(file)
	if err != nil {
		return nil, err
	}
	return NewDemo(script), nil
}

// Draw invokes all demo script commands in order
func (d *Demo) Draw(s *tl.Screen) {
	if d.currentCommandIndex >= len(d.script) {
		return
	}
	command := d.script[d.currentCommandIndex]
	ctx := &GameContext{
		Controls: findControls(s),
		Round:    findRound(s),
		Dt:       s.TimeDelta(),
		Memory:   d.memory,
	}
	if command.Eval(ctx) {
		d.currentCommandIndex++
		d.memory.ClearAll()
	}
}

// Tick does nothing now
func (d *Demo) Tick(e tl.Event) {}

// Restart will make demo run all script commands from start again
func (d *Demo) Restart() {
	d.currentCommandIndex = 0
}

// findControls will look for gorched.Controls entity in current screen
func findControls(s *tl.Screen) *gorched.Controls {
	for _, e := range s.Entities {
		if c, ok := e.(*gorched.Controls); ok {
			return c
		}
	}
	panic("Controls not found")
}

// findControls will look for gorched.Round entity in current screen
func findRound(s *tl.Screen) *gorched.Round {
	for _, e := range s.Entities {
		if c, ok := e.(*gorched.Round); ok {
			return c
		}
	}
	panic("Round not found")
}
