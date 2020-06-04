package gorched

import (
	tl "github.com/JoelOtter/termloop"
)

// TODO: cleanup, split and move to more meaningfully named files / packages

// DrawPallette draws all available colors as column with 6 colors per line.
// Given shift will change starting color and will result in printing less than 256 colors.
// Debug only !
func DrawPallette(s *tl.Screen, shift int) {
	for c := 0; c < 256-shift; c++ {
		x := c % 6
		y := c / 6
		s.RenderCell(x, y, &tl.Cell{Fg: tl.Attr(c + shift), Ch: '█'})
	}
}
