package core

// Player holds stats and attributes of player
type Player struct {
	// Name is the name of this player
	Name string
	// Stats are the statistics about player from previous rounds
	Stats Stats
}

// NewPlayer creates player with given name and with default attributes
func NewPlayer(name string) *Player {
	return &Player{
		Name: name,
	}
}
