package gorched

import (
	"math/rand"

	tl "github.com/JoelOtter/termloop"
)

// WorldOptions provide configuration needed for generating game world (one round).
type WorldOptions struct {
	// Width of game world in number of console pixels (cells)
	Width int
	// Height of game world in number of console pixels (cells)
	Height int
	// Seed is number used as random seed and if it is reused it allows to create same game looking game world with same positions for players
	Seed int64
}

// NewWorld creates new game world with all entities 
func NewWorld(game *Game, o WorldOptions) *tl.BaseLevel {
	// random positions in the world are seeded too
	rnd := rand.New(rand.NewSource(o.Seed))
	
	// create terrain
	terrain := GenerateTerrain(&TerrainGenerator{
		Seed:      o.Seed,
		Width:     o.Width,
		Height:    o.Height,
		Roughness: 7.5,
	})

	// create clouds
	clouds := GenerateClouds(&CloudsGenerator{seed: o.Seed, width: o.Width, height: o.Height})

	// create players
	// TODO: update for different player counts
	tanks := []*Tank{
		NewTank(
			game.players[0],
			terrain.GetPositionOn(10 + rnd.Intn(10)),
			0,
			tl.RgbTo256Color(200, 0, 0),
			game.options.ASCIIOnly,
		),
		NewTank(
			game.players[1],
			terrain.GetPositionOn(o.Width - 10 - rnd.Intn(10)),
			180,
			tl.ColorBlack,
			game.options.ASCIIOnly,
		),
	}
	
	// create controls
	controls := &Controls{
		game: game, 
		tanks: tanks, 
		showInfo: game.CurrentRound() == 1,
		activeTankIndex: game.startingPlayerIndex,
	}

	// create level with all entities
	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlue,
	})
	level.AddEntity(clouds)
	level.AddEntity(terrain)
	for _, c := range terrain.GetColliders() {
		level.AddEntity(c)
	}
	for _, t := range tanks {
		level.AddEntity(t)
	}
	level.AddEntity(controls)

	return level
}
