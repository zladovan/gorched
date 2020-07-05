package entities

import (
	"math"
	"time"

	tl "github.com/JoelOtter/termloop"
	osx "github.com/ojrac/opensimplex-go"
	"github.com/zladovan/gorched/debug"
	"github.com/zladovan/gorched/draw"
	"github.com/zladovan/gorched/entities/terrain"
	"github.com/zladovan/gorched/gmath"
)

// Explosion represent effect of explosion.
type Explosion struct {
	// Center defines point where explosion is starting
	Center gmath.Vector2i
	// Strength defines maximum radius which explosion can take
	Strength float64
	// radius is actual radius of explosion
	radius float64
	// speed is given in number of explosion cycles per second
	// one explosion cycle contains explosion and implosion
	speed float64
	// t is time in nanoseconds since explosion was created
	t float64
	// noise holds noise function used for calculating explosion pattern
	noise osx.Noise
	// terrainCollided is flag for marking that collision with terrain was already applied
	terrainCollided bool
	// collided contains flag for each entity which explosion already collided with
	collided map[tl.Physical]bool
	// shooter is tank who caused this explosion and will be rewarded if this explosion will take some damage, can be nil
	shooter *Tank
}

// NewExplosion creates new explosion in the given center point.
// Given strength defines maximum radius which explosion reaches at it's peak.
// Optionally you can specify shooter to tank who caused this explosion and will be rewarded if this explosion will take some damage.
func NewExplosion(center gmath.Vector2i, strength int, shooter *Tank) *Explosion {
	return &Explosion{
		Center:   center,
		Strength: float64(strength),
		speed:    1,
		noise:    osx.NewNormalized(time.Now().UTC().UnixNano()),
		collided: map[tl.Physical]bool{},
		shooter:  shooter,
	}
}

// Draw is drawing explosion sprite.
func (e *Explosion) Draw(s *tl.Screen) {
	// increase time of explosion
	e.t += s.TimeDelta()

	// radius is growing with time up to maximum  given by strength and after then it's decreasing
	e.radius = math.Sin(e.speed*math.Pi*e.t) * (e.Strength + 1)
	if e.t > 1/e.speed {
		s.Level().RemoveEntity(e)
		return
	}

	// gradient with shades of yellow
	gradient := draw.RadialGradient{
		A:         e.Center,
		B:         *e.Center.Translate(0, int(e.Strength)),
		ColorA:    tl.Attr(232),
		Step:      -1,
		StepCount: 3,
	}

	// we need integer version of radius
	r := int(e.radius)

	// maximal distance of explosion points for current radius
	maxd := e.Center.DistanceF(e.Center.Translate(r, r))

	// for each x we will print one column
	for ix := -r + 1; ix < r; ix++ {
		// radius on y axis for current x is scale by 0.5 to reduce terminal's cells ratio 2:1 for height:width
		ry := int(math.Sqrt(math.Pow(float64(r-1), 2)-math.Pow(float64(ix), 2)) * 0.5)

		for iy := -ry; iy <= ry; iy++ {
			// vector used as gradient and noise intput
			// it is scaled back by 1/0.5 to do not break gradient and noise patterns
			// it is rotated around center to produce rotating effect of explosion
			v := e.Center.Translate(ix, int(float64(iy)/0.5)).RotateAround(&e.Center, e.t).As2F()

			// init cell as point with maximum intensity (center of explosion)
			cell := &tl.Cell{Fg: gradient.Color(*v.As2I()), Ch: '█'}
			red := tl.Attr(221)
			if IsLowColor(s) {
				cell.Fg = tl.ColorYellow
				red = tl.ColorRed
			}

			// structure of explosion pattern should be affected by distance from center
			// if it's in the center it should be no chance to use other cell then initial
			// more far from center more change to use less visible cells (or cells with red background near center)
			n := e.noise.Eval2(v.X*0.2, v.Y*0.2) * v.Distance(e.Center.As2F()) / maxd * 2
			minn := 0.4 // higher = bigger and stronger center
			switch {
			case n > minn+(1-minn)*0.875:
				cell.Ch = ' '
			case n > minn+(1-minn)*0.75:
				cell.Ch = '░'
			case n > minn+(1-minn)*0.625:
				cell.Ch = '▒'
				cell.Bg = red
			case n > minn+(1-minn)*0.5:
				cell.Fg = red
			case n > minn:
				cell.Ch = '▓'
				cell.Bg = red
			}

			// draw
			s.RenderCell(e.Center.X+ix, e.Center.Y+iy, cell)
		}

	}
}

// Position returns postion of collider
func (e *Explosion) Position() (int, int) {
	return e.Center.X - int(e.radius) + 1, e.Center.Y - int(e.radius/2) + 1
}

// Size returns size of collider
func (e *Explosion) Size() (int, int) {
	return int(e.radius)*2 - 1, int(e.radius) - 1
}

// Collide hadnles collisions with other objects
func (e *Explosion) Collide(collision tl.Physical) {
	// do not care about collisions if explosion is not yet after peak
	if !e.afterPeak() {
		return
	}

	// process collision with terrain only once
	if !e.terrainCollided {
		if target, ok := collision.(*terrain.Column); ok {
			debug.Logf("Explosion collides with terrain")
			e.terrainCollided = true
			target.MakeHole(e.Center.X, e.Center.Y, int(e.Strength))
		}
	}

	// process collisions with each tank only once
	if !e.collided[collision] {
		if target, ok := collision.(*Tank); ok {
			// get middle point of tank
			tx, ty := target.Position()
			tw, th := target.Size()
			tp := &gmath.Vector2f{X: float64(tx) + float64(tw-1)/2, Y: float64(ty) + float64(th-1)/2}

			// calculate distance of explosion center from tank's middle point
			d := e.Center.As2F().Distance(tp.Translate(0, -(float64(e.Center.Y)-tp.Y)/2))

			// damage to be taken is affected by distance of explosion and by the strength of explosion
			damage := int(math.Max(0, (e.Strength-d)/e.Strength) * e.MaxDamage())
			target.TakeDamage(damage, e.shooter)
			debug.Logf("Explosion collides with tank damage=%d", damage)
		}
		e.AddAlreadyCollided(collision)
	}
}

// AddAlreadyCollided will make ignore later collisions with this explosion and given object p
// This is useful if you do not want to cause any or any additional damage to p by this explosion.
func (e *Explosion) AddAlreadyCollided(p tl.Physical) {
	e.collided[p] = true
}

// MaxDamage returns maximum amount of damage which can be taken by this explosion
func (e *Explosion) MaxDamage() float64 {
	return 100 + e.Strength*10
}

// afterPeak returns true if explosion was already in it's biggest radius
func (e *Explosion) afterPeak() bool {
	return e.t > 1/e.speed/2
}

// Tick is not used now
func (e *Explosion) Tick(ev tl.Event) {}

// ZIndex returns z-index of explosion which should be higher than most of other z-indexes
func (e *Explosion) ZIndex() int {
	return 10001
}
