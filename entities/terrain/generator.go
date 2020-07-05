package terrain

import (
	osx "github.com/ojrac/opensimplex-go"
)

// Generator holds configuration for terrain generating logic
type Generator struct {
	// Seed to the noise function
	Seed int64
	// Width of the terrain
	Width int
	// Height of the terrain holds maximum height of terrain hill
	Height int
	// Roughness configures how much will be terrain "wavy"
	Roughness float64
	// LowColor generates terrain in only 8 colors mode when true
	LowColor bool
}

// Generate will generate new terrain using noise function (open simplex)
func Generate(g *Generator) *Terrain {
	noise := osx.NewNormalized(g.Seed)
	heights := make([]int, g.Width)
	for x := 0; x < g.Width; x++ {
		// reduce height to keep 5 cells space for tank on the highest hill top
		heights[x] = 5 + int(float64(g.Height-5)*noise.Eval2(g.Roughness/float64(g.Width)*float64(x), 0.5))
	}
	return NewTerrain(heights, g.Height, g.LowColor)
}
