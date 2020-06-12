package gorched

import (
	"math"
	"strings"

	tl "github.com/JoelOtter/termloop"
	osx "github.com/ojrac/opensimplex-go"

	"github.com/zladovan/gorched/draw"
	"github.com/zladovan/gorched/gmath"
)

// Tree represents tree in the scenery.
// Use NewTree for creating new Tree.
// Use GenerateWood for random generation of multiple Trees.
type Tree struct {
	// it extends from termloop.Entity
	*tl.Entity
	// body is physical body of the tree used for falling simulation
	body *Body
}

// TreeKind represents the type of tree.
// Different sprites are drawn based on TreeKind.
type TreeKind uint8

const (
	// SpruceTree looks like christmas tree.
	SpruceTree TreeKind = iota
	// OakTree has oval crown.
	OakTree
	// PopulusTree is thin and long tree.
	PopulusTree

	// CountOfTreeKind holds the count of all different elements of TreeKind enum.
	// It must be always last element!
	CountOfTreeKind
)

// NewTree creates new Tree
func NewTree(position gmath.Vector2i, kind TreeKind, size int, lowColor bool, asciiOnly bool) *Tree {
	canvas := createTreeCanvas(kind, size, lowColor, asciiOnly)
	return &Tree{
		Entity: tl.NewEntityFromCanvas(position.X-len(canvas)/2, position.Y-len(canvas[0]), canvas),
		body: &Body{
			Position: *position.As2F(),
			Mass:     5,
		},
	}
}

// Create sprite for tree of given kind
func createTreeCanvas(kind TreeKind, size int, lowColor bool, asciiOnly bool) tl.Canvas {
	switch kind {
	case SpruceTree:
		return createSpruceTreeCanvas(size, lowColor, asciiOnly)
	case PopulusTree:
		return createPopulusTreeCanvas(size, lowColor, asciiOnly)
	case OakTree:
		return createOakTreeCanvas(size, lowColor)
	}
	panic("Invalid tree kind")
}

// Create sprite for tree of kind SpruceTree with given size.
// Spruce tree is growing in width with size and it's trunk is always only 1 cell high.
//
// See below for how it's shape looks for size 1 - 3:
//
//    ▓
//   ▓▓▓
//    }
//
//    ▓
//   ▓▓▓
//  ▓▓▓▓▓
//    █
//
//    ▓
//   ▓▓▓
//  ▓▓▓▓▓
// ▓▓▓▓▓▓▓
//    █
//
func createSpruceTreeCanvas(size int, lowColor bool, asciiOnly bool) tl.Canvas {
	// calculate dimensions
	trunkHeight := 1
	crownHeight := 1 + size
	height := trunkHeight + crownHeight
	width := 3 + (size-1)*2

	// create printer with canvas
	p := draw.BlankPrinter(width, height)

	// print trunk
	if size > 1 {
		if lowColor {
			p.Fg = tl.ColorMagenta | tl.AttrBold
			p.Bg = p.Fg
		} else {
			p.Fg = 101
			p.Bg = 95
		}
		p.WriteHorizontalUp(p.CenterX(), p.MaxY(), strings.Repeat("}", trunkHeight))
	} else {
		if lowColor {
			p.Fg = tl.ColorMagenta | tl.AttrBold
		} else {
			p.Fg = 95 | tl.AttrBold
		}
		p.WriteCenterX(p.MaxY(), "}")
	}

	// gradient with shades of green for crown
	gradient := draw.RadialGradient{
		A:         gmath.Vector2i{X: width - 2 - (size-1)/2, Y: 0},
		B:         gmath.Vector2i{X: 0, Y: crownHeight - 1 + (size+1)/4},
		ColorA:    tl.Attr(42),
		Step:      -6,
		StepCount: 3,
	}

	// character used for drawing a crown pixel
	crownCh := "≡"
	if asciiOnly {
		crownCh = "="
	}

	// draw crown
	for i := 0; i < crownHeight; i++ {
		// width (w) of crown is increasing linearly with size
		w := 1 + i*2
		for j := 0; j < w; j++ {
			v := gmath.Vector2i{X: width/2 - w/2 + j, Y: i}
			if lowColor {
				p.Fg = tl.ColorCyan | tl.AttrBold
				p.Bg = tl.ColorDefault
			} else {
				p.Fg = gradient.Color(v) | tl.AttrBold
				p.Bg = tl.Attr(24)
			}
			p.Write(v.X, v.Y, crownCh)
		}
	}

	return *p.Canvas
}

// Create sprite for tree of kind PopulusTree with given size.
// Pupulus tree is growing in heigh with size while it's width is same (except size 1).
// I's trunk is growing by one cell for each next 4 in size.
//
// See below  fot how it looks for sizes 1 - 5:
//
//  M
//  W
//  |
//
// /M\
// \W/
//  █
//
// /M\
// UUU
// \W/
//  █
//
// /M\
// UUU
// UUU
// \W/
//  █
//  █
//
// /M\
// UUU
// UUU
// UUU
// \W/
//  █
//  █
//
func createPopulusTreeCanvas(size int, lowColor bool, asciiOnly bool) tl.Canvas {
	// calculate dimensions
	trunkHeight := 1 + size/4
	crownHeight := 2
	if size > 1 {
		crownHeight = size
	}
	height := trunkHeight + crownHeight

	// character used for drawing trunk pixel
	trunkCh := "≡"
	if asciiOnly {
		trunkCh = "="
	}

	// create printer with canvas
	p := draw.BlankPrinter(3, height)

	// print trunk
	if size > 1 {
		if lowColor {
			p.Fg = tl.ColorWhite | tl.AttrBold
		} else {
			p.Fg = 243
			p.Bg = 245
		}
		p.WriteHorizontalUp(p.CenterX(), p.MaxY(), strings.Repeat(trunkCh, trunkHeight))
	} else {
		if lowColor {
			p.Fg = tl.ColorWhite | tl.AttrBold
		} else {
			p.Fg = 245 | tl.AttrBold
		}
		p.WriteCenterX(p.MaxY(), "|")
	}

	// print crown
	for i := 0; i < crownHeight; i++ {
		var texture string
		switch i {
		case 0:
			// upper part of crown
			texture = "/M\\"
		case crownHeight - 1:
			// bottom part of crown
			texture = "\\W/"
		default:
			// middle part of crown
			texture = "bUd"
		}
		if lowColor {
			p.Fg = tl.ColorGreen | tl.AttrBold
			p.Bg = tl.ColorDefault
		} else {
			p.Bg = tl.Attr(29)
			p.Fg = tl.Attr(35)
		}
		if size > 1 {
			// side columns of crown
			p.Write(0, i, string(texture[0]))
			p.Write(2, i, string(texture[2]))
		}
		// middle column of crown
		if !lowColor {
			p.Fg = tl.Attr(41)
		}
		p.Write(1, i, string(texture[1]))
	}

	return *p.Canvas
}

// Create sprite for tree of kind OakTree with given size.
// Oak tree has round crown which is growing with size.
// I's trunk is growing by one cell for each next 3 in size after size higher then 4.
//
// See below  fot how it looks for sizes 1 - 5:
//
//
//     ▓▓▓
//      /
//
//     ▓▓▓
//    ▓▓▓▓▓
//     ▓▓▓
//      █
//
//    ▓▓▓▓▓
//   ▓▓▓▓▓▓▓
//    ▓▓▓▓▓
//      █
//
//    ▓▓▓▓▓
//   ▓▓▓▓▓▓▓
//  ▓▓▓▓▓▓▓▓▓
//   ▓▓▓▓▓▓▓
//    ▓▓▓▓▓
//      █
//      █
//
//   ▓▓▓▓▓▓▓
//  ▓▓▓▓▓▓▓▓▓
// ▓▓▓▓▓▓▓▓▓▓▓
//  ▓▓▓▓▓▓▓▓▓
//   ▓▓▓▓▓▓▓
//      █
//      █
//
func createOakTreeCanvas(size int, lowColor bool) tl.Canvas {
	// calculate dimensions
	trunkHeight := 1 + ((size - 1) / 3)
	crownHeight := 1 + (size/2)*2
	height := trunkHeight + crownHeight
	width := 3 + (size-1)*2

	// create printer with canvas
	p := draw.BlankPrinter(width, height)

	// print trunk
	if size > 1 {
		if lowColor {
			p.Fg = tl.ColorMagenta | tl.AttrBold
		} else {
			p.Fg = 245
			p.Bg = 95
		}
		p.WriteHorizontalUp(p.CenterX(), p.MaxY(), strings.Repeat("║", trunkHeight))
	} else {
		if lowColor {
			p.Fg = tl.ColorMagenta | tl.AttrBold
		} else {
			p.Fg = 95 | tl.AttrBold
		}
		p.WriteCenterX(p.MaxY(), "/")
	}

	// gradient with shades of green for the crown
	gradient := draw.RadialGradient{
		A:         gmath.Vector2i{X: width - 1, Y: 0},
		B:         gmath.Vector2i{X: 0, Y: crownHeight - 1},
		ColorA:    tl.Attr(41),
		Step:      -6,
		StepCount: 3,
	}

	// print crown
	p.Bg = tl.Attr(23)
	for i := 0; i < crownHeight; i++ {
		// width (w) of crown is increasing by two until half of height, after half is decreasing by two
		// when crown height is even, width of crown is same on two rows around half height
		w := width - int(math.Abs(float64(crownHeight/2-i)))*2
		for j := 0; j < w; j++ {
			v := gmath.Vector2i{X: width/2 - w/2 + j, Y: i}
			if lowColor {
				p.Fg = tl.ColorGreen | tl.AttrBold
				p.Bg = tl.ColorDefault
			} else {
				p.Fg = gradient.Color(v) | tl.AttrBold
				p.Bg = tl.Attr(23)
			}
			p.Write(v.X, v.Y, "@")
		}
	}

	return *p.Canvas
}

// Draw only updates entity position based on physical body and draws it
func (t *Tree) Draw(s *tl.Screen) {
	w, h := t.Entity.Size()
	t.Entity.SetPosition(int(t.body.Position.X)-w/2, int(t.body.Position.Y)-h)
	t.Entity.Draw(s)
}

// Size returns 0 to make trees not collidable yet
func (t *Tree) Size() (int, int) {
	return 0, 0
}

// ZIndex return z-index of tree.
// It should be bigger than z-index of terrain and lower than z-index of tank.
// Trees with higher position (lower y) will have lower z-index to be far from the screen.
// Sorting trees by y position will avoid weird looking trunks over crowns.
func (t *Tree) ZIndex() int {
	_, y := t.Position()
	_, h := t.Entity.Size()
	return 1000 + y + h
}

// Body returns physical body of the tank used for falling simulation
func (t *Tree) Body() *Body {
	return t.body
}

// BottomLine returns line x coordinates for collision with the ground when falling
func (t *Tree) BottomLine() (int, int) {
	// just single point (trunk)
	return 0, 0
}

// Wood is just array of Tree
type Wood []*Tree

// WoodGenerator holds some parameters which control how wood will look like
type WoodGenerator struct {
	// Line is terrain line where index is x coordinate and value at this index is height.
	// All trees positions will be within this line.
	Line []int
	// Seed for the noise function
	Seed int64
	// Density controls how much trees there will be.
	// It should be number between 0 and 1.
	// If 0 is used there will be no trees.
	// If 1 is used there will be tree on every available position (considering MinSpace parameter).
	Density float64
	// MaxSize controls the maximus size of generated tree.
	// It should be bigger than 0.
	MaxSize uint
	// MinSpace controls how far (on x axis) need to be trees minimally from each other.
	MinSpace uint
	// LowColor generates trees in only 8 colors mode when true
	LowColor bool
	// ASCIIOnly identifies that only ASCII characters can be used for generated trees
	ASCIIOnly bool
}

// constants which are too magic to become generator paramters
const (
	// woodMagicGrouping allows control how tight groups of trees are
	// when lower values are used (e.g. 0.1) then trees are grouped into groups with tree size changing more gradient-like
	woodMagicGrouping = 0.875
	// woodMagicGrouping allows to control how large will be continuous areas with same kind of tree
	// when higher values are used (close to 1) then more likely all kind of trees will be generated in one generation and trees will look more like mixed
	// when lower values are used (e.g. 0.01) then more likely there will be larger separated groups or just only one type of trees in one generation
	woodMagicKindVariability = 0.01
)

// GenerateWood generates trees for given generator g
func GenerateWood(g *WoodGenerator) Wood {
	wood := Wood{}

	// zero density means no trees
	if g.Density <= 0 {
		return wood
	}

	//init noise
	noise := osx.NewNormalized(g.Seed)

	// threshold is the minimal value for which there will be a tree
	threshold := 1 - g.Density
	if threshold < 0 {
		threshold = 0
	}

	// just some random point for noise evaluation coordinates which are not changed during one generation
	r := float64(g.Line[0]) / float64(g.Line[len(g.Line)-1])

	// init last assigned x to something out of range for minimal space
	lastX := math.MinInt8

	for x, y := range g.Line {
		// check if we are far enough from last tree
		if x-lastX <= int(g.MinSpace) {
			continue
		}

		// size and kind are separately evaluated from noise function to be independent
		// size is used to determine if tree will be generated or not depending on if it's higher than threshold
		size := noise.Eval2(woodMagicGrouping*float64(x), r)
		if size > threshold {
			size := (size-threshold)/(1-threshold)*float64(g.MaxSize) + 1
			kind := TreeKind(noise.Eval3(r, r, woodMagicKindVariability*float64(x)) * float64(CountOfTreeKind))
			tree := NewTree(gmath.Vector2i{x, y}, kind, int(size), g.LowColor, g.ASCIIOnly)
			wood = append(wood, tree)
			lastX = x
		}
	}

	return wood
}

// CutAround return new Wood without trees which were in collision with rectangle defined by given x, y coordinates, width and height
func (w Wood) CutAround(x, y, width, heigh int) Wood {
	cut := Wood{}
	for _, t := range w {
		tx, ty := t.Position()
		tw, th := t.Entity.Size()
		if tx <= x+width && tx+tw >= x && ty <= y+heigh && ty+th >= y {
			continue
		}
		cut = append(cut, t)
	}
	return cut
}
