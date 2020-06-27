package core

// Player holds stats and attributes of player
type Player struct {
	// Name is the name of this player
	Name string
	// Stats are the statistics about player from previous rounds
	Stats Stats
	// Attributes are the player's attributes
	Attributes Attributes
}

// NewPlayer creates player with given name and with default attributes
func NewPlayer(name string) *Player {
	return &Player{
		Name: name,
		Attributes: Attributes{
			Attack:  1,
			Defense: 1,
		},
	}
}

// AddStats add to each of player stats values given by s.
// Use to report statistics after round finish
func (p *Player) AddStats(s Stats) {
	p.Stats.Kills += s.Kills
	p.Stats.Deaths += s.Deaths
	p.Stats.Suicides += s.Suicides
}

// Players is array of multiple players
type Players []*Player

// Attributes holds players's attributes.
//
// There are two base attributes Attack and Defense.
// Attack affects explosion size and maximum shooting power.
// Defense affects starting amount of armour.
//
// Additionally there is attribute Points.
// It holds number of points possible to redistribute between other attributes.
type Attributes struct {
	Attack, Defense, Points int
}

// Explosion is value used to calculate size of the bullet explosion.
// Bigger value = bigger explosion size.
func (s *Attributes) Explosion() int {
	return s.Attack
}

// Power is maximum shooting power
func (s *Attributes) Power() int {
	return 100 + (s.Attack-1)*5
}

// Armour is amout of armour / health on round start
func (s *Attributes) Armour() int {
	return 100 + (s.Defense-1)*5
}

// Stats are player statistics collected during multiple rounds
type Stats struct {
	// how many times player hits some enemy
	Kills int
	// how many times player was hit by some enemy
	Deaths int
	// how many times player killed himself
	Suicides int
}
