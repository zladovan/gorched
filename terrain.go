package gorched

import (
	"math"

	tl "github.com/JoelOtter/termloop"
	osx "github.com/ojrac/opensimplex-go"
)

// Terrain represents "hills" in the game world.
type Terrain struct {
	// line is array where index is x coordinate and value is top y coordinate
	line []int
	// height is max height of hill top
	height int
	// lowColor turns only 8 colors mode on
	lowColor bool
}

// Draw draws terrain
func (t *Terrain) Draw(s *tl.Screen) {
	for x, baseY := range t.line {
		for y := baseY; y <= t.height; y++ {
			s.RenderCell(x, y, &tl.Cell{Bg: chooseColor(y-baseY, t.lowColor), Ch: ' '})
		}
	}
}

// Tick is not used yet
func (t *Terrain) Tick(e tl.Event) {}

// GetHeightOn returns y coordinate which will be "on the terrain" for given x
func (t *Terrain) GetHeightOn(x int) int {
	return t.line[x]
}

// GetPositionOn returns position which will be "on the terrain" for given x
func (t *Terrain) GetPositionOn(x int) Position {
	return Position{x, t.GetHeightOn(x)}
}

// GetColliders returns all colliders needed to calculate collisions with terrain
func (t *Terrain) GetColliders() []*TerrainColumn {
	columns := make([]*TerrainColumn, len(t.line))
	for x, baseY := range t.line {
		columns[x] = NewTerrainColumn(t, x, baseY, t.height-baseY)
	}
	return columns
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
}

// GenerateTerrain will generate new terrain using noise function (open simplex)
func GenerateTerrain(g *TerrainGenerator) *Terrain {
	noise := osx.NewNormalized(g.Seed)
	heights := make([]int, g.Width)
	for x := 0; x < g.Width; x++ {
		heights[x] = int(float64(g.Height) * noise.Eval2(g.Roughness/float64(g.Width)*float64(x), 0.5))
	}
	return &Terrain{line: heights, height: g.Height}
}

// TerrainColumn is collider represented by 1 console pixel wide rectangle with height for it's corresponding part in terrain.
type TerrainColumn struct {
	*tl.Entity
	terrain *Terrain
}

// NewTerrainColumn creates new TerrainColumn collider for given terrain.
// Position defined by x and y should be from terrain line.
// Height is the distance between terrain line and the bottom for this column.
func NewTerrainColumn(terrain *Terrain, x, y, height int) *TerrainColumn {
	return &TerrainColumn{
		Entity:  tl.NewEntity(x, y, 1, height),
		terrain: terrain,
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
