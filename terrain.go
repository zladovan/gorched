package gorched

import (
	"math"

	tl "github.com/JoelOtter/termloop"
	osx "github.com/ojrac/opensimplex-go"
	"github.com/zladovan/gorched/draw"
)

// Terrain represents "hills" in the game world.
// Terrain consists of multiple columns with width 1 cell pixel.
// Use NewTerrain to create new instance if you already have terrain line.
// Use GenerateTerrain to create random terrain.
type Terrain struct {
	// columns holds all column entities which is terrain split to
	columns []*TerrainColumn
}

// NewTerrain creates new Terrain for given terrain line and height.
// Terrain line is array where index is x coordinate and value is top y coordinate.
// Terrain height is maximum y value of terrain (lowest on the screen).
func NewTerrain(line []int, height int, lowColor bool) *Terrain {
	terrain := new(Terrain)
	terrain.columns = make([]*TerrainColumn, len(line))

	// create column for each point in terrain line
	for x, baseY := range line {
		p := draw.BlankPrinter(1, height-baseY)

		// print each pixel of column along it's height on canvas
		for y := baseY; y <= height; y++ {
			p.Bg = chooseColor(y-baseY, lowColor)
			p.WritePoint(0, y-baseY, ' ')
		}

		// use canvas to create new column
		terrain.columns[x] = NewTerrainColumn(terrain, x, baseY, p.Canvas)
	}

	return terrain
}

// HeightOn returns y coordinate which will be "on the terrain" for given x
func (t *Terrain) HeightOn(x int) int {
	_, y := t.columns[x].Position()
	return y
}

// PositionOn returns position which will be "on the terrain" for given x
func (t *Terrain) PositionOn(x int) Position {
	return Position{x, t.HeightOn(x)}
}

// Entities returns all entities (columns) which is terrain made of
func (t *Terrain) Entities() []*TerrainColumn {
	return t.columns
}

// CutAround will modify terrain line between x and x+w to be above given y
func (t *Terrain) CutAround(x, y, w int) {
	for i := x; i < x+w; i++ {
		colY := t.HeightOn(i)
		if colY < y {
			t.columns[i].CutFromTop(y - colY)
		}
	}
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

// TerrainGenerator holds configuration for terrain generating logic
type TerrainGenerator struct {
	// Seed to the noise function
	Seed int64
	// Width of the terrain
	Width int
	// Height of the terrain holds maximum height of terrain hill
	Height int
	// Roughness configures how much will be terrain "wavy"
	Roughness float64
	// LowColor generates terrain in only 8 colors mode when true
	LowColor bool
}

// GenerateTerrain will generate new terrain using noise function (open simplex)
func GenerateTerrain(g *TerrainGenerator) *Terrain {
	noise := osx.NewNormalized(g.Seed)
	heights := make([]int, g.Width)
	for x := 0; x < g.Width; x++ {
		// reduce height to keep 5 cells space for tank on the highest hill top
		heights[x] = 5 + int(float64(g.Height-5)*noise.Eval2(g.Roughness/float64(g.Width)*float64(x), 0.5))
	}
	return NewTerrain(heights, g.Height, g.LowColor)
}

// TerrainColumn is collider represented by 1 console pixel wide rectangle with height for it's corresponding part in terrain.
type TerrainColumn struct {
	*tl.Entity
	terrain *Terrain
	canvas  *tl.Canvas
}

// NewTerrainColumn creates new TerrainColumn collider for given terrain.
// Position defined by x and y should be from terrain line.
// Height is the distance between terrain line and the bottom for this column.
func NewTerrainColumn(terrain *Terrain, x, y int, canvas *tl.Canvas) *TerrainColumn {
	return &TerrainColumn{
		Entity:  tl.NewEntityFromCanvas(x, y, *canvas),
		terrain: terrain,
		canvas:  canvas,
	}
}

// Position returns top-left position of collider
func (t *TerrainColumn) Position() (int, int) {
	return t.Entity.Position()
}

// Size returns size of collider
func (t *TerrainColumn) Size() (int, int) {
	return t.Entity.Size()
}

// CutFromTop will cut h pixel cells from top of this column.
func (t *TerrainColumn) CutFromTop(h int) {
	// update canvas
	newCanvas := *t.canvas
	newCanvas[0] = newCanvas[0][h:]
	t.canvas = &newCanvas

	// update entity
	x, y := t.Position()
	t.Entity = tl.NewEntityFromCanvas(x, y+h, newCanvas)
}
