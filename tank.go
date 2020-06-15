package gorched

import (
	"math/rand"

	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/debug"
	"github.com/zladovan/gorched/draw"
	"github.com/zladovan/gorched/gmath"
	"github.com/zladovan/gorched/physics"
)

// Tank represents player's entity.
// It's tank which can change the angle of it's cannon.
// It can choose the shooting power and shoot bullet with given angle and power.
type Tank struct {
	// it extends from termloop.Entity
	*tl.Entity
	// player is reference to Player controlling this Tank
	player *Player
	// body is physical body of the tank used for falling simulation
	body *physics.Body
	// angle of cannon, 0 points to the right, 180 to the left
	angle int
	// power which will be used to shoot bullet, can be 0 - 100
	power float64
	// color of this tank
	color tl.Attr
	// state describes the current state of Tank
	state TankState
	// previousState holds state in previous frame, it's useful for actions on state transitions
	previousState TankState
	// callback called when shooted bullet finishes his path
	onShootingFinished func()
	// label is used to display info about angle, power or to show some message
	label *Label
	// asciiOnly if true will change sprite of the tank to the one containing no unicode characters
	asciiOnly bool
}

// TankState describes the state of Tank
type TankState uint8

const (
	// Idle is the state when tank is doing nothing but it's ready to go
	Idle TankState = iota
	// Loading is the state when tank is preparing to shoot and it's power is changing
	Loading
	// Shooting is the state when tank will shoot a bullet
	Shooting
	// Dead is the state after tank was hit and he is out of game
	Dead
)

// NewTank creates tank for given player.
func NewTank(player *Player, position gmath.Vector2i, angle int, color tl.Attr, asciiOnly bool) *Tank {
	return &Tank{
		Entity: tl.NewEntityFromCanvas(position.X-2, position.Y-3, *createCanvas(angle, color, asciiOnly)),
		player: player,
		body: &physics.Body{
			Position: *position.As2F(),
			Mass:     3,
		},
		angle:     angle,
		color:     color,
		label:     NewLabel(position.X+1, position.Y-4, color),
		asciiOnly: asciiOnly,
	}
}

// create canvas with tank model
func createCanvas(angle int, color tl.Attr, asciiOnly bool) *tl.Canvas {
	p := draw.BlankPrinter(6, 3).WithFg(color)
	if asciiOnly {
		printModelASCIIOnly(p, angle)
	} else {
		printModel(p, angle)
	}
	return p.Canvas
}

// Draw tank in one of the folowing positions depending on it's angle
//
// "  ▄▂▂"		 0 - 14
// "[██]"
// "◥@@◤"
//
// "  ▄▂▬"		15 - 44
// "[██]"
// "◥@@◤"
//
// "  ▄▬▀"		45 - 74
// "[██]"
// "◥@@◤"
//
// "  ▋ "		 75 - 104
// "[██]"
// "◥@@◤"
//
// "▀▬▄"		105 - 134
// "  [██]"
// "  ◥@@◤"
//
// "▬▂▄"		135 - 164
// "  [██]"
// "  ◥@@◤"
//
// "▂▂▄"		 165 - 180
// "  [██]"
// "  ◥@@◤"
func printModel(p *draw.Printer, angle int) {
	// Draw cannon
	switch {
	case angle < 15:
		p.Write(3, 0, "▄▂▂")
	case angle < 45:
		p.Write(3, 0, "▄▂▬")
	case angle < 75:
		p.Write(3, 0, "▄▬▀")
	case angle < 105:
		p.Write(3, 0, "▋")
	case angle < 135:
		p.Write(0, 0, "▀▬▄")
	case angle < 165:
		p.Write(0, 0, "▬▂▄")
	case angle < 181:
		p.Write(0, 0, "▂▂▄")
	}

	// Draw body
	p.Write(1, 1, "[██]")

	// Draw chasis
	p.Write(1, 2, "◥@@◤")
}

// Draw tank using only ASCII characters in one of the folowing positions depending on it's angle
//
// "  ▄▬■"		 0 - 14
// "[██]"
// "{@@}"
//
// "  ▄▬▀"		15 - 44
// "[██]"
// "{@@}"
//
// "  ▄▀ "		45 - 74
// "[██]"
// "{@@}"
//
// "  ▄ "		 75 - 104
// "[██]"
// "{@@}"
//
// " ▀▄"		105 - 134
// "  [██]"
// "  {@@}"
//
// "▀▬▄"		135 - 164
// "  [██]"
// "  {@@}"
//
// ■▬▄"		 165 - 180
// "  [██]"
// "  {@@}"
func printModelASCIIOnly(p *draw.Printer, angle int) {
	// Draw cannon
	switch {
	case angle < 15:
		p.Write(3, 0, "▄▬■")
	case angle < 45:
		p.Write(3, 0, "▄▬▀")
	case angle < 75:
		p.Write(3, 0, "▄▀")
	case angle < 105:
		p.Write(3, 0, "▄")
	case angle < 135:
		p.Write(0, 0, " ▀▄")
	case angle < 165:
		p.Write(0, 0, "▀▬▄")
	case angle < 181:
		p.Write(0, 0, "■▬▄")
	}

	// Draw body
	p.Write(1, 1, "[██]")

	// Draw chasis
	p.Write(1, 2, "{@@}")
}

// MoveUp increase cannon's angle
func (t *Tank) MoveUp() {
	t.updateAngle(1)
}

// MoveDown decrease cannon's angle
func (t *Tank) MoveDown() {
	t.updateAngle(-1)
}

// updates cannon's angle by given change
func (t *Tank) updateAngle(change int) {
	// TODO: angle should be updated by delta time to avoid lags
	t.angle += change
	if t.angle > 180 {
		t.angle = 180
	} else if t.angle < 0 {
		t.angle = 0
	}
	t.label.ShowNumber(t.angle)
	t.Entity.SetCanvas(createCanvas(t.angle, t.color, t.asciiOnly))
}

// Shoot will start loading when called first time and shoot bullet when started second time.
// Given onFinish callback is called  when shooted bullet finishes his path and hit to some obstacle or disapears out of world.
func (t *Tank) Shoot(onFinish func()) {
	switch t.state {
	case Idle:
		t.state = Loading
		t.power = 0
	case Loading:
		t.state = Shooting
		t.onShootingFinished = func() {
			if t.state != Dead {
				t.state = Idle
			}
			onFinish()
		}
	}
}

// phrases which are shown when tank's bullet hit some enemy
var phrasesAfterHit = []string{
	// TODO: more phrases
	"Yeeha !",
	"Take that !",
	"¡Hasta la vista!",
	"Bang !",
	"Rest in pieces !",
}

// Hit should be called when this tank hit some enemy
func (t *Tank) Hit() {
	t.label.Show(phrasesAfterHit[rand.Intn(len(phrasesAfterHit))])
	t.player.hits++
}

// TakeDamage should be called when this tank was hit by some enemy
func (t *Tank) TakeDamage() {
	t.state = Dead
	t.player.takes++
}

// IsAlive returns wether this tank is still in game
func (t *Tank) IsAlive() bool {
	return t.state != Dead
}

// Tick is not used now
func (t *Tank) Tick(e tl.Event) {}

// Draw tank
func (t *Tank) Draw(s *tl.Screen) {
	// TODO: simplify by creating label with relative position
	// update entity and label positions based on body position
	y := int(t.body.Position.Y) - 3
	t.Entity.SetPosition(int(t.body.Position.X)-2, y)
	t.label.position.Y = y - 1
	lx, _ := t.label.Text.Position()
	t.label.Text.SetPosition(lx, t.label.position.Y)

	switch t.state {
	case Shooting:
		if t.previousState != Shooting {
			// create new bullet
			debug.Logf("Tank shooting angle=%d power=%f", t.angle, t.power)
			// TODO: choose strength of bullet based on player stats
			s.Level().AddEntity(NewBullet(t, t.getBulletInitPos(), float64(int(t.power)), t.angle, 4, t.onShootingFinished))
		}
	case Loading:
		// increase shooting power
		// idea is that increase should be faster for each next 5 points
		t.power += (10 + t.power/5) * s.TimeDelta()
		if t.power >= 100 {
			t.power = 1
		}
		t.label.ShowNumber(int(t.power))
	case Dead:
		if t.previousState != Dead {
			explosion := NewExplosion(*t.body.Position.Translate(0, -2).As2I(), 6)
			tomb := NewTomb(*t.body.Position.As2I(), t.color)
			s.Level().AddEntity(explosion)
			s.Level().AddEntity(&SpawnAfter{Entity: tomb, After: explosion})
			s.Level().RemoveEntity(t)
		}
	}
	t.previousState = t.state

	// draw underlying entity
	t.Entity.Draw(s)
	// draw label above tank
	t.label.Draw(s)
}

// calculates initial position of the bullet
func (t *Tank) getBulletInitPos() gmath.Vector2i {
	x, y := t.Entity.Position()
	x += 2 // move to the center (almost) of the tank
	if t.angle >= 75 && t.angle < 105 {
		y--
	}
	if t.angle < 75 {
		x += 3
	}
	if t.angle >= 105 {
		x -= 2
	}
	return gmath.Vector2i{X: x, Y: y}
}

// Position returns collider position
func (t *Tank) Position() (int, int) {
	// position for collider is moved to do not include cannon edge
	x, y := t.Entity.Position()
	return x + 1, y
}

// Size returns collider size
func (t *Tank) Size() (int, int) {
	// collider is little bit smaller than 6x3 canvas to do not include cannon edge
	return 4, 3
}

// ZIndex return z-index of tank.
// It should be bigger than z-index of terrain and trees.
func (t *Tank) ZIndex() int {
	return 2000
}

// Body returns physical body of the tank used for falling simulation
func (t *Tank) Body() *physics.Body {
	return t.body
}

// BottomLine returns line x coordinates for collision with the ground when falling
func (t *Tank) BottomLine() (int, int) {
	return 0, 1
}

// Angle returns angle of tank's cannon
func (t *Tank) Angle() int {
	return t.angle
}

// Power returns power which will be used to shoot bullet, can be 0 - 100
func (t *Tank) Power() int {
	return int(t.power)
}

// IsIdle returns true if tank is in Idle state
func (t *Tank) IsIdle() bool {
	return t.state == Idle
}

// IsLoading returns true if tank is loading now
func (t *Tank) IsLoading() bool {
	return t.state == Loading
}

// Tomb is entity representing tomb stone shown on position where tank was killed
type Tomb struct {
	body   *physics.Body
	canvas *tl.Canvas
}

// NewTomb creates new Tomb entity on given position with given color.
// It should have same color as tank instead of which it was added to world.
func NewTomb(position gmath.Vector2i, color tl.Attr) *Tomb {
	return &Tomb{
		body: &physics.Body{
			Position: *position.As2F(),
			Mass:     3,
		},
		canvas: createTombCanvas(color),
	}
}

// createTombCanvas creates canvas with tomb sprite
func createTombCanvas(color tl.Attr) *tl.Canvas {
	p := draw.BlankPrinter(6, 3).WithFg(color)
	p.WriteLines(1, 1, []string{
		" ▄█▄",
		"  █",
	})
	return p.Canvas
}

// Draw draws tomb stone
func (t *Tomb) Draw(s *tl.Screen) {
	offsetx := -1
	offsety := -3
	for i := 0; i < len(*t.canvas); i++ {
		for j := 0; j < len((*t.canvas)[0]); j++ {
			s.RenderCell(int(t.body.Position.X)+i+offsetx, int(t.body.Position.Y)+j+offsety, &(*t.canvas)[i][j])
		}
	}
}

// Tick does nothing now
func (t *Tomb) Tick(e tl.Event) {}

// ZIndex return z-index of the tomb.
// It should be lower than z-index of tank.
func (t *Tomb) ZIndex() int {
	return 1999
}

// Body returns physical body of the tomb used for falling simulation
func (t *Tomb) Body() *physics.Body {
	return t.body
}

// BottomLine returns line x coordinates for collision with the ground when falling
func (t *Tomb) BottomLine() (int, int) {
	return 1, 1
}
