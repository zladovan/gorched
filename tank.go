package gorched

import (
	"math/rand"

	tl "github.com/JoelOtter/termloop"
)

// Tank represents player's entity.
// It's tank which can change the angle of it's cannon.
// It can choose the shooting power and shoot bullet with given angle and power.
type Tank struct {
	// it extends from termloop.Entity
	*tl.Entity
	// player is reference to Player controlling this Tank
	player *Player
	// angle of cannon, 0 points to the right, 180 to the left
	angle int
	// power which will be used to shoot bullet, can be 0 - 100
	power float64
	// color of this tank
	color tl.Attr
	// if true tank was already hit and it's out of game
	isDead bool
	// if true shooting power is increasing and tank is preparing to shoot
	isLoading bool
	// if true tank is shootin now
	isShooting bool
	// if true tank cannot shot now because it's waiting for bullet finishes his path
	cannotShoot bool
	// callback called when shooted bullet finishes his path
	onShootingFinished func()
	// label is used to display info about angle, power or to show some message
	label     *Label
	asciiOnly bool
}

// NewTank creates tank for given player.
func NewTank(player *Player, position Position, angle int, color tl.Attr, asciiOnly bool) *Tank {
	return &Tank{
		Entity:    tl.NewEntityFromCanvas(position.x-2, position.y-3, *createCanvas(angle, color, asciiOnly)),
		player:    player,
		angle:     angle,
		color:     color,
		label:     NewLabel(position.x+1, position.y-4, color),
		asciiOnly: asciiOnly,
	}
}

// create canvas with tank model
func createCanvas(angle int, color tl.Attr, asciiOnly bool) *tl.Canvas {
	canvas := tl.NewCanvas(6, 3)
	p := &Printer{canvas: &canvas, fg: color}
	if asciiOnly {
		printModelAsciiOnly(p, angle)
	} else {
		printModel(p, angle)
	}
	return &canvas
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
func printModel(p *Printer, angle int) {
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
func printModelAsciiOnly(p *Printer, angle int) {
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

// draw dead tank
func createDeadCanvas(color tl.Attr) *tl.Canvas {
	canvas := tl.NewCanvas(6, 3)
	p := &Printer{canvas: &canvas, fg: color}
	p.WriteLines(1, 1, []string{
		" ▄█▄",
		"  █",
	})
	return &canvas
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
	// TODO: clenup
	if t.cannotShoot {
		return
	}
	if t.isLoading {
		// 2nd call = shoot bullet
		t.isLoading = false
		t.isShooting = true
		t.cannotShoot = true
		t.onShootingFinished = func() {
			t.cannotShoot = false
			onFinish()
		}
	} else if !t.isShooting {
		// 1st call start loading
		t.isLoading = true
		t.power = 0
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
	t.isDead = true
	t.player.takes++
	t.Entity.SetCanvas(createDeadCanvas(t.color))
}

// IsAlive returns wether this tank is still in game
func (t *Tank) IsAlive() bool {
	return !t.isDead
}

// Tick is not used now
func (t *Tank) Tick(e tl.Event) {}

// Draw tank
func (t *Tank) Draw(s *tl.Screen) {
	// draw underlying entity
	t.Entity.Draw(s)

	// create new bullet if tank is shooting
	if t.isShooting {
		s.Level().AddEntity(NewBullet(t, t.getBulletInitPos(), t.power, t.angle, t.onShootingFinished))
		t.isShooting = false
	}

	// increase shooting power if tank is loading
	if t.isLoading {
		// idea is that increase should be faster for each next 5 points
		t.power += (10 + t.power/5) * s.TimeDelta()
		if t.power >= 100 {
			t.power = 1
		}
		t.label.ShowNumber(int(t.power))
	}

	// draw label above tank
	t.label.Draw(s)
}

// calculates initial position of the bullet
func (t *Tank) getBulletInitPos() Position {
	x, y := t.Entity.Position()
	x += 2 // move to the center (almost) of the tank
	if t.angle >= 45 && t.angle <= 135 {
		y--
	}
	if t.angle < 75 {
		x += 3
	}
	if t.angle > 105 {
		x -= 3
	}
	return Position{x, y}
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
