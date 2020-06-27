package hud

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/core"
	"github.com/zladovan/gorched/gmath"
	"github.com/zladovan/gorched/hud/ui"
)

var attributesPageLayout = Trim(`
            ╔═╗┌┬┐┌┬┐┬─┐┬┌┐ ┬ ┬┌┬┐┌─┐┌─┐
            ╠═╣ │  │ ├┬┘│├┴┐│ │ │ ├┤ └─┐
            ╩ ╩ ┴  ┴ ┴└─┴└─┘└─┘ ┴ └─┘└─┘

Player 1

Attack    [  1] + -     Explosion strength    [  1]
                        Shooting power        [100]
Defense   [  1] + -     Armour                [100]

Points    [  2]

| Press [Tab] to change focus.
| Press [Enter] to do action. 

                               Previous Next Finish
`)

// AttributesForm shows player's attributes and allows to modify them.
//
// It contains multiple pages, one for each player.
// There are also buttons for navigation between pages.
//
// Each page shows values for all attributes.
// If form is not in read only mode there are also + / - buttons
// which can be used to modify base attributes.
type AttributesForm struct {
	*ui.BaseForm
	// players holds players which attributes are shown
	players core.Players
	// activePage is index of currently visible page
	activePage int
	// pages hold one container for each player
	pages []ui.Container
	// readOnly if true disables attribute modification buttons
	readOnly bool
}

// NewAttributesForm creates new form for showing and modifying players attributes.
// If readOnly is set to true it will not allow to modify attributes.
func NewAttributesForm(players core.Players, readOnly bool) *AttributesForm {
	f := &AttributesForm{
		BaseForm: ui.NewForm(),
		players:  players,
		readOnly: readOnly,
	}
	f.initPages()
	return f
}

// initPages creates container for each player
func (f *AttributesForm) initPages() {
	f.pages = make([]ui.Container, len(f.players))
	for i, player := range f.players {
		f.pages[i] = f.createPage(i, player)
	}
	if len(f.pages) > 0 {
		f.SetContainer(f.pages[0])
	}
}

// createPage creates one container with all components for one player
func (f *AttributesForm) createPage(pageIndex int, player *core.Player) *ui.BaseContainer {
	// label with player name
	name := ui.NewText(player.Name)
	name.Colors.Fg = name.Colors.Fg | tl.AttrBold

	// attribute values
	attack := ui.NewValue(player.Attributes.Attack)
	explosion := ui.NewValue(player.Attributes.Explosion())
	power := ui.NewValue(player.Attributes.Power())
	defense := ui.NewValue(player.Attributes.Defense)
	armour := ui.NewValue(player.Attributes.Armour())
	points := ui.NewValue(0)
	points.Add(player.Attributes.Points)
	attrs := []ui.Component{attack, explosion, power, defense, armour, points}

	// plus / minus buttons
	attackPlus := ui.NewButton("+", func() {
		if points.Get() > 0 && attack.Get() < 100 {
			attack.Add(1)
			points.Add(-1)
			player.Attributes.Attack++
			player.Attributes.Points--
			explosion.Add(player.Attributes.Explosion() - explosion.Get())
			power.Add(player.Attributes.Power() - power.Get())
		}
	})
	attackMinus := ui.NewButton("-", func() {
		if attack.Addition() > 0 {
			attack.Add(-1)
			points.Add(1)
			player.Attributes.Attack--
			player.Attributes.Points++
			explosion.Add(player.Attributes.Explosion() - explosion.Get())
			power.Add(player.Attributes.Power() - power.Get())
		}
	})
	defensePlus := ui.NewButton("+", func() {
		if points.Get() > 0 && defense.Get() < 100 {
			defense.Add(1)
			points.Add(-1)
			player.Attributes.Defense++
			player.Attributes.Points--
			armour.Add(player.Attributes.Armour() - armour.Get())
		}
	})
	defenseMinus := ui.NewButton("-", func() {
		if defense.Addition() > 0 {
			defense.Add(-1)
			points.Add(1)
			player.Attributes.Defense--
			player.Attributes.Points++
			armour.Add(player.Attributes.Armour() - armour.Get())
		}
	})
	buttons := []ui.Component{attackPlus, attackMinus, defensePlus, defenseMinus}

	// container for all components
	p := ui.NewFormatPane(attributesPageLayout, []*ui.ComponentBuilder{
		{
			Pattern: `Player \d`,
			Build: func(i int, s string) ui.Component {
				return name
			},
		},
		{
			Pattern: `\[\s*\d+\]`,
			Build: func(i int, s string) ui.Component {
				return attrs[i]
			},
		},
		{
			Pattern: `[\\+\\-]`,
			Skip:    f.readOnly,
			Build: func(i int, s string) ui.Component {
				return buttons[i]
			},
		},
		{
			Pattern: "Next",
			Skip:    pageIndex == len(f.players)-1,
			Build: func(i int, str string) ui.Component {
				b := ui.NewButton("Next", func() { f.Next() })
				b.ActionKey = 'N'
				return b
			},
		},
		{
			Pattern: "Finish",
			Build: func(i int, str string) ui.Component {
				b := ui.NewButton("Finish", func() { f.Close() })
				b.ActionKey = 'F'
				return b
			},
		},
		{
			Pattern: "Previous",
			Skip:    pageIndex == 0,
			Build: func(i int, str string) ui.Component {
				b := ui.NewButton("Previous", func() { f.Previous() })
				b.ActionKey = 'P'
				return b
			},
		},
	})
	p.Style().CopyFrom(f.Style())
	return p
}

// Next changes active page to the next page of the form
func (f *AttributesForm) Next() {
	f.changePage(1)
}

// Previous changes active page to the previous page of the form
func (f *AttributesForm) Previous() {
	f.changePage(-1)
}

// changePage adds given change to activePage index and switch form's container according it
func (f *AttributesForm) changePage(change int) {
	next := gmath.Min(len(f.players)-1, f.activePage+change)
	if next == f.activePage {
		return
	}
	f.activePage = next
	f.SetContainer(f.pages[f.activePage])
}
