package gorched

import (
	"fmt"
	"math"

	tl "github.com/JoelOtter/termloop"
)

// Bullet is entity representing bullet shooted from tank.
type Bullet struct {
	// tank who shooted this bullet
	shooter *Tank
	// position
	x, y float64
	// velocity
	vx, vy float64
	// true if bullet hit to something
	isInCollision bool
	// will be called when bullet finished his path
	onFinish func()
}

// NewBullet creates new bullet.
// Given onFinish will be  called when bullet finished his path.
func NewBullet(shooter *Tank, p Position, speed float64, angle int, onFinish func()) *Bullet {
	theta := 2.0 * math.Pi * (float64(angle) / 360.0)
	return &Bullet{
		shooter:  shooter,
		x:        float64(p.x),
		y:        float64(p.y),
		vx:       math.Cos(theta) * speed,
		vy:       math.Sin(theta) * -speed,
		onFinish: onFinish,
	}
}

// Draw bullet
func (b *Bullet) Draw(s *tl.Screen) {
	// draw bullet symbol
	s.RenderCell(int(b.x), int(b.y), &tl.Cell{Fg: tl.ColorYellow, Ch: '■'})

	// update velocity by gravity
	b.vy += 9.8100 * s.TimeDelta()

	// update position by velocity
	b.x += b.vx * s.TimeDelta()
	b.y += b.vy * s.TimeDelta()

	// check if below the screen
	sw, sh := s.Size()
	if int(b.y) > sh {
		b.die(s)
		return
	}

	// check if out of screen
	if int(b.y) < 0 || int(b.x) < 0 || int(b.x) > sw {
		x := math.Min(math.Max(0, b.x), float64(sw))
		y := math.Min(math.Max(0, b.y), float64(sh))
		d := math.Sqrt(math.Pow(b.x-x, 2) + math.Pow(b.y-y, 2))
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
	return int(b.x), int(b.y)
}

// Size returns size of collider
func (b *Bullet) Size() (int, int) {
	return 1, 1
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
}
