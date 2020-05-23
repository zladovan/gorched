package gorched

import (
	//"time"

	"fmt"
	"strings"

	tl "github.com/JoelOtter/termloop"
)

// TODO: extract message boxes

// Controls holds data and logic for controlling game world.
type Controls struct {
	// reference to game
	game *Game
	// all tanks in game
	tanks []*Tank
	// index of tank which is active and controlled
	activeTankIndex int
	// true if info message box should be displayed
	showInfo bool
	// true if score board should be displayed
	showScore bool
}

// Tick handles all key events
func (c *Controls) Tick(e tl.Event) {
	// TODO: simplify

	// when info is show it's just possible to hide it
	if c.showInfo {
		if e.Type == tl.EventKey {
			c.ToggleInfo()
		}
		return
	}

	// when score is show it's just possible to hide it
	if c.showScore {
		if e.Type == tl.EventKey {
			if c.NumberOfTanksAlive() <= 1 {
				c.game.NextRound()
			} else {
				c.ToggleScore()
			}
		}
		return
	}

	// otherwise handle in-game controls
	switch e.Key {
	case tl.KeyArrowLeft:
		c.ActiveTank().MoveUp()
	case tl.KeyArrowRight:
		c.ActiveTank().MoveDown()
	case tl.KeySpace:
		c.ActiveTank().Shoot(func() {
			if c.NumberOfTanksAlive() <= 1 {
				c.showScore = true
				c.showInfo = false
			} else {
				c.ActivateNextTank()
			}
		})
	case tl.KeyCtrlR:
		c.game.RestartRound()
	case tl.KeyCtrlN:
		c.game.NextRound()
	}
	switch e.Ch {
	case 'h':
		c.ToggleInfo()
	case 's':
		c.ToggleScore()
	}
}

// ActivateNextTank moves turn to nearest tank which is alive.
func (c *Controls) ActivateNextTank() {
	c.activeTankIndex = (c.activeTankIndex + 1) % len(c.tanks)
	if !c.ActiveTank().IsAlive() && c.NumberOfTanksAlive() > 0 {
		c.ActivateNextTank()
	}
}

// ActiveTank returns tank which is currenlty active / on turn.
func (c *Controls) ActiveTank() *Tank {
	return c.tanks[c.activeTankIndex]
}

// ToggleInfo shows / hides info message box.
func (c *Controls) ToggleInfo() {
	c.showInfo = !c.showInfo
}

// ToggleScore shows / hides score board.
func (c *Controls) ToggleScore() {
	c.showScore = !c.showScore
}

// NumberOfTanksAlive returns how many tanks is still alive (in game).
func (c *Controls) NumberOfTanksAlive() int {
	alive := 0
	for _, t := range c.tanks {
		if t.IsAlive() {
			alive++
		}
	}
	return alive
}

// text of info message box
const infoText = `
╔═════════════════════════════════════════════════════╗
║                ╔═╗╔═╗┬─┐┌─┐┬ ┬┌─┐┌┬┐                ║
║                ║ ╦║ ║├┬┘│  ├─┤├┤  ││                ║
║                ╚═╝╚═╝┴└─└─┘┴ ┴└─┘─┴┘                ║
║                                                     ║
║ Left / Right   change cannon angle                  ║
║ SPACE          start loading (1st) and shoot (2nd)  ║
║ Ctrl+C         exit game                            ║
║ Ctrl+R         restart current round                ║
║ Ctrl+N         start next round                     ║
║   S            show score                           ║
║   H            show help                            ║
║                                                     ║
║                  © 2020, Zladovan                   ║
╚═════════════════════════════════════════════════════╝
`

// header of scoreboard
const scoreHeader = `
╔═════════════════════════════════════════════════════╗
║                   ╔═╗┌─┐┌─┐┬─┐┌─┐                   ║
║                   ╚═╗│  │ │├┬┘├┤                    ║
║                   ╚═╝└─┘└─┘┴└─└─┘                   ║
║                                                     ║
║                  Kills        Deaths                ║`

// format string used for showing score for each player, expects player's name, number of kills, number of deaths
const scoreRow = `
║ %-10s        %4d             %-4d             ║`

// footer of scoreboard
var scoreFooter = strings.TrimSpace(`
║                                                     ║
╚═════════════════════════════════════════════════════╝
`)

// Draw controls
func (c *Controls) Draw(s *tl.Screen) {
	if c.showInfo {
		drawMessage(infoText, s, c.game.options.LowColor)
	}
	if c.showScore {
		b := &strings.Builder{}
		fmt.Fprint(b, scoreHeader)
		for i, p := range c.game.players {
			fmt.Fprintf(b, scoreRow, fmt.Sprintf("Player %d", i+1), p.hits, p.takes)
		}
		fmt.Fprintln(b)
		fmt.Fprint(b, scoreFooter)
		drawMessage(b.String(), s, c.game.options.LowColor)
	}
}

// draw message box with given text
func drawMessage(message string, s *tl.Screen, lowColor bool) {
	bg := tl.RgbTo256Color(50, 50, 50)
	fg := tl.RgbTo256Color(200, 200, 200)
	if lowColor {
		bg = tl.ColorBlack
		fg = tl.ColorWhite
	}
	info := NewMessage(message, fg, bg)
	MoveToScreenCenter(info, s)
	info.Draw(s)
}
