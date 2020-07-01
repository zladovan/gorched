package gorched

import (
	"fmt"
	"math"

	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/debug"
	"github.com/zladovan/gorched/gmath"
	"github.com/zladovan/gorched/physics"
	"github.com/zladovan/gorched/terrain"
)

// Bullet is entity representing bullet shooted from tank.
type Bullet struct {
	// tank who shooted this bullet
	shooter *Tank
	// body is physical body
	body *physics.Body
	// strength of the explosion
	strength int
	// explosion is created after bullet hit to something
	explosion *Explosion
}

// NewBullet creates new bullet.
func NewBullet(shooter *Tank, p gmath.Vector2i, speed float64, angle int, strength int) *Bullet {
	theta := 2.0 * math.Pi * (float64(angle) / 360.0)
	return &Bullet{
		shooter: shooter,
		body: &physics.Body{
			Position: gmath.Vector2f{X: float64(p.X), Y: float64(p.Y)},
			Velocity: gmath.Vector2f{X: math.Cos(theta) * speed, Y: math.Sin(theta) * -speed},
			Mass:     1,
		},
		strength: strength,
	}
}

// Draw bullet
func (b *Bullet) Draw(s *tl.Screen) {
	// color of the bullet
	color := tl.Attr(221)
	if IsLowColor(s) {
		color = tl.ColorYellow
	}

	// draw bullet symbol
	s.RenderCell(int(b.body.Position.X), int(b.body.Position.Y), &tl.Cell{Fg: color, Ch: '■'})

	// remove if below the screen or too far on the left/right of the screen
	sw, sh := s.Size()
	if int(b.body.Position.Y) > sh || int(b.body.Position.X) < -100 || int(b.body.Position.X) > sw+100 {
		b.die(s)
		return
	}

	// check if out of screen
	if int(b.body.Position.Y) < 0 || int(b.body.Position.X) < 0 || int(b.body.Position.X) > sw {
		x := gmath.Clampf(0, float64(sw), b.body.Position.X)
		y := gmath.Clampf(0, float64(sh), b.body.Position.Y)
		d := b.body.Position.Translate(-x, -y).Length()
		dstr := fmt.Sprintf("%d", int(d))

		// adjust x position to do not draw number out of screen
		maxx := x + float64(len(dstr))
		if maxx >= float64(sw) {
			x -= maxx - float64(sw)
		}

		// draw number with how far is bullet out of screen
		// TODO: use label for this info
		i := 0
		for _, c := range dstr {
			s.RenderCell(int(x)+i, int(y), &tl.Cell{Fg: color | tl.AttrBold, Ch: c})
			i++
		}
	}

	// if bullet hit somewhere it's dead
	if b.explosion != nil {
		s.Level().AddEntity(b.explosion)
		b.die(s)
	}
}

// Tick is not used yet
func (b *Bullet) Tick(e tl.Event) {}

// Position returns postion of collider
func (b *Bullet) Position() (int, int) {
	return int(b.body.Position.X), int(b.body.Position.Y)
}

// Size returns size of collider
func (b *Bullet) Size() (int, int) {
	return 1, 1
}

// ZIndex return z-index of bullet.
// It should be higher as in most other entities.
func (b *Bullet) ZIndex() int {
	return 10000
}

// bullet finished his path
func (b *Bullet) die(s *tl.Screen) {
	s.Level().RemoveEntity(b)
}

// Collide check the collisions
func (b *Bullet) Collide(collision tl.Physical) {
	b.explosion = NewExplosion(*b.body.Position.As2I(), b.strength+3, b.shooter)
	b.body.Locked = true

	// collision with tank
	if target, ok := collision.(*Tank); ok {
		target.TakeDamage(int(b.explosion.MaxDamage()), b.shooter)
		b.explosion.AddAlreadyCollided(target)
	}

	// colision with terrain
	if _, ok := collision.(*terrain.Column); ok {
		bx := int(b.body.Position.X)
		by := int(b.body.Position.Y)
		debug.Logf("Ground was hit x=%d y=%d", bx, by)
	}
}

// Body returns physical body of this bullet
func (b *Bullet) Body() *physics.Body {
	return b.body
}
