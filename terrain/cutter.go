package terrain

import (
	"math"

	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/gmath"
)

// Cutter is entity which performs cutting of terrain columns.
// It is able to update terrain columns
type Cutter struct {
	terrain *Terrain
	cuts    []Cut
}

// Cut represents horizontal line on given X coordinate going from MinY to MaxY which should be cut from the terrain column.
type Cut struct{ X, MinY, MaxY int }

// CutHole will create hole in terrain with center at cx and cy coordinates with given radius r.
func (c *Cutter) CutHole(cx, cy, r int) {
	for ix := -r + 1; ix < r; ix++ {
		// y coordinate is scaled by 0.5 to reduce terminal's cells ratio 2:1 for height:width
		iy := int(math.Sqrt(math.Pow(float64(r-1), 2)-math.Pow(float64(ix), 2)) * 0.5)
		miny := cy - iy
		maxy := cy + iy
		x := cx + ix
		if x < 0 || x >= len(c.terrain.columns) {
			continue
		}
		c.Cut(x, miny, maxy)
	}
}

// Cut will cut column at given x by horizontal line going from miny to maxy.
// Cutting can result to removing column and to replacing it with zero, one or more new columns.
// Number of new columns depends on the position of the intersection of line and column.
// Effects of Cut will be applied on nex frame Draw.
func (c *Cutter) Cut(x, miny, maxy int) {
	c.cuts = append(c.cuts, Cut{X: x, MinY: miny, MaxY: maxy})
}

// CutFromTop will cut h pixel cells from top column on given x.
// It has immediate effect in contrast to Cut.
func (c *Cutter) CutFromTop(x, cells int) {
	// skip empty columns
	if len(c.terrain.columns[x]) == 0 {
		return
	}

	// get the top column on x
	column := c.terrain.columns[x][0]

	// update canvas
	newCanvas := *column.canvas
	newCanvas[0] = newCanvas[0][cells:]
	column.canvas = &newCanvas

	// update entity
	x, y := column.Position()
	//t.body.Position.Y += float64(cells)
	column.Entity = tl.NewEntityFromCanvas(x, y+cells, newCanvas)
}

// Draw is processing all pending cuts
func (c *Cutter) Draw(s *tl.Screen) {
	// process all pending cuts
	for _, cut := range c.cuts {

		// for each x which is should be cut we will create new columns updated by current cut
		newcols := []*Column{}
		for _, column := range c.terrain.columns[cut.X] {

			// do the cutting logic on one column
			cuttingParts, isCut := doCut(column, cut)

			// no cut means that column can be just moved to new columns and we can go to next column
			if !isCut {
				newcols = append(newcols, column)
				continue
			}

			// add columns created during cut to new columns and to the level
			for _, p := range cuttingParts {
				newcols = append(newcols, p)
				s.Level().AddEntity(p)
			}

			// remove this column from level as it was replaced by new columns or cut off
			s.Level().RemoveEntity(column)
		}
		// replace old columns with new updated columns
		c.terrain.columns[cut.X] = newcols
	}

	// clear cuts as they were already processed
	c.cuts = []Cut{}
}

// Tick does nothing now
func (c *Cutter) Tick(e tl.Event) {}

// doCut performs cutting Column t by Cut c.
// It returns new columns created with this cut and boolean flag if there was some cut or not.
// If there was cut new columns can be also empty which means that whole column was destroyed by this cut.
func doCut(t *Column, c Cut) ([]*Column, bool) {
	// get dimensions
	x, y := t.Position()
	_, h := t.Size()

	// is cut line out of column ?
	if c.MaxY < y || c.MinY > y+h-1 {
		return nil, false
	}

	// init cut
	cuttingParts := []*Column{}

	// local y coordinates of cut hole
	topy := gmath.Max(0, c.MinY-y)
	bottomy := gmath.Min(h-1, c.MaxY-y)

	// create two new columns around cut hole
	canvas := *t.canvas
	topCol := NewColumn(t.terrain, x, y, &tl.Canvas{canvas[0][:topy]})
	bottomCol := NewColumn(t.terrain, x, y+bottomy+1, &tl.Canvas{canvas[0][bottomy+1:]})

	// add top column if it has non zero height
	_, h1 := topCol.Size()
	if h1 > 0 {
		cuttingParts = append(cuttingParts, topCol)
	}

	// add bottom column if it has non zero height
	_, h2 := bottomCol.Size()
	if h2 > 0 {
		cuttingParts = append(cuttingParts, bottomCol)
	}

	return cuttingParts, true
}
