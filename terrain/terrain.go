package terrain

import (
	"math"

	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/debug"
	"github.com/zladovan/gorched/draw"
	"github.com/zladovan/gorched/gmath"
)

// Terrain represents "hills" in the game world.
// Terrain consists of multiple columns with width 1 cell pixel.
// There can be no column or many columns for each x coordinate.
// Use NewTerrain to create new instance if you already have terrain line.
// Use GenerateTerrain to create random terrain.
type Terrain struct {
	// height is max height of hill top
	height int
	// columns holds all column entities which is terrain split to
	columns [][]*Column
	// cutter provides terrain destruction logic
	cutter *Cutter
}

// NewTerrain creates new Terrain for given terrain line and height.
// Terrain line is array where index is x coordinate and value is top y coordinate.
// Terrain height is maximum y value of terrain (lowest on the screen).
func NewTerrain(line []int, height int, lowColor bool) *Terrain {
	terrain := &Terrain{height: height}
	terrain.columns = make([][]*Column, len(line))
	terrain.cutter = &Cutter{terrain: terrain}

	// create column for each point in terrain line
	for x, baseY := range line {
		p := draw.BlankPrinter(1, height-baseY)

		// print each pixel of column along it's height on canvas
		for y := baseY; y <= height; y++ {
			p.Bg = chooseColor(y-baseY, lowColor)
			p.WritePoint(0, y-baseY, ' ')
		}

		// use canvas to create new column
		terrain.columns[x] = []*Column{NewColumn(terrain, x, baseY, p.Canvas)}
	}

	return terrain
}

// HeightOn returns y coordinate which will be "on the terrain" for given x
func (t *Terrain) HeightOn(x int) int {
	return t.HeightInside(x, 0)
}

// HeightInside returns y coordinate which will be "in the terrain" on the nearest column under given y for given x
// It allows to find y coordinate inside terrain hole.
func (t *Terrain) HeightInside(x, y int) int {
	if len(t.columns[x]) == 0 {
		return t.height
	}
	for _, c := range t.columns[x] {
		_, cy := c.Position()
		if y <= cy {
			return cy
		}
	}
	return t.height
}

// PositionOn returns position which will be "on the terrain" for given x
func (t *Terrain) PositionOn(x int) gmath.Vector2i {
	return gmath.Vector2i{X: x, Y: t.HeightOn(x)}
}

// Entities returns all entities (columns) which is terrain made of
func (t *Terrain) Entities() []tl.Drawable {
	entities := []tl.Drawable{t.cutter}
	for _, cs := range t.columns {
		for _, c := range cs {
			entities = append(entities, c)
		}
	}
	return entities
}

// CutAround will modify terrain line between x and x+w to be above given y
func (t *Terrain) CutAround(x, y, w int) {
	for i := x; i < x+w; i++ {
		if len(t.columns[i]) == 0 {
			continue
		}
		colY := t.HeightOn(i)
		if colY < y {
			t.cutter.CutFromTop(i, y-colY)
		}
	}
}

// MakeHole will create hole in terrain with center at cx and cy coordinates with given radius r.
func (t *Terrain) MakeHole(cx, cy, r int) {
	debug.Logf("Hole in the terrain centerx=%d, centery=%d", cx, cy)
	t.cutter.CutHole(cx, cy, r)
}

// Line returns terrain line array where index is x coordinate and value is top y coordinate.
func (t *Terrain) Line() []int {
	line := make([]int, len(t.columns))
	for x := range t.columns {
		line[x] = t.HeightOn(x)
	}
	return line
}

// palette holds colors of terrain - shades of green
var palette = []int{41, 35, 29, 23}

// chooseColor selects color from palette for given height
// if lowColor mode is on it will return just one color all the time
// otherwise it will create gradient with each deeper level wider than level above
func chooseColor(height int, lowColor bool) tl.Attr {
	if lowColor {
		return tl.ColorGreen
	}
	idx := int(math.Sqrt(float64(height+1))) - 1
	if idx >= len(palette) {
		idx = len(palette) - 1
	}
	return tl.Attr(palette[idx])
}
