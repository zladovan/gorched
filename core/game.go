package core

// Game is only holder of the players for now
type Game interface {
	// Players returns all players in the game
	Players() Players
}
