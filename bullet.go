package gorched

import (
	"fmt"
	"math"

	tl "github.com/JoelOtter/termloop"
	"github.com/zladovan/gorched/gmath"
)

// Bullet is entity representing bullet shooted from tank.
type Bullet struct {
	// tank who shooted this bullet
	shooter *Tank
	// body is physical body
	body *Body
	// strength of the explosion
	strength int
	// true if bullet hit to something
	isInCollision bool
	// will be called when bullet finished his path
	onFinish func()
}

// NewBullet creates new bullet.
// Given onFinish will be  called when bullet finished his path.
func NewBullet(shooter *Tank, p Position, speed float64, angle int, strength int, onFinish func()) *Bullet {
	theta := 2.0 * math.Pi * (float64(angle) / 360.0)
	return &Bullet{
		shooter: shooter,
		body: &Body{
			Position: gmath.Vector2f{X: float64(p.x), Y: float64(p.y)},
			Velocity: gmath.Vector2f{X: math.Cos(theta) * speed, Y: math.Sin(theta) * -speed},
			Mass:     1,
		},
		strength: strength,
		onFinish: onFinish,
	}
}

// Draw bullet
func (b *Bullet) Draw(s *tl.Screen) {
	// draw bullet symbol
	s.RenderCell(int(b.body.Position.X), int(b.body.Position.Y), &tl.Cell{Fg: tl.ColorYellow, Ch: '■'})

	// check if below the screen
	sw, sh := s.Size()
	if int(b.body.Position.Y) > sh {
		b.die(s)
		return
	}

	// check if out of screen
	if int(b.body.Position.Y) < 0 || int(b.body.Position.X) < 0 || int(b.body.Position.X) > sw {
		x := math.Min(math.Max(0, b.body.Position.X), float64(sw))
		y := math.Min(math.Max(0, b.body.Position.Y), float64(sh))
		d := math.Sqrt(math.Pow(b.body.Position.X-x, 2) + math.Pow(b.body.Position.Y-y, 2))
		dstr := fmt.Sprintf("%d", int(d))

		if x >= float64(sw) {
			x -= float64(len(dstr))
		}

		// draw number with how far is bullet out of screen
		// TODO: use label for this info
		i := 0
		for _, c := range dstr {
			s.RenderCell(int(x)+i, int(y), &tl.Cell{Fg: tl.ColorYellow | tl.AttrBold, Ch: c})
			i++
		}
	}

	// if bullet hit somewhere it's dead
	if b.isInCollision {
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
	b.onFinish()
}

// Collide check the collisions
func (b *Bullet) Collide(collision tl.Physical) {
	b.isInCollision = true
	if target, ok := collision.(*Tank); ok {
		target.TakeDamage()
		if target != b.shooter {
			b.shooter.Hit()
		}
	}
	if target, ok := collision.(*TerrainColumn); ok {
		bx := int(b.body.Position.X)
		by := int(b.body.Position.Y)
		Debug.Logf("Ground was hit x=%d y=%d", bx, by)
		target.terrain.MakeHole(bx, by, b.strength)
	}
}

// Body returns physical body of this bullet
func (b *Bullet) Body() *Body {
	return b.body
}
