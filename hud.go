package gorched

import (
	tl "github.com/JoelOtter/termloop"
)

// HUD stands for heads-up display and it holds entities which are drawn on  the screen always over all the level entities.
type HUD struct {
	// game reference
	game *Game
	// messageBox holds some message box which can be displayed in the center of the screen
	messageBox *MessageBox
}

// NewHUD creates new HUD for given game
func NewHUD(game *Game) *HUD {
	return &HUD{game: game}
}

// IsMessageBoxShown informs wether some message box is currently displayed
func (h *HUD) IsMessageBoxShown() bool {
	return h.messageBox != nil
}

// ShowMessageBox shows given box on the center of the screen
func (h *HUD) ShowMessageBox(box *MessageBox) {
	h.messageBox = box
}

// HideMessageBox hides any visible message box on the screen
func (h *HUD) HideMessageBox() {
	h.messageBox = nil
}

// ShowInfo shows message box with main game information
func (h *HUD) ShowInfo() {
	h.ShowMessageBox(NewInfoBox(h.game.options.BrowserMode, h.game.options.LowColor))
}

// ShowScore shows message box with actual score
func (h *HUD) ShowScore() {
	h.ShowMessageBox(NewScoreBox(h.game))
}

// Draw draws all entities of HUD
func (h *HUD) Draw(s *tl.Screen) {
	if h.messageBox != nil {
		h.messageBox.Draw(s)
	}
}

// Tick does nothing now
func (h *HUD) Tick(e tl.Event) {}
