package hud

import (
	"fmt"
	"strings"

	"github.com/zladovan/gorched/core"
	"github.com/zladovan/gorched/hud/ui"
)

// header of scoreboard
var scoreHeader = Trim(`
                  ╔═╗┌─┐┌─┐┬─┐┌─┐                  
                  ╚═╗│  │ │├┬┘├┤                   
                  ╚═╝└─┘└─┘┴└─└─┘                  
                                                
                 Kills        Deaths      Suicides
`)

// format string used for showing score for each player, expects player's name, number of kills, deaths and suicides
var scoreRow = Trim(`
%-10s        %4d             %-4d          %-4d
`)

// NewScoreBox creates MessageBox with current score
func NewScoreBox(players core.Players) *ui.MessageBox {
	b := &strings.Builder{}
	fmt.Fprint(b, scoreHeader)
	for _, p := range players {
		fmt.Fprintln(b)
		fmt.Fprintf(b, scoreRow, p.Name, p.Stats.Kills, p.Stats.Deaths, p.Stats.Suicides)
	}
	fmt.Fprintln(b)
	return ui.NewMessageBox(b.String())
}
