package gorched

import (
	"math"

	tl "github.com/JoelOtter/termloop"
	osx "github.com/ojrac/opensimplex-go"
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
	columns [][]*TerrainColumn
}

// NewTerrain creates new Terrain for given terrain line and height.
// Terrain line is array where index is x coordinate and value is top y coordinate.
// Terrain height is maximum y value of terrain (lowest on the screen).
func NewTerrain(line []int, height int, lowColor bool) *Terrain {
	terrain := &Terrain{height: height}
	terrain.columns = make([][]*TerrainColumn, len(line))

	// create column for each point in terrain line
	for x, baseY := range line {
		p := draw.BlankPrinter(1, height-baseY)

		// print each pixel of column along it's height on canvas
		for y := baseY; y <= height; y++ {
			p.Bg = chooseColor(y-baseY, lowColor)
			p.WritePoint(0, y-baseY, ' ')
		}

		// use canvas to create new column
		terrain.columns[x] = []*TerrainColumn{NewTerrainColumn(terrain, x, baseY, p.Canvas)}
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
func (t *Terrain) PositionOn(x int) Position {
	return Position{x, t.HeightOn(x)}
}

// Entities returns all entities (columns) which is terrain made of
func (t *Terrain) Entities() []*TerrainColumn {
	cols := []*TerrainColumn{}
	for _, cs := range t.columns {
		cols = append(cols, cs...)
	}
	return cols
}

// CutAround will modify terrain line between x and x+w to be above given y
func (t *Terrain) CutAround(x, y, w int) {
	for i := x; i < x+w; i++ {
		if len(t.columns[i]) == 0 {
			continue
		}
		colY := t.HeightOn(i)
		if colY < y {
			t.columns[i][0].CutFromTop(y - colY)
		}
	}
}

// MakeHole will create hole in terrain with center at cx and cy coordinates with given radius r.
func (t *Terrain) MakeHole(cx, cy, r int) {
	Debug.Logf("Hole in the terrain centerx=%d, centery=%d", cx, cy)
	for ix := -r + 1; ix < r; ix++ {
		iy := int(math.Sqrt(math.Pow(float64(r-1), 2) - math.Pow(float64(ix), 2)))
		miny := cy - iy
		maxy := cy + iy
		x := cx + ix
		for _, c := range t.columns[x] {
			c.Cut(miny, maxy)
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
	body    *Body
	terrain *Terrain
	canvas  *tl.Canvas
	// isCut if is true this column will be removed and replaced by columns from cuttingParts array on next frame Draw
	isCut        bool
	cuttingParts []*TerrainColumn
	bodyLocker   *TimeLocker
}

// NewTerrainColumn creates new TerrainColumn collider for given terrain.
// Position defined by x and y should be from terrain line.
func NewTerrainColumn(terrain *Terrain, x, y int, canvas *tl.Canvas) *TerrainColumn {
	body := &Body{
		Position: gmath.Vector2f{X: float64(x), Y: float64(y + len((*canvas)[0]))},
		Mass:     5,
		Locked:   true,
	}
	return &TerrainColumn{
		Entity:     tl.NewEntityFromCanvas(x, y, *canvas),
		body:       body,
		terrain:    terrain,
		canvas:     canvas,
		bodyLocker: &TimeLocker{BodyToRelock: body, RemainingSeconds: 0.25},
	}
}

// Draw draws this column and process possible cuttings
func (t *TerrainColumn) Draw(s *tl.Screen) {
	// update body locker
	t.bodyLocker.Update(s.TimeDelta())

	// update position of entity based on body position if not locked
	if !t.body.Locked {
		t.Entity.SetPosition(int(t.body.Position.X), int(t.body.Position.Y)-len((*t.canvas)[0]))
	}

	// draw entity
	t.Entity.Draw(s)

	// process cut
	if t.isCut {
		x := int(t.body.Position.X)

		// find index of this column in all columns with position x
		idx := 0
		for i, c := range t.terrain.columns[x] {
			if c == t {
				idx = i
				break
			}
		}

		// create new columns with all columns before this column
		newcols := []*TerrainColumn{}
		newcols = append(newcols, t.terrain.columns[x][:idx]...)

		// add columns created during cut to new columns and to the level
		for _, p := range t.cuttingParts {
			newcols = append(newcols, p)
			s.Level().AddEntity(p)
		}

		// add rest of columns after this column
		newcols = append(newcols, t.terrain.columns[x][idx+1:]...)

		// update terrain columns
		t.terrain.columns[x] = newcols

		// remove this columns from level
		s.Level().RemoveEntity(t)
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
// It has immediate effect in contrast to Cut.
func (t *TerrainColumn) CutFromTop(cells int) {
	// update canvas
	newCanvas := *t.canvas
	newCanvas[0] = newCanvas[0][cells:]
	t.canvas = &newCanvas

	// update entity
	x, y := t.Position()
	//t.body.Position.Y += float64(cells)
	t.Entity = tl.NewEntityFromCanvas(x, y+cells, newCanvas)
}

// Cut cuts horizontal line given by miny and maxy from this column.
// Cutting can result to removing this column and to replacing it with zero, one or two new columns.
// Number of new columns depends on the position of the intersection of line and column.
// Effects of cut are applied on next frame Draw.
func (t *TerrainColumn) Cut(miny, maxy int) {
	// get dimensions
	x, y := t.Position()
	_, h := t.Size()

	// is cut line out of column ?
	if maxy < y || miny > y+h-1 {
		return
	}

	// init cut
	t.isCut = true
	t.cuttingParts = []*TerrainColumn{}

	// local y coordinates of cut hole
	topy := int(math.Max(0, float64(miny-y)))
	bottomy := int(math.Min(float64(h-1), float64(maxy-y)))

	// create two new columns around cut hole
	canvas := *t.canvas
	topCol := NewTerrainColumn(t.terrain, x, y, &tl.Canvas{canvas[0][:topy]})
	bottomCol := NewTerrainColumn(t.terrain, x, y+bottomy+1, &tl.Canvas{canvas[0][bottomy+1:]})

	// add top column if it has non zero height
	_, h1 := topCol.Size()
	if h1 > 0 {
		t.cuttingParts = append(t.cuttingParts, topCol)
	}

	// add bottom column if it has non zero height
	_, h2 := bottomCol.Size()
	if h2 > 0 {
		t.cuttingParts = append(t.cuttingParts, bottomCol)
	}
}

// Body returns physical body for falling processing
func (t *TerrainColumn) Body() *Body {
	return t.body
}

// BottomLine returns single point at column's body position for collision with the ground when falling
func (t *TerrainColumn) BottomLine() (int, int) {
	return 0, 0
}
