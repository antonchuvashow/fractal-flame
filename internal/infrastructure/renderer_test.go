package infrastructure

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"hw4-fractal-flame/internal/domain"
)

func TestNewRenderer(t *testing.T) {
	palette := domain.FirePalette()
	gamma := 2.2
	brightness := 1.5

	r := NewRenderer(palette, gamma, brightness)

	assert.NotNil(t, r)
	assert.Equal(t, palette, r.Palette)
	assert.InDelta(t, gamma, r.Gamma, 1e-9)
	assert.InDelta(t, brightness, r.Brightness, 1e-9)
}

func Test_clamp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		value    float64
		expected uint8
	}{
		{"below zero", -10.5, 0},
		{"zero", 0.0, 0},
		{"normal value", 128.7, 128},
		{"255", 255.0, 255},
		{"above 255", 300.1, 255},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, clamp(tc.value))
		})
	}
}

func TestRenderer_Render(t *testing.T) {
	t.Parallel()

	palette := [256]color.RGBA{}
	palette[255] = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	r := NewRenderer(palette, 1.0, 1.0)

	t.Run("renders an empty histogram", func(t *testing.T) {
		t.Parallel()

		hist := domain.NewHistogram(10, 10)
		img := r.Render(hist)

		assert.Equal(t, image.Rect(0, 0, 10, 10), img.Bounds())

		for y := range 10 {
			for x := range 10 {
				assert.Equal(t, color.RGBA{R: 0, G: 0, B: 0, A: 0}, img.At(x, y))
			}
		}
	})

	t.Run("renders a histogram with a single point", func(t *testing.T) {
		t.Parallel()

		hist := domain.NewHistogram(10, 10)
		hist.Increment(5, 5, 1.0) // Increment with color index 1.0

		// With a single point, maxCount is 1, logMaxCount is 0.
		// This leads to division by zero if not handled. Let's set another pixel to avoid this.
		hist.Increment(0, 0, 0.0)
		hist.Increment(0, 0, 0.0) // maxCount is now 2

		hist.Count[5][5] = 2
		hist.ColorSum[5][5] = 2.0
		img := r.Render(hist)
		pixelColor := img.At(5, 5)

		// logDensity = log(2)/log(2) = 1
		// alpha = 1^1 * 2.0 = 2.0
		// Base color is white (255, 255, 255)
		// Final color = (255*2.0, 255*2.0, 255*2.0) clamped to (255, 255, 255)
		expected := color.RGBA{R: 255, G: 255, B: 255, A: 255}
		assert.Equal(t, expected, pixelColor)
	})
}
