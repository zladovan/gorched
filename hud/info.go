package hud

import (
	"strings"

	"github.com/zladovan/gorched/hud/ui"
)

// text of info message box
var infoText = Trim(`
               ╔═╗╔═╗┬─┐┌─┐┬ ┬┌─┐┌┬┐              
               ║ ╦║ ║├┬┘│  ├─┤├┤  ││              
               ╚═╝╚═╝┴└─└─┘┴ ┴└─┘─┴┘              
                                            
Left / Right   change cannon angle                
SPACE          start loading (1st) and shoot (2nd)
Ctrl+C         exit game                          
Ctrl+R         restart current round              
Ctrl+N         start next round                   
  S            show score       
  A            show player's attributes
  H            show help                          
                                            
                 © 2020, Zladovan                 
`)

// NewInfoBox creates MessageBox with main game info
func NewInfoBox(browserMode bool) *ui.MessageBox {
	text := infoText
	// in browser mode some controls are different to do not collide with browser shortcuts
	if browserMode {
		text = strings.ReplaceAll(text, "Ctrl+R", "  R   ")
		text = strings.ReplaceAll(text, "Ctrl+N", "  N   ")
	}
	return ui.NewMessageBox(text)
}
