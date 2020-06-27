package ui

import (
	tl "github.com/JoelOtter/termloop"
)

// MessageBox shows form with some message in the screen center.
type MessageBox struct {
	*BaseForm
}

// NewMessageBox creates new message box from given string with multiple lines.
// Create MessageBox will have dimensions based on the lines resolved from given string.
func NewMessageBox(msg string) *MessageBox {
	form := &MessageBox{BaseForm: NewForm()}
	form.Add(NewText(msg))
	return form
}

// Tick handles controls for this MessageBox
func (m *MessageBox) Tick(e tl.Event) {
	// Message box is closed on any key press
	if e.Type == tl.EventKey {
		m.Close()
	}
}
