package hud

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/core"
	"github.com/zladovan/gorched/hud/ui"
)

// HUD stands for heads-up display and it holds entities which are drawn on  the screen always over all the level entities.
type HUD struct {
	// game reference
	game core.Game
	// options for HUD
	options Options
	// form holds some ui form which can be displayed in the center of the screen
	// there is always only one form shown at the same time
	form ui.Form
	// skipTick if true will cause Tick not processed until next frame redrawn
	// this is needed to avoid closing message boxes right after their are shown
	skipTick bool
}

// Options holds flags affecting how HUD should look like
type Options struct {
	// AsciiOnly identifies that only ASCII characters can be used for HUD graphics
	ASCIIOnly bool
	// LowColor identifies that only 8 colors can be used for HUD graphics
	LowColor bool
	// BrowserMode identifies that game was run in browser and some controls need to be modified to do not collide with usual browser shortcuts
	BrowserMode bool
}

// NewHUD creates new HUD for given game
func NewHUD(game core.Game, options Options) *HUD {
	// switch ui pallette for low color mode if needed
	if options.LowColor {
		ui.ActivePallette = ui.LowColorPallette
	}
	return &HUD{game: game, options: options}
}

// IsFormShown informs wether some message box is currently displayed
func (h *HUD) IsFormShown() bool {
	return h.form != nil && !h.form.Closed()
}

// ShowForm shows given box on the center of the screen
func (h *HUD) ShowForm(box ui.Form) {
	h.HideForm()
	box.OnClose(func() {
		h.form = nil
	})
	h.form = box
	h.skipTick = true
}

// HideForm hides any visible message box on the screen
func (h *HUD) HideForm() {
	if h.form != nil {
		h.form.Close()
		h.form = nil
	}
}

// ShowInfo shows message box with main game information
func (h *HUD) ShowInfo() *ui.MessageBox {
	info := NewInfoBox(h.options.BrowserMode)
	h.ShowForm(info)
	return info
}

// ShowScore shows message box with actual score
func (h *HUD) ShowScore() *ui.MessageBox {
	score := NewScoreBox(h.game.Players())
	h.ShowForm(score)
	return score
}

// Draw draws all entities of HUD
func (h *HUD) Draw(s *tl.Screen) {
	// reset skipTick flag if needed
	if h.skipTick {
		h.skipTick = false
	}
	// no form means nothing to draw now
	if h.form == nil {
		return
	}
	// draw active form
	h.form.Draw(s)
}

// Tick does nothing now
func (h *HUD) Tick(e tl.Event) {
	// prevents auto closing of message boxes
	if h.skipTick {
		return
	}
	if h.form != nil {
		h.form.Tick(e)
	}
}
