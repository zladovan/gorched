package gorched

import tl "github.com/JoelOtter/termloop"

// Round represents one round in the game.
// It is responsible for managing state of the round.
// It also provides functionality for transitions to next rounds.
// Use NewRound to create new instance.
// Add it to the termloop.Screen entities after creating.
// Call Restart to restart this round.
// Call Next to go to the next round.
type Round struct {
	// game refers to the main game object
	game *Game
	// world refers to actual game world
	world *World
	// index is the index of current round, starting from zero for first round
	index int
	// state holds the state of this round
	state RoundState
	// tanks contains all tanks in game
	tanks []*Tank
	// startingPlayerIndex holds index of the player which is first on turn in this round
	startingPlayerIndex int
	// onTurnPlayerIndex is index of the player currently on turn
	onTurnPlayerIndex int
}

// RoundState represents state of the round
type RoundState uint8

const (
	// Started round is state right after round was created / started / restarted
	Started RoundState = iota
	// PlayerOnTurn is state when some player is on turn but it has not done his move yet
	PlayerOnTurn
	// WaitForTurnFinish is state when some player did his move and we are waiting for all consequences of the move
	WaitForTurnFinish
	// Finished is state when round was finished which means there is only one or zero tanks alive
	Finished
)

// NewRound creates new round.
// Created round will be in Started state.
// It will add World as level after added to the screen.
func NewRound(game *Game) *Round {
	round := &Round{game: game}
	round.Restart()
	return round
}

// Draw is processing Round states
func (r *Round) Draw(s *tl.Screen) {
	switch r.state {
	case Started:
		s.SetLevel(r.world)
		r.onTurnPlayerIndex = r.startingPlayerIndex
		r.state = PlayerOnTurn
	case PlayerOnTurn:
		if r.ActiveTank().IsShooting() {
			r.state = WaitForTurnFinish
		}
	case WaitForTurnFinish:
		if r.IsTurnFinished() {
			if r.NumberOfTanksAlive() <= 1 {
				r.finishRound()
			} else {
				r.state = PlayerOnTurn
				r.ActivateNextTank()
			}
		}
	}
}

// finishRound is called when round is finished
func (r *Round) finishRound() {
	r.state = Finished

	// states gained during this round are added to players on round finish
	// all players that didn't made suicide gain one point
	// winner gain one more point
	for pi, player := range r.game.players {
		tank := r.tanks[pi]
		player.AddStats(tank.stats)
		if tank.stats.Suicides == 0 {
			player.Attributes.Points++
		}
		if tank.IsAlive() {
			player.Attributes.Points++
		}
	}

	// score board following by attributes form are shown at the end of round
	score := r.game.Hud().ShowScore()
	score.OnClose(func() {
		attrs := r.game.Hud().ShowAttributes(false)
		attrs.OnClose(func() {
			r.Next()
		})
	})
}

// Tick does nothing now
func (r *Round) Tick(e tl.Event) {}

// Restart will put state of this round to the same state as when it was started.
func (r *Round) Restart() {
	// create world
	r.world = NewWorld(r.game, WorldOptions{
		Width:     r.game.options.Width,
		Height:    r.game.options.Height,
		Seed:      r.game.options.Seed + int64(r.index),
		ASCIIOnly: r.game.options.ASCIIOnly,
		LowColor:  r.game.options.LowColor,
	})

	// collect tanks for players
	r.tanks = make([]*Tank, len(r.game.players))
	for i, player := range r.game.players {
		for _, e := range r.world.Entities {
			if tank, ok := e.(*Tank); ok && tank.player == player {
				r.tanks[i] = tank
			}
		}
	}

	// round is started again
	r.state = Started
}

// Next will go to the next round.
func (r *Round) Next() {
	r.index++
	r.startingPlayerIndex = (r.startingPlayerIndex + 1) % len(r.game.players)
	r.Restart()
}

// Number returns number of this round starting with 1 for the first round
func (r *Round) Number() int {
	return r.index + 1
}

// ActiveTank returns tank which is currently active / on turn.
func (r *Round) ActiveTank() *Tank {
	return r.tanks[r.onTurnPlayerIndex]
}

// IsTurnFinished returns true if there are no bullets and explosions in world
func (r *Round) IsTurnFinished() bool {
	for _, e := range r.world.Entities {
		if _, ok := e.(*Bullet); ok {
			return false
		}
		if _, ok := e.(*Explosion); ok {
			return false
		}
	}
	return true
}

// NumberOfTanksAlive returns how many tanks is still alive (in game).
func (r *Round) NumberOfTanksAlive() int {
	alive := 0
	for _, t := range r.tanks {
		if t.IsAlive() {
			alive++
		}
	}
	return alive
}

// ActivateNextTank moves turn to nearest tank which is alive.
func (r *Round) ActivateNextTank() {
	r.onTurnPlayerIndex = (r.onTurnPlayerIndex + 1) % len(r.tanks)
	if !r.ActiveTank().IsAlive() && r.NumberOfTanksAlive() > 0 {
		r.ActivateNextTank()
	}
}

// IsFinished returns true when round was already finished
func (r *Round) IsFinished() bool {
	return r.state == Finished
}

// IsPlayerOnTurn returns turn when some player is on turn now and he didn't made his move yet
func (r *Round) IsPlayerOnTurn() bool {
	return r.state == PlayerOnTurn
}
