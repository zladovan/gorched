package gorched

import (
	tl "github.com/JoelOtter/termloop"
)

// Controls holds data and logic for controlling game world.
type Controls struct {
	// reference to game
	game *Game
	// all tanks in game
	tanks []*Tank
	// index of tank which is active and controlled
	activeTankIndex int
}

// Tick handles all key events
func (c *Controls) Tick(e tl.Event) {
	// TODO: show some message box after resize about restart round is needed to be applied
	// on resize update game options to be applied on round restart or on next round
	if e.Type == tl.EventResize {
		w, h := c.game.engine.Screen().Size()
		c.game.options.Width = w
		c.game.options.Height = h
	}

	// when message box is show it's just possible to hide it
	if c.game.Hud().IsMessageBoxShown() {
		if e.Type == tl.EventKey {
			c.game.Hud().HideMessageBox()
			if c.NumberOfTanksAlive() <= 1 {
				c.game.NextRound()
			}
		}
		return
	}

	// otherwise handle in-game controls
	switch e.Key {
	case tl.KeyArrowLeft:
		c.ActiveTank().MoveUp()
	case tl.KeyArrowRight:
		c.ActiveTank().MoveDown()
	case tl.KeySpace:
		c.ActiveTank().Shoot(func() {
			if c.NumberOfTanksAlive() <= 1 {
				c.game.Hud().ShowScore()
			} else {
				c.ActivateNextTank()
			}
		})
	case tl.KeyCtrlR:
		c.game.RestartRound()
	case tl.KeyCtrlN:
		c.game.NextRound()
	}
	switch e.Ch {
	case 'h':
		c.game.Hud().ShowInfo()
	case 's':
		c.game.Hud().ShowScore()
	}
	// for the browser mode we cannot use ctrl+n and ctr+r as we would leave the window
	if c.game.options.BrowserMode {
		switch e.Ch {
		case 'n':
			c.game.NextRound()
		case 'r':
			c.game.RestartRound()
		}
	}
}

// ActivateNextTank moves turn to nearest tank which is alive.
func (c *Controls) ActivateNextTank() {
	c.activeTankIndex = (c.activeTankIndex + 1) % len(c.tanks)
	if !c.ActiveTank().IsAlive() && c.NumberOfTanksAlive() > 0 {
		c.ActivateNextTank()
	}
}

// ActiveTank returns tank which is currently active / on turn.
func (c *Controls) ActiveTank() *Tank {
	return c.tanks[c.activeTankIndex]
}

// NumberOfTanksAlive returns how many tanks is still alive (in game).
func (c *Controls) NumberOfTanksAlive() int {
	alive := 0
	for _, t := range c.tanks {
		if t.IsAlive() {
			alive++
		}
	}
	return alive
}

// Draw does nothing now
func (c *Controls) Draw(s *tl.Screen) {}
