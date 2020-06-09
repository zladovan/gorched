package gorched

import (
	"fmt"
	"strings"

	tl "github.com/JoelOtter/termloop"

	"github.com/zladovan/gorched/draw"
)

// MessageBox shows box with some message in the screen center.
type MessageBox struct {
	*tl.Entity
}

// NewMessageBox creates new message box from given string with multiple lines.
// Create MessageBox will have dimensions based on the lines resolved from given string.
func NewMessageBox(msg string, lowColor bool) *MessageBox {
	// resolve lines and dimensions based on the lines
	lines := strings.Split(strings.TrimSpace(msg), "\n")
	h := len(lines)
	w := len([]rune(lines[0]))

	// resolve colors
	bg := tl.RgbTo256Color(50, 50, 50)
	fg := tl.RgbTo256Color(200, 200, 200)
	if lowColor {
		bg = tl.ColorBlack
		fg = tl.ColorWhite
	}

	p := draw.BlankPrinter(w, h).WithColors(fg, bg)
	p.WriteLines(0, 0, lines)
	return &MessageBox{
		Entity: tl.NewEntityFromCanvas(0, 0, *p.Canvas),
	}
}

// Draw draws message box on the center of the screen
func (m *MessageBox) Draw(s *tl.Screen) {
	w, h := m.Entity.Size()
	sw, sh := s.Size()
	m.Entity.SetPosition(sw/2-w/2, sh/2-h/2)
	m.Entity.Draw(s)
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

// NewInfoBox creates MessageBox with main game info
func NewInfoBox(browserMode bool, lowColor bool) *MessageBox {
	text := infoText
	// in browser mode some controls are different to do not collide with browser shortcuts
	if browserMode {
		text = strings.ReplaceAll(text, "Ctrl+R", "  R   ")
		text = strings.ReplaceAll(text, "Ctrl+N", "  N   ")
	}
	return NewMessageBox(text, lowColor)
}

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

// NewScoreBox creates MessageBox with current score
func NewScoreBox(game *Game) *MessageBox {
	b := &strings.Builder{}
	fmt.Fprint(b, scoreHeader)
	for i, p := range game.players {
		fmt.Fprintf(b, scoreRow, fmt.Sprintf("Player %d", i+1), p.hits, p.takes)
	}
	fmt.Fprintln(b)
	fmt.Fprint(b, scoreFooter)
	return NewMessageBox(b.String(), game.options.LowColor)
}
