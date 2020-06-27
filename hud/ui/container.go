package ui

import (
	"strings"

	"github.com/zladovan/gorched/draw"
	"github.com/zladovan/gorched/gmath"
)

// Container is component which can hold another components.
// It is responsible for drawing all it's child components to it's canvas.
type Container interface {
	Parent
	Component
	// Style can be used to modify how container looks like
	Style() *Style
}

// Style is collection of attributes affecting how container looks like.
type Style struct {
	// Padding defines size of "space" between container sides and child components
	Padding Padding
	// Border defines border style
	Border Border
	// Colors defines default colors
	Colors Colors
}

// CopyFrom will copy all attributes of given style o to this style
func (s *Style) CopyFrom(o *Style) {
	s.Padding = o.Padding
	s.Border = o.Border
	s.Colors = o.Colors
}

// Border is line around containers
type Border struct {
	// Size can be 0 or 1 for now
	// If size is zero border is disabled
	Size int
	// Colors defines colors of this border
	Colors Colors
}

// Padding defines size of "space" between container sides and child components
type Padding struct {
	Top, Left, Right, Bottom int
}

// PaddingAll creates new padding with same space on all sides
func PaddingAll(a int) Padding {
	return Padding{Top: a, Left: a, Right: a, Bottom: a}
}

// BaseContainer is implementation of Container with size derived from it's child components.
// It should be used when you have exact positions and dimensions of components
// and you want to just create some container around.
// It can also be used as base for building other types of containers.
type BaseContainer struct {
	*BaseComponent
	style Style

	components []Component
}

// NewBaseContainer creates empty container.
// Call Add to add some components.
func NewBaseContainer() *BaseContainer {
	return &BaseContainer{BaseComponent: &BaseComponent{}}
}

// Add will add given components to this containers.
// If this container has some parent it will be notified about change after all components are added.
func (b *BaseContainer) Add(components ...Component) {
	b.components = append(b.components, components...)
	for _, cp := range components {
		cp.AddedToParent(b)
	}
	b.NotifyChildChanged()
}

// Components retruns all components of this container
func (b *BaseContainer) Components() []Component {
	return b.components
}

// Dimensions calculate width and height of this container
func (b *BaseContainer) Dimensions() gmath.Vector2i {
	var maxw, maxh int

	// find maximum x and y coordinates taken by some child component
	for _, c := range b.components {
		p := c.Position()
		d := c.Dimensions()
		fw := p.X + d.X
		if fw > maxw {
			maxw = fw
		}
		fh := p.Y + d.Y
		if fh > maxh {
			maxh = fh
		}
	}

	// apply style to maximum values found
	return gmath.Vector2i{
		X: maxw + b.style.Padding.Left + b.style.Padding.Right + b.style.Border.Size*2,
		Y: maxh + b.style.Padding.Top + b.style.Padding.Bottom + b.style.Border.Size*2,
	}
}

// Refresh redraws container and all it's child components to canvas.
// Child components are not always refreshed before redraw but their are refreshed their canvas is nil (first time refresh).
func (b *BaseContainer) Refresh() {
	// get dimensions
	d := b.Dimensions()
	w := d.X
	h := d.Y

	// create canvas printer
	p := draw.BlankPrinter(w, h)
	p.Bg = b.style.Colors.Bg
	p.Fg = b.style.Colors.Fg
	p.Fill(' ')

	// draw border
	if b.style.Border.Size > 0 {
		p.Bg = b.style.Border.Colors.Bg
		p.Fg = b.style.Border.Colors.Fg

		// print border
		p.Write(0, 0, "╔")
		p.Write(w-1, 0, "╗")
		p.Write(0, h-1, "╚")
		p.Write(w-1, h-1, "╝")
		p.Write(1, 0, strings.Repeat("═", w-2))
		p.Write(1, h-1, strings.Repeat("═", w-2))
		p.WriteHorizontalDown(0, 1, strings.Repeat("║", h-2))
		p.WriteHorizontalDown(w-1, 1, strings.Repeat("║", h-2))
	}

	// draw all components
	offsetX := b.style.Padding.Left + b.style.Border.Size
	offsetY := b.style.Padding.Top + b.style.Border.Size
	for _, c := range b.components {
		// potential first time redraw of this component
		if c.Canvas() == nil {
			c.Refresh()
		}
		p.Draw(c.Position().X+offsetX, c.Position().Y+offsetY, c.Canvas())
	}

	// change canvas
	b.SetCanvas(p.Canvas)
}

// NotifyChildChanged  should be called when one of children is changed
func (b *BaseContainer) NotifyChildChanged() {
	b.BaseComponent.NotifyParentAboutChange()
}

// Style returns style of this container.
// If modified you need to call Refresh to be applied.
func (b *BaseContainer) Style() *Style {
	return &b.style
}
