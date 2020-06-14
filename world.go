package gorched

import (
	"math/rand"
	"sort"

	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/debug"
	"github.com/zladovan/gorched/physics"
	"github.com/zladovan/gorched/terrain"
)

// World represents game world with all entities.
// It extends from termloop.BaseLevel so it can be added to the screen as termloop.Level.
type World struct {
	*tl.BaseLevel
	terrain *terrain.Terrain
	physics *physics.Physics
	options WorldOptions
	// entitiesToRemove holds references to entities which will be removed on next Tick
	entitiesToRemove []tl.Drawable
}

// WorldOptions provide configuration needed for generating game world (one round).
type WorldOptions struct {
	// Width of game world in number of console pixels (cells)
	Width int
	// Height of game world in number of console pixels (cells)
	Height int
	// Seed is number used as random seed and if it is reused it allows to create same game looking game world with same positions for players
	Seed int64
	// AsciiOnly identifies that only ASCII characters can be used for world graphics
	ASCIIOnly bool
	// LowColor identifies that only 8 colors can be used for world graphics
	LowColor bool
}

// NewWorld creates new game world with all entities
func NewWorld(game *Game, o WorldOptions) *World {
	// random positions in the world are seeded too
	rnd := rand.New(rand.NewSource(o.Seed))

	// create terrain
	terrain := terrain.Generate(&terrain.Generator{
		Seed:      o.Seed,
		Width:     o.Width,
		Height:    o.Height,
		Roughness: 7.5,
		LowColor:  game.options.LowColor,
	})

	// create clouds
	clouds := GenerateClouds(&CloudsGenerator{seed: o.Seed, width: o.Width, height: o.Height})

	// create trees
	trees := GenerateWood(&WoodGenerator{
		Line:      terrain.Line(),
		Seed:      o.Seed,
		Density:   0.2,
		MaxSize:   6,
		MinSpace:  1,
		LowColor:  game.options.LowColor,
		ASCIIOnly: game.options.ASCIIOnly,
	})

	// create players
	// TODO: update for different player counts
	tanks := []*Tank{
		NewTank(
			game.players[0],
			terrain.PositionOn(10+rnd.Intn(10)),
			0,
			tl.ColorRed,
			game.options.ASCIIOnly,
		),
		NewTank(
			game.players[1],
			terrain.PositionOn(o.Width-10-rnd.Intn(10)),
			180,
			tl.ColorBlack,
			game.options.ASCIIOnly,
		),
	}

	// cut the trees and terrain around the tanks
	for _, tank := range tanks {
		x, y := tank.Position()
		w, h := tank.Size()
		trees = trees.CutAround(x, y, w, h)
		terrain.CutAround(x, y+h, w)
	}

	// create controls
	controls := &Controls{
		game:            game,
		tanks:           tanks,
		activeTankIndex: game.startingPlayerIndex,
	}

	// create level with all entities
	bg := tl.Attr(111)
	if game.options.LowColor {
		bg = tl.ColorBlue
	}
	world := &World{
		BaseLevel: tl.NewBaseLevel(tl.Cell{Bg: bg}),
		terrain:   terrain,
		physics:   &physics.Physics{Gravity: 9.81, Ground: terrain.HeightInside},
		options:   o,
	}
	world.AddEntity(clouds)
	for _, c := range terrain.Entities() {
		world.AddEntity(c)
	}
	for _, t := range trees {
		world.AddEntity(t)
	}
	for _, t := range tanks {
		world.AddEntity(t)
	}
	world.AddEntity(controls)

	debug.Logf("New world created width=%d height=%d seed=%d", o.Width, o.Height, o.Seed)

	return world
}

// RemoveEntity only registers entity to remove.
// Entity will be removed in next Tick.
// This is needed for be able to remove entities from Draw method (where level is accessible).
func (w *World) RemoveEntity(e tl.Drawable) {
	w.entitiesToRemove = append(w.entitiesToRemove, e)
}

// Draw draws all entities in the world.
// Drawing takes into account z-index of entity which can be specified by implementing ZIndexer interface.
// Gravity is also applied here as there is access to delta time from screen.
func (w *World) Draw(s *tl.Screen) {
	// depthField  will contain entities distributed by z-index
	depthField := make(map[int][]tl.Drawable, len(w.Entities))

	for _, e := range w.BaseLevel.Entities {

		// apply physics to all entities with bodies
		if entity, ok := e.(physics.HasBody); ok {
			w.physics.Apply(entity, s.TimeDelta())
		}

		// find z-index of entity where 0 is the default
		zIndex := 0
		if entity, ok := e.(ZIndexer); ok {
			zIndex = entity.ZIndex()
		}

		// add entity to depth field
		depthField[zIndex] = append(depthField[zIndex], e)
	}

	// find all different values of z-index and sort them
	zIndexes := []int{}
	for z := range depthField {
		zIndexes = append(zIndexes, z)
	}
	sort.Ints(zIndexes)

	// draw entities in order defined by z-index (lower first)
	for _, z := range zIndexes {
		for _, e := range depthField[z] {
			e.Draw(s)
		}
	}
}

// Tick first removes all entity previously registered to be removed.
// Then calls original Tick logic.
func (w *World) Tick(e tl.Event) {
	for _, entity := range w.entitiesToRemove {
		w.BaseLevel.RemoveEntity(entity)
	}
	w.BaseLevel.Tick(e)
}

// IsLowColor is helper function for quick access to LowColor world option in Draw methods
func IsLowColor(s *tl.Screen) bool {
	if world, ok := s.Level().(*World); ok {
		return world.options.LowColor
	}
	panic("Invalid level. World ecxpected")
}

// ZIndexer if implemented by entity allows to be sorted by z-index.
// Lower z-index will be drawn first.
// Zero is used as default value for entities do not implementing ZIndexer.
type ZIndexer interface {
	ZIndex() int
}
