package gorched

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/debug"
)

// Game holds information which is kept during whole session.
type Game struct {
	// players holds all players in the game
	players []*Player
	// engine references to termloop's game
	engine *tl.Game
	// options holds game options used to create new game
	options GameOptions
	// hud holds games HUD
	hud *HUD
	// controls contains processing of input
	controls *Controls
	// round is responsible for creating and managing state of game rounds
	round *Round
}

// Player holds stats for player which are aggregated during whole session.
type Player struct {
	// how many times player hits some enemy
	hits int
	// how many times player was hit by some enemy
	takes int
}

// GameOptions provide configuration needed for creating new game
type GameOptions struct {
	// Width of game world in number of console pixels (cells)
	Width int
	// Height of game world in number of console pixels (cells)
	Height int
	// PlayerCount is number of players which will be added to game
	PlayerCount int
	// Seed is number used as random seed and if it is reused it allows to play same game with same looking rounds
	Seed int64
	// Fps sets screen framerate
	Fps int
	// AsciiOnly identifies that only ASCII characters can be used for all graphics
	ASCIIOnly bool
	// LowColor identifies that only 8 colors can be used for all graphics
	LowColor bool
	// BrowserMode identifies that game was run in browser and some controls need to be modified to do not collide with usual browser shortcuts
	BrowserMode bool
	// Debug turns on debug mode if set to true
	Debug bool
}

// NewGame creates new game object.
// Game is not started yet. You need to call Start().
func NewGame(o GameOptions) *Game {
	game := &Game{}
	game.options = o

	// init engine
	game.engine = tl.NewGame()
	game.engine.Screen().SetFps(float64(o.Fps))

	// init debug
	if o.Debug {
		debug.Attach(game.engine)
	}

	// init players
	game.players = make([]*Player, o.PlayerCount)
	for pi := range game.players {
		game.players[pi] = &Player{}
	}

	// init HUD with info visible at startup
	game.hud = NewHUD(game)
	game.hud.ShowInfo()
	game.engine.Screen().AddEntity(game.hud)

	// init round
	game.round = NewRound(game)
	game.engine.Screen().AddEntity(game.round)

	// init controls
	game.controls = &Controls{game: game}
	game.engine.Screen().AddEntity(game.controls)

	return game
}

// Start starts the game which means that game engine is started and first round is set up.
func (g *Game) Start() {
	g.engine.Start()
}

// InitialSeed returns seed used for the first level.
func (g *Game) InitialSeed() int64 {
	return g.options.Seed
}

// LastSeed returns seed used for the last (current active) level.
func (g *Game) LastSeed() int64 {
	return g.options.Seed + int64(g.round.Number()-1)
}

// Hud returns games HUD
func (g *Game) Hud() *HUD {
	return g.hud
}

// Engine returns reference to underlying game engine
func (g *Game) Engine() *tl.Game {
	return g.engine
}
