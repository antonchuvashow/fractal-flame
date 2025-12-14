package domain

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLerp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		a, b, T  float64
		expected float64
	}{
		{"t is zero", 10, 20, 0, 10},
		{"t is one", 10, 20, 1, 20},
		{"t is midpoint", 10, 20, 0.5, 15},
		{"t is negative", 10, 20, -0.5, 5},
		{"t is greater than one", 10, 20, 1.5, 25},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := lerp(tc.a, tc.b, tc.T)
			assert.InDelta(t, tc.expected, result, 1e-9)
		})
	}
}

func TestClampByte(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		value    float64
		expected uint8
	}{
		{"value is negative", -10.5, 0},
		{"value is zero", 0.0, 0},
		{"value is normal", 128.7, 128},
		{"value is 255", 255.0, 255},
		{"value is above 255", 300.1, 255},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := clampByte(tc.value)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFirePalette(t *testing.T) {
	t.Parallel()

	palette := FirePalette()

	assert.Len(t, palette, 256, "Palette should have 256 colors")

	// Check start color (black)
	assert.Equal(t, color.RGBA{R: 0, G: 0, B: 0, A: 255}, palette[0])

	// Check a color in the middle (yellowish)
	// t = 127/255 ~= 0.498
	// R = 255
	// G = (0.498 - 0.33) * 3 * 255 ~= 128
	// B = 0
	assert.Equal(t, uint8(255), palette[127].R)
	assert.InDelta(t, 128, palette[127].G, 2)
	assert.Equal(t, uint8(0), palette[127].B)
	assert.Equal(t, uint8(255), palette[127].A)

	// Check end color (white)
	assert.Equal(t, color.RGBA{R: 255, G: 255, B: 255, A: 255}, palette[255])
}

func TestRainbowPalette(t *testing.T) {
	t.Parallel()

	palette := RainbowPalette()

	assert.Len(t, palette, 256, "Palette should have 256 colors")

	// Check key colors in the rainbow spectrum
	assert.Equal(t, color.RGBA{R: 255, G: 0, B: 0, A: 255}, palette[0], "Should start with Red")
	assert.Equal(t, color.RGBA{R: 0, G: 255, B: 0, A: 255}, palette[85], "Should have Green around 1/3")
	assert.Equal(t, color.RGBA{R: 0, G: 0, B: 255, A: 255}, palette[170], "Should have Blue around 2/3")

	// The end of the rainbow spectrum (h=6) wraps back to red.
	assert.Equal(t, color.RGBA{R: 255, G: 0, B: 0, A: 255}, palette[255], "Should end with Red")

	// Check alpha on a random color
	assert.Equal(t, uint8(255), palette[100].A)
}

func TestOceanPalette(t *testing.T) {
	t.Parallel()

	palette := OceanPalette()

	assert.Len(t, palette, 256, "Palette should have 256 colors")

	// Check start color (dark blue)
	assert.Equal(t, color.RGBA{R: 5, G: 10, B: 40, A: 255}, palette[0])

	// Check color in the middle
	assert.Equal(t, color.RGBA{R: 29, G: 99, B: 169, A: 255}, palette[127])

	// Check end color (light cyan/white)
	assert.Equal(t, color.RGBA{R: 180, G: 230, B: 255, A: 255}, palette[255])
}
