package ui

import (
	"unicode"

	tl "github.com/JoelOtter/termloop"
)

// Form represents ui window with some components.
//
// It is a connection between ui module and termloop engine.
// It holds components, it can be drawn on screen and it will receive input events.
// It's responsible for drawing components (theirs canvases) to the screen.
type Form interface {
	tl.Drawable
	Parent
	// OnClose will add callback called after form is closed
	OnClose(func())
	// Close form will make this form closed which means it stops to draw itself to screen and process input events
	Close()
	// Closed returns true if this from is already closed
	Closed() bool
}

// BaseForm is basic implementation of form with BaseContainer.
//
// It implements simple focus management.
// Pressing `Tab` will move focus to next component.
// Components are in the same order as they were added to the form.
// Additionally form listens for one of focus keys which can be defined by components.
// Focus management does not support nested containers yet !
//
// You can use it directly to show some components.
// Just call NewForm to create form and than call Add to add components.
//
// You can also use it as base for building other forms.
// If BaseContainer is not enough for you you can change it to something else by calling SetContainer.
type BaseForm struct {
	// entity is used to draw this form on the screen
	entity *tl.Entity
	// container holds form's components
	container Container
	// onClose contains functions which will be called after form is closed
	onClose []func()

	// focusers holds all components which can gain focus
	focusers []Focuser
	// focusIndex is the index of component which currently has focus
	focusIndex int

	refreshed bool
	closed    bool
}

// NewForm crates new form with default style.
func NewForm() *BaseForm {
	f := &BaseForm{}
	f.SetContainer(NewBaseContainer())
	f.Style().Border = Border{Size: 1}
	f.Style().Colors = ActivePallette.Standard
	f.Style().Padding.Left = 1
	f.Style().Padding.Right = 1
	return f
}

// Draw draws form always on the center of the screen
func (f *BaseForm) Draw(s *tl.Screen) {
	if f.closed {
		return
	}
	// redraw to canvas when needed
	if !f.refreshed {
		f.Refresh()
	}
	w, h := f.entity.Size()
	sw, sh := s.Size()
	f.entity.SetPosition(sw/2-w/2, sh/2-h/2)
	f.entity.Draw(s)
}

// Refresh redraws form to it's canvas
func (f *BaseForm) Refresh() {
	f.refreshFocusers()
	f.container.Refresh()
	f.entity = tl.NewEntityFromCanvas(0, 0, *f.container.Canvas())
	f.refreshed = true
}

// refreshFocusers re-collects all components which can gain focus
func (f *BaseForm) refreshFocusers() {
	noFocusersBefore := len(f.focusers) == 0
	f.focusers = []Focuser{}
	for _, c := range f.container.Components() {
		if focuser, ok := c.(Focuser); ok {
			f.focusers = append(f.focusers, focuser)
		}
		// TODO: support for nested containers
	}
	// first component has focus at the start or after focus cleared
	if noFocusersBefore && len(f.focusers) > 0 {
		f.focusIndex = 0
		f.focusers[f.focusIndex].GainFocus()
	}
}

// Tick is processing all input events of this form
func (f *BaseForm) Tick(e tl.Event) {
	// closed form does not process any events
	if f.closed {
		return
	}

	// move focus index
	if e.Key == tl.KeyTab {
		newFocusIndex := f.focusIndex + 1
		if newFocusIndex >= len(f.focusers) {
			newFocusIndex = 0
		}
		f.setFocusIndex(newFocusIndex)
	}

	// process focus keys
	if e.Ch != 0 {
		for fi, focuser := range f.focusers {
			if unicode.ToLower(e.Ch) == unicode.ToLower(focuser.FocusKey()) {
				f.setFocusIndex(fi)
			}
		}
	}

	// process events by currently focused component
	if f.focusIndex < len(f.focusers) {
		f.focusers[f.focusIndex].Tick(e)
	}
}

// SetContainer will change form's container with components
// It will remove all added components and change form's style !
func (f *BaseForm) SetContainer(c Container) {
	f.container = c
	c.AddedToParent(f)
	f.ClearFocus()
}

// NotifyChildChanged is called when any of container's child components is changed
func (f *BaseForm) NotifyChildChanged() {
	f.refreshed = false
}

// Components returns all child components
func (f *BaseForm) Components() []Component {
	return f.container.Components()
}

// Add adds given components to this form
func (f *BaseForm) Add(components ...Component) {
	f.container.Add(components...)
}

// Style can be used to modify form's style
func (f *BaseForm) Style() *Style {
	return f.container.Style()
}

// ClearFocus will remove focus from current focused component.
// It will cause re collecting of focusable components and setting focus to first component.
func (f *BaseForm) ClearFocus() {
	if len(f.focusers) == 0 {
		return
	}
	f.focusers[f.focusIndex].LooseFocus()
	f.focusers = []Focuser{}
	f.Refresh()
}

// Close makes this form closed.
// It will be no more drawn to the screen and it will stop processing input events.
func (f *BaseForm) Close() {
	f.closed = true
	for _, onClose := range f.onClose {
		onClose()
	}
}

// Closed returns true if this from is closed.
func (f *BaseForm) Closed() bool {
	return f.closed
}

// OnClose adds given callback fn to be called after this form will be closed.
// You can call it multiple times to add multiple callbacks.
// Callbacks will be called in same order as they were added.
func (f *BaseForm) OnClose(fn func()) {
	f.onClose = append(f.onClose, fn)
}

// setFocusIndex changes focus to component with index given by i
func (f *BaseForm) setFocusIndex(i int) {
	if f.focusIndex != i {
		f.focusers[f.focusIndex].LooseFocus()
		f.focusers[i].GainFocus()
		f.focusIndex = i
	}
}
