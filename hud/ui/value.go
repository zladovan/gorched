package ui

import (
	"fmt"

	"github.com/zladovan/gorched/draw"
	"github.com/zladovan/gorched/gmath"
)

// Value is component showing some numeric value.
// It has support for changing value and for letting know that value was changed by highlighting it.
//
// Use NewValue to create component with some base value.
// Than call Add to change this value and cause highlight.
//
// Changes are counted separately and if they are non zero value will be highlighted.
// If changes become zero again value highlight will be removed.
type Value struct {
	*BaseComponent
	// Colors defines colors used to draw this component
	Colors         ValueColors
	base, addition int
}

// ValueColors holds colors for Value component
type ValueColors struct {
	// Standard are colors used for unchanged value
	Standard Colors
	// Highlight are colors used for changed value
	Highlight Colors
}

// NewValue creates new Value component.
// Given base value will be considered as unchanged value.
// If you want to change it and cause value to be highlighted call Add with non zero value.
func NewValue(base int) *Value {
	return &Value{
		BaseComponent: &BaseComponent{},
		Colors: ValueColors{
			Standard:  ActivePallette.Standard,
			Highlight: ActivePallette.Highlight,
		},
		base: base,
	}
}

// Dimensions return always X: 5 and Y: 1
func (v *Value) Dimensions() gmath.Vector2i {
	return gmath.Vector2i{X: 5, Y: 1}
}

// Refresh redraws this component to it's canvas
func (v *Value) Refresh() {
	p := draw.BlankPrinter(5, 1)

	// draw box for value, colors are inverted
	p.Bg = v.Colors.Standard.Fg
	p.Fg = v.Colors.Standard.Bg
	p.Write(0, 0, "[   ]")

	// change colors when value was changed and should be highlighted
	if v.addition != 0 {
		if v.Colors.Highlight.Fg != 0 {
			p.Fg = v.Colors.Highlight.Fg
		}
		if v.Colors.Highlight.Bg != 0 {
			p.Bg = v.Colors.Highlight.Bg
		}
	}

	// print value
	p.Write(1, 0, fmt.Sprintf("%3d", v.Get()))

	// change canvas
	v.SetCanvas(p.Canvas)
}

// Get returns current value
func (v *Value) Get() int {
	return v.base + v.addition
}

// Add modifies current value by adding given x to it
func (v *Value) Add(x int) {
	v.addition += x
	v.Refresh()
}

// Addition return change added to base value used to create this component
func (v *Value) Addition() int {
	return v.addition
}
