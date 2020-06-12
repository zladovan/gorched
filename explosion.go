package gorched

import (
	"math"
	"time"

	tl "github.com/JoelOtter/termloop"
	osx "github.com/ojrac/opensimplex-go"
	"github.com/zladovan/gorched/draw"
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
}

// NewExplosion creates new explosion in the given center point.
// Given strength defines maximum radius which explosion reaches at it's peak.
func NewExplosion(center gmath.Vector2i, strength int) *Explosion {
	return &Explosion{
		Center:   center,
		Strength: float64(strength),
		speed:    1,
		noise:    osx.NewNormalized(time.Now().UTC().UnixNano()),
	}
}

// Draw is drawing explosion sprite.
func (e *Explosion) Draw(s *tl.Screen) {
	if e.t == 0 {
		Debug.Logf("Explosion started")
	}
	// increase time of explosion
	e.t += s.TimeDelta()

	// radius is growing with time up to maximum  given by strength and after then it's decreasing
	e.radius = math.Sin(e.speed*math.Pi*e.t) * (e.Strength + 1)
	if e.t > 1/e.speed {
		s.Level().RemoveEntity(e)
		Debug.Logf("Explosion finished")
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
	return e.Center.X, e.Center.Y
}

// Size returns size of collider
func (e *Explosion) Size() (int, int) {
	return int(e.radius), int(e.radius)
}

// Collide hadnles collisions with other objects
func (e *Explosion) Collide(collision tl.Physical) {
	// do not care about collisions if explosion is not yet after peak
	if !e.afterPeak() {
		return
	}

	// process collision with terrain only once
	if !e.terrainCollided {
		if target, ok := collision.(*TerrainColumn); ok {
			Debug.Logf("Explosion collides with terrain")
			e.terrainCollided = true
			target.terrain.MakeHole(e.Center.X, e.Center.Y, int(e.Strength))
		}
	}
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
