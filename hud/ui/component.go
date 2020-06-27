package ui

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/gmath"
)

// Component is base building block of ui.
// It need to define it's position and dimensions (width on X, height on Y).
// It need to provide canvas where it is drawn.
// It can be added to some parent and it will be notified about it by calling AddedToParent function.
type Component interface {
	// Position returns position of this component relative to it's parent postion
	Position() *gmath.Vector2i
	// Dimensions returns width and height of this component (width on X, height on Y)
	Dimensions() gmath.Vector2i
	// Canvas returns canvas where this component is drawn
	Canvas() *tl.Canvas
	// AddedToParent is called after this component is added to some parent
	AddedToParent(p Parent)
	// Refresh should cause redrawing itself to it's canvas
	Refresh()
}

// Focuser should be used to extend some component with possibility to gain focus.
// It will be notified about focus gain and loose.
// When it has focus it will start to receive input events by calling Tick method for each event.
type Focuser interface {
	// GainFocus is called when focus was changed and added to this focuser
	GainFocus()
	// LoosFocus is called when focus was changed and removed from this focuser
	LooseFocus()
	// FocusKey defines character which when typed this focuser will gain focus.
	FocusKey() rune
	// Tick is called for each event when focuser has focus
	Tick(e tl.Event)
}

// Parent holds one or more child components.
// It provides support for creating hierarchy of components.
// But Parent itself does not need to be component.
//
// It should provide all it's components and allow add new components.
// It will be notified about changes in one of its' children by calling NotifyChildChanged function
type Parent interface {
	// Components returns all components having this parent
	Components() []Component
	// Add adds component to this parent
	Add(components ...Component)
	// NotifyChildChanged should be called when one of children is changed
	NotifyChildChanged()
}

// BaseComponent is empty component with implementation of all functions required by Component.
// It's purpose is to be wrapped by another component which then does not need to implement all functions.
//
// Wrapper component then need just implement Refresh function.
// Wrapper should call SetCanvas to store canvas changes usually at the end of Refresh function.
// Wrapper could also call NotifyParentAboutChange to let know it's parent that it's canvas was changed and it need to be redrawn.
type BaseComponent struct {
	position gmath.Vector2i
	canvas   *tl.Canvas
	parent   Parent
}

// Position returns position of this component
func (c *BaseComponent) Position() *gmath.Vector2i {
	return &c.position
}

// Dimensions returns size of component's canvas
func (c *BaseComponent) Dimensions() gmath.Vector2i {
	if c.canvas == nil {
		return gmath.Vector2i{X: 0, Y: 0}
	}
	var w, h int
	w = len(*c.canvas)
	if w > 0 {
		h = len((*c.canvas)[0])
	}
	return gmath.Vector2i{X: w, Y: h}
}

// Canvas returns canvas where this component is drawn
func (c *BaseComponent) Canvas() *tl.Canvas {
	return c.canvas
}

// SetCanvas will change canvas of thi component and it will notify parent about it.
func (c *BaseComponent) SetCanvas(canvas *tl.Canvas) {
	c.canvas = canvas
	c.NotifyParentAboutChange()
}

// AddedToParent will store given parent.
// This parent will be later notified about each canvas change
func (c *BaseComponent) AddedToParent(p Parent) {
	c.parent = p
}

// NotifyParentAboutChange will notify stored parent about canvas changes
func (c *BaseComponent) NotifyParentAboutChange() {
	if c.parent != nil {
		c.parent.NotifyChildChanged()
	}
}

// Refresh does nothing
func (c *BaseComponent) Refresh() {}
