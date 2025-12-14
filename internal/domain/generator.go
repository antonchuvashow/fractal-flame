package domain

import (
	"image/color"
	"math"
	"math/rand"
)

// First iterations that will be skipped.
const warmupIterations = 20

type Generator struct {
	// Transforms is a set of weighted non-linear composite transformations.
	Transforms []Transform

	// Palette is a 256 color array used for colorizing the picture.
	Palette [256]color.RGBA

	// SymmetryLevel is a number of duplicates around the center
	SymmetryLevel int

	// totalWeight is the sum of all transform weights, used for random selection.
	totalWeight float64
}

// NewGenerator creates and initializes a new fractal generator.
func NewGenerator(variations []Transform, palette [256]color.RGBA, symmetryLevel int) *Generator {
	var totalWeight float64
	for _, t := range variations {
		totalWeight += t.Weight()
	}
	return &Generator{
		Transforms:    variations,
		Palette:       palette,
		SymmetryLevel: symmetryLevel,
		totalWeight:   totalWeight,
	}
}

// Generate runs the Chaos Game algorithm to produce a fractal histogram.
// It takes the image dimensions, number of iterations, and a random number source.
func (g *Generator) Generate(width, height, iterations int, rnd *rand.Rand) *Histogram {
	hist := NewHistogram(width, height)
	p := Point{rnd.Float64()*2 - 1, rnd.Float64()*2 - 1, rnd.Float64()}
	angles := g.precomputeAngles()

	for i := -warmupIterations; i < iterations; i++ {
		transform := g.chooseRandomTransform(rnd)

		// Apply the chosen non-linear transform, then the final affine transform.
		p = transform.Apply(p)
		p.ColorIndex = (p.ColorIndex + transform.ColorIndex()) / 2

		if i >= 0 {
			g.plotWithSymmetry(hist, p, width, height, angles)
		}
	}

	return hist
}

// chooseRandomTransform selects a transformation based on its weight.
// Transformations with higher weights are chosen more frequently.
func (g *Generator) chooseRandomTransform(rnd *rand.Rand) Transform {
	r := rnd.Float64() * g.totalWeight

	var currentWeight float64
	for _, t := range g.Transforms {
		currentWeight += t.Weight()
		if r < currentWeight {
			return t
		}
	}

	// Fallback for the last transform just in case.
	return g.Transforms[len(g.Transforms)-1]
}

// precomputeAngles get sin and cos of the division.
func (g *Generator) precomputeAngles() []struct{ sin, cos float64 } {
	angles := make([]struct{ sin, cos float64 }, g.SymmetryLevel)

	for i := range g.SymmetryLevel {
		angle := 2.0 * math.Pi * float64(i) / float64(g.SymmetryLevel)
		angles[i].sin = math.Sin(angle)
		angles[i].cos = math.Cos(angle)
	}

	return angles
}

// plotWithSymmetry plot the point with given angles.
func (g *Generator) plotWithSymmetry(hist *Histogram, point Point,
	width, height int, angles []struct{ sin, cos float64 }) {

	for _, a := range angles {
		rotatedX := point.X*a.cos - point.Y*a.sin
		rotatedY := point.X*a.sin + point.Y*a.cos

		px := int(float64(width) * (rotatedX + 1) / 2)
		py := int(float64(height) * (rotatedY + 1) / 2)

		if hist.Contains(px, py) {
			hist.Increment(px, py, point.ColorIndex)
		}
	}
}
