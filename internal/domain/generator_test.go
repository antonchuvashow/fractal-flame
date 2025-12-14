package domain

import (
	"image/color"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGenerator(t *testing.T) {
	t.Run("creates a new generator with correct properties", func(t *testing.T) {
		transforms := []Transform{
			Affine{W: 1.0},
			Sinusoidal{W: 2.0},
		}

		var palette [256]color.RGBA

		symmetryLevel := 3
		g := NewGenerator(transforms, palette, symmetryLevel)

		assert.Equal(t, transforms, g.Transforms)
		assert.Equal(t, palette, g.Palette)
		assert.Equal(t, symmetryLevel, g.SymmetryLevel)
		assert.InDelta(t, 3.0, g.totalWeight, 0.01)
	})
}

func TestGenerator_Generate(t *testing.T) {
	t.Run("runs the chaos game and produces a histogram", func(t *testing.T) {
		transforms := []Transform{
			Affine{A: 0.5, E: 0.5, W: 1.0, Color: 0.2},
		}
		g := NewGenerator(transforms, [256]color.RGBA{}, 1)
		rnd := rand.New(rand.NewSource(99))
		iterations := 100

		hist := g.Generate(100, 100, iterations, rnd)

		// Check that the histogram is not empty and has some hits
		var totalHits uint32

		for i := range hist.Height {
			for j := range hist.Width {
				totalHits += hist.GetCount(i, j)
			}
		}

		assert.Positive(t, totalHits, "Expected histogram to have some hits")
		assert.Equal(t, iterations, int(totalHits))
	})
}
