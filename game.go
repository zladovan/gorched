package gorched

import tl "github.com/JoelOtter/termloop"

// Game holds information which is kept during whole session.
type Game struct {
	// players holds all players in the game
	players []*Player
	// engine references to termloop's game
	engine *tl.Game
	// options holds game options used to create new game
	options GameOptions
	// rounds holds number of finished rounds
	rounds int
	// startingPlayerIndex holds index of player which was first on turn in current round
	startingPlayerIndex int
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
	// Debug turns on debug mode if set to true
	Debug bool
}

// NewGame creates new game object.
// Game is not started yet. You need to call Start().
func NewGame(o GameOptions) *Game {
	game := &Game{}
	game.options = o
	game.engine = tl.NewGame()
	game.engine.Screen().SetFps(float64(o.Fps))
	if o.Debug {
		Debug.Attach(game.engine)
	}
	game.players = make([]*Player, o.PlayerCount)
	for pi := range game.players {
		game.players[pi] = &Player{}
	}
	return game
}

// Start starts the game which means that game engine is started and first round is set up.
func (g *Game) Start() {
	g.RestartRound()
	g.engine.Start()
}

// NextRound finish current round as it is and switches to new round.
func (g *Game) NextRound() {
	g.rounds++
	g.startingPlayerIndex = (g.startingPlayerIndex + 1) % len(g.players)
	g.RestartRound()
}

// RestartRound regenerates again current round to the same state as when it was started.
func (g *Game) RestartRound() {
	world := NewWorld(g, WorldOptions{
		Width:  g.options.Width,
		Height: g.options.Height,
		Seed:   g.LastSeed(),
	})
	g.engine.Screen().SetLevel(world)
}

// InitialSeed returns seed used for the first level.
func (g *Game) InitialSeed() int64 {
	return g.options.Seed
}

// LastSeed returns seed used for the last (current active) level.
func (g *Game) LastSeed() int64 {
	return g.options.Seed + int64(g.rounds)
}

// CurrentRound returns number of actual round starting with 1 for the first row.
func (g *Game) CurrentRound() int {
	return g.rounds + 1
}
