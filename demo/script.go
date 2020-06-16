package demo

import (
	"github.com/zladovan/gorched"
)

// Script is just array of commands.
type Script []Command

// Command represents some activity which can be done during game play.
// Each command can be executed multiple times.
// How many times will be command executed depends on its return value.
// If command returns true it is considered as finished. Otherwise it will be called again.
type Command interface {
	Eval(c *GameContext) bool
}

// GameContext holds access to the game controls and state.
type GameContext struct {
	// Controls refers to the game controls which can be used to do some action
	Controls *gorched.Controls
	// Round refers to current game round which can be used to getting the state information
	Round *gorched.Round
	// Dt is delta time between two executions (game screen redraws)
	Dt float64
	// Memory is local memory which "lives" during execution of single command
	Memory *Memory
}

// Memory is simple memory which can be used by commands to store some values under string identifiers
type Memory struct {
	content map[string]interface{}
}

// NewMemory creates new memory with no values
func NewMemory() *Memory {
	return &Memory{content: make(map[string]interface{})}
}

// Read reads value stored on given id from memory
func (m *Memory) Read(id string) interface{} {
	return m.content[id]
}

// Write stores value to memory on given id
func (m *Memory) Write(id string, v interface{}) {
	m.content[id] = v
}

// ReadFloat reads float64 value stored on given id from memory.
// If there is no float64 value stored on id 0 is returned.
func (m *Memory) ReadFloat(id string) float64 {
	if f, ok := m.content[id].(float64); ok {
		return f
	}
	return 0
}

// WriteFloat stores float64 value to memory on given id.
// If there was already some value stored on id it will be overwritten.
func (m *Memory) WriteFloat(id string, f float64) {
	m.content[id] = f
}

// UpdateFloat reads value from memory stored on given id, applies function f to this value and writes result back on id.
// Result of function apply is returned too.
func (m *Memory) UpdateFloat(id string, f func(f float64) float64) float64 {
	r := f(m.ReadFloat(id))
	m.WriteFloat(id, r)
	return r
}

// AddFloat adds given add float to the float value stored in memory on given id.
func (m *Memory) AddFloat(id string, add float64) float64 {
	return m.UpdateFloat(id, func(f float64) float64 {
		return f + add
	})
}

// Clear removes value stored on given id.
func (m *Memory) Clear(id string) {
	m.content[id] = nil
}

// ClearAll removes all values stored in memory
func (m *Memory) ClearAll() {
	m.content = make(map[string]interface{})
}
