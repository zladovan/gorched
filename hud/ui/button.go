package ui

import (
	"strings"
	"unicode"

	tl "github.com/JoelOtter/termloop"

	"github.com/zladovan/gorched/draw"
	"github.com/zladovan/gorched/gmath"
)

// Button is component represent by text label.
// It can gain focus and "be pressed" to call some action.
//
// Button is pressed if:
//
//  - it has focus and SPACE or ENTER is hit
//  - it has ActionKey defined and it was hit
//
// To create new button call NewButton.
// Then set button's attributes and add it to some container.
//
// When you change attributes after Button is added to the container, call Refresh.
type Button struct {
	// it extends from BaseComponent
	*BaseComponent
	// Colors defines colors of the button
	Colors ButtonColors
	// ActionKey is character which when typed this button will gain focus and call it's action.
	// It is case insensitive.
	// If button's text contains this character it's first occurrence will be underlined.
	ActionKey rune
	// Action will be called when button is pressed
	Action func()
	// Text is string representation of this button and it will be drawn on component's position
	Text string

	// focus is flag defining if this button has focus now
	focus bool
}

// ButtonColors holds colors for all button states
type ButtonColors struct {
	// Standard are colors used when button is out of focus
	Standard Colors
	// Focus are colors used when button has focus
	Focus Colors
}

// NewButton creates button with given text and action to be called when pressed.
func NewButton(text string, action func()) *Button {
	return &Button{
		BaseComponent: &BaseComponent{},
		Colors: ButtonColors{
			Standard: Colors{Fg: ActivePallette.Standard.Fg | tl.AttrBold},
			Focus:    ActivePallette.Focus,
		},
		Text:   text,
		Action: action,
	}
}

// Dimensions returns width and height of this button.
// Buttons width depends on it's text and height is always 1 cell pixel.
func (b *Button) Dimensions() gmath.Vector2i {
	return gmath.Vector2i{X: len(b.Text), Y: 1}
}

// GainFocus will add focus to this button
func (b *Button) GainFocus() {
	b.focus = true
	b.Refresh()
}

// LooseFocus will remove focus from this button
func (b *Button) LooseFocus() {
	b.focus = false
	b.Refresh()
}

// Tick handles inputs for this button.
func (b *Button) Tick(e tl.Event) {
	// nothing to handle if button is out of focus
	if !b.focus {
		return
	}
	// calling action
	switch e.Key {
	case tl.KeySpace, tl.KeyEnter:
		b.Action()
		return
	}
	if e.Ch != 0 && unicode.ToLower(e.Ch) == unicode.ToLower(b.ActionKey) {
		b.Action()
	}
}

// FocusKey defines character which when typed should bring focus to this button
func (b *Button) FocusKey() rune {
	return b.ActionKey
}

// Refresh will redraw this button to it's canvas
func (b *Button) Refresh() {
	p := draw.BlankPrinter(len([]rune(b.Text)), 1)
	p.Fg = b.Colors.Standard.Fg
	p.Bg = b.Colors.Standard.Bg
	if b.focus {
		p.Fg = b.Colors.Focus.Fg
		p.Bg = b.Colors.Focus.Bg
	}
	p.Write(0, 0, b.Text)
	b.highlightActionKey(p.Canvas)
	b.SetCanvas(p.Canvas)
}

// underline action key character in button's text if defined
func (b *Button) highlightActionKey(canvas *tl.Canvas) {
	if b.ActionKey == 0 {
		return
	}
	ai := strings.IndexRune(b.Text, b.ActionKey)
	if ai == -1 {
		return
	}
	(*canvas)[0][ai].Fg |= tl.AttrUnderline
}
