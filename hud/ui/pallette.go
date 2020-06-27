package ui

import tl "github.com/JoelOtter/termloop"

// Pallette holds default colors used in ui components
type Pallette struct {
	// Standard are default colors
	Standard Colors
	// Action are colors used for "actionable" components like buttons
	Action Colors
	// Focus are colors used for components which has focus
	Focus Colors
	// Highlight are colors used to show something more important or changed
	Highlight Colors
}

// Colors hold foreground and background colors
type Colors struct {
	Fg, Bg tl.Attr
}

// DefaultPallette defines colors used as default option
var DefaultPallette = Pallette{
	Standard:  Colors{Fg: tl.Attr(254), Bg: tl.Attr(242)},
	Action:    Colors{Fg: tl.Attr(235) | tl.AttrBold},
	Focus:     Colors{Fg: tl.Attr(254) | tl.AttrBold, Bg: tl.Attr(235)},
	Highlight: Colors{Fg: tl.Attr(125)},
}

// LowColorPallette defines colors used when only 8 colors are available
var LowColorPallette = Pallette{
	Standard:  Colors{Fg: tl.ColorWhite, Bg: tl.ColorBlack},
	Action:    Colors{Fg: tl.ColorWhite | tl.AttrBold},
	Focus:     Colors{Fg: tl.ColorBlack, Bg: tl.ColorWhite},
	Highlight: Colors{Fg: tl.ColorRed},
}

// ActivePallette holds pallette which will be used in ui components.
// Ideally change it before creating any ui components.
// Otherwise you need to recreate all ui components to be applied.
var ActivePallette = DefaultPallette
