package demo

import (
	"os"

	"github.com/nsf/termbox-go"
)

// Wait will wait for given Seconds to be finished.
// It is not blocking all other entities during wait only next commands.
type Wait struct {
	Seconds float64
}

// Eval evaluates command
func (w *Wait) Eval(c *GameContext) bool {
	return c.Memory.AddFloat("t", c.Dt) > w.Seconds
}

// HideMessageBox will hide any messagebox if visible
type HideMessageBox struct{}

// Eval evaluates command
func (t *HideMessageBox) Eval(c *GameContext) bool {
	c.Controls.HideMessageBox()
	return true
}

// SetAngle will change tank's angle to specified Angle.
// Changing is not done immediately but it tries to simulate continuos change like when done via kayboard.
type SetAngle struct {
	Angle int
}

// Eval evaluates command
func (a *SetAngle) Eval(c *GameContext) bool {
	if c.Round.ActiveTank().Angle() == a.Angle {
		return true
	}
	t := c.Memory.AddFloat("t", c.Dt)
	if t > 0.01 {
		if c.Round.ActiveTank().Angle() < a.Angle {
			c.Round.ActiveTank().MoveUp()
		} else {
			c.Round.ActiveTank().MoveDown()
		}
		c.Memory.Clear("t")
	}
	return false
}

// Shoot will make tank to load and shoot with given Power
type Shoot struct {
	Power int
}

// Eval evaluates command
func (s *Shoot) Eval(c *GameContext) bool {
	if c.Round.ActiveTank().IsIdle() {
		c.Controls.Shoot()
	}
	if c.Round.ActiveTank().IsLoading() && c.Round.ActiveTank().Power() >= s.Power {
		c.Controls.Shoot()
		return true
	}
	return false
}

// WaitForFinishTurn will wait until there are no bullets and explisions in the game world.
// Use it after Shoot if you want to invoke next actions for next player.
// It is not blocking all other entities during wait only next commands.
type WaitForFinishTurn struct{}

// Eval evaluates command
func (w *WaitForFinishTurn) Eval(c *GameContext) bool {
	return c.Round.IsTurnFinished()
}

// NextRound will switch game to the next round.
type NextRound struct{}

// Eval evaluates command
func (n *NextRound) Eval(c *GameContext) bool {
	c.Controls.NextRound()
	return true
}

// Exit exits the game
type Exit struct{}

// Eval evaluates command
func (e *Exit) Eval(c *GameContext) bool {
	// TODO: find safer way how stop the game
	termbox.Close()
	os.Exit(0)
	return true
}

// MoveFocus moves focus on ui form to next component
type MoveFocus struct{}

// Eval evaluates command
func (m *MoveFocus) Eval(c *GameContext) bool {
	c.Controls.MoveFocus()
	return true
}

// PressButton sends action event to currently focused ui component
type PressButton struct{}

// Eval evaluates command
func (p *PressButton) Eval(c *GameContext) bool {
	c.Controls.PressButton()
	return true
}
