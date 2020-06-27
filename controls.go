package gorched

import (
	tl "github.com/JoelOtter/termloop"
)

// Controls holds data and logic for controlling game world.
type Controls struct {
	// reference to game
	game *Game
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

	// when message box is shown it is in control
	if c.game.Hud().IsFormShown() {
		return
	}

	// otherwise handle in-game controls
	switch e.Key {
	case tl.KeyArrowLeft:
		c.MoveUp()
	case tl.KeyArrowRight:
		c.MoveDown()
	case tl.KeySpace:
		c.Shoot()
	case tl.KeyCtrlR:
		c.RestartRound()
	case tl.KeyCtrlN:
		c.NextRound()
	}
	switch e.Ch {
	case 'h':
		c.ShowInfo()
	case 's':
		c.ShowScore()
	case 'a':
		c.ShowAttributes()
	}
	// for the browser mode we cannot use ctrl+n and ctr+r as we would leave the window
	if c.game.options.BrowserMode {
		switch e.Ch {
		case 'n':
			c.NextRound()
		case 'r':
			c.RestartRound()
		}
	}
}

// Draw does nothing now
func (c *Controls) Draw(s *tl.Screen) {}

// Following methods can be also used from outside as a support for other external controller

// MoveUp increase cannon's angle of active tank
func (c *Controls) MoveUp() {
	if c.game.round.IsPlayerOnTurn() {
		c.game.round.ActiveTank().MoveUp()
	}
}

// MoveDown decreases cannon's angle of active tank
func (c *Controls) MoveDown() {
	if c.game.round.IsPlayerOnTurn() {
		c.game.round.ActiveTank().MoveDown()
	}
}

// Shoot will start loading or shoot with active tank if it's already loading
func (c *Controls) Shoot() {
	if c.game.round.IsPlayerOnTurn() {
		c.game.round.ActiveTank().Shoot()
	}
}

// HideMessageBox will hide any active message box
func (c *Controls) HideMessageBox() {
	c.game.Hud().HideForm()
}

// NextRound will switch game to the next round
func (c *Controls) NextRound() {
	c.game.round.Next()
}

// RestartRound will restart current round
func (c *Controls) RestartRound() {
	c.game.round.Restart()
}

// ShowInfo shows main game information
func (c *Controls) ShowInfo() {
	c.game.Hud().ShowInfo()
}

// ShowScore shows actual score board
func (c *Controls) ShowScore() {
	c.game.Hud().ShowScore()
}

// Show attributes dialog
func (c *Controls) ShowAttributes() {
	c.game.hud.ShowAttributes(true)
}
