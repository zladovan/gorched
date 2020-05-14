package gorched

import (
	"math"

	tl "github.com/JoelOtter/termloop"
	osx "github.com/ojrac/opensimplex-go"
)

// Clouds is entity used to shown moving clouds on the sky.
// Clouds are calculated using open simplex noise function.
type Clouds struct {
	// 2d array of normalized numbers describing "how cloudy" is given pixel
	points [][]float64
	// offset in noise function
	offsetXGlobal float64
	// offset in current points array
	offsetXCurrent int
	// generator of clouds
	generator *CloudsGenerator
}

// CloudsGenerator holds configuration for generating clouds
type CloudsGenerator struct {
	// seed is used as seed to noise function
	seed int64
	// width of visible clouds
	width int
	// height of visible clouds
	height int
	// noise function
	noise osx.Noise
}

// GenerateClouds will initialize noise function in given generator and use it to create new clouds.
func GenerateClouds(g *CloudsGenerator) *Clouds {
	g.noise = osx.NewNormalized(g.seed)
	return &Clouds{points: generate(g, 0), generator: g}
}

func generate(g *CloudsGenerator, offsetX int) [][]float64 {
	// we are generating always 2 times longer cloud array to avoid regenerating it on every move
	points := make([][]float64, g.width*2)
	for x := 0; x < g.width*2; x++ {
		points[x] = make([]float64, g.height)
		for y := 0; y < g.height; y++ {
			// clouds are gradientally fading out with higher y to avoid sharp cut
			// otherwise there are some magic numbers which I just found while trying
			points[x][y] = (1.05 - math.Pow(float64(y), 1.5)/float64(g.height)) * g.noise.Eval2(0.03*float64(x+offsetX), 0.2*float64(y))
		}
	}
	return points
}

// Draw clouds
func (c *Clouds) Draw(s *tl.Screen) {
	for x, columns := range c.points[int(c.offsetXCurrent):] {
		for y, c := range columns {
			switch {
			case c > 0.9:
				s.RenderCell(x, y, &tl.Cell{Fg: tl.ColorWhite, Ch: '▓'})
			case c > 0.7:
				s.RenderCell(x, y, &tl.Cell{Fg: tl.ColorWhite, Ch: '▒'})
			case c > 0.5:
				s.RenderCell(x, y, &tl.Cell{Fg: tl.ColorWhite, Ch: '░'})
			}
		}
	}
	// move clouds
	// TODO: parametrize speed (wind)
	c.offsetXGlobal += 0.5 * s.TimeDelta()
}

// Tick updates clouds points if needed
func (c *Clouds) Tick(e tl.Event) {
	c.offsetXCurrent = int(c.offsetXGlobal) % len(c.points)
	// if we shown already more than half points we need to generate new points
	if c.offsetXCurrent >= len(c.points)/2 {
		c.points = generate(c.generator, int(c.offsetXGlobal))
		c.offsetXCurrent = 0
	}
}
