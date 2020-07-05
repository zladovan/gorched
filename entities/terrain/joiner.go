package terrain

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/debug"
)

// Joiner is responsible for joining terrain columns.
// If two columns with same x position "touch" themselves and there is no hole between them they can be joined to one.
// Joiner is by default not enabled. You need to activate it by calling Enable. Then it will be enabled for 1 second.
type Joiner struct {
	// terrain is reference to the terrain which columns will be joined
	terrain *Terrain
	// ttl is number of seconds till joiner will be not active
	ttl float64

}

// Draw will perform joining logic if this joiner is enabled
func (j *Joiner) Draw(s *tl.Screen) {
	// early exit if not enabled
	if j.ttl <= 0 {
		return
	}
	j.ttl -= s.TimeDelta()

	for x, columns := range j.terrain.columns {
		// nothing to join
		if len(columns) <= 1 {
			continue
		}

		joined := []*Column{columns[0]}
		for _, column := range columns[1:] {
			last := joined[len(joined)-1]
			_, ly := last.Position()
			_, lh := last.Size()
			_, y := column.Position()

			// no need to join
			if ly + lh < y {
				joined = append(joined, column)
				continue
			}

			// joining
			// canvas of last column is added to the canvas of current colum
			debug.Logf("Joining terrain columns x=%d y1=%d y2=%d", x, ly, lh)
			column.SetCanvas(&tl.Canvas{append((*last.canvas)[0], (*column.canvas)[0]...)})
			
			// replace and remove last column
			joined[len(joined) - 1] = column
			s.Level().RemoveEntity(last)
		}

		// replace original columns with joined
		j.terrain.columns[x] = joined
	}
}

// Tick does nothing now
func (j *Joiner) Tick(e tl.Event) {}

// Enable this joiner for one second
func (j *Joiner) Enable() {
	j.ttl = 1
}
