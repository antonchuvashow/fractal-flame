package infrastructure

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"hw4-fractal-flame/internal/domain"
)

func TestGetTransformByName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		transformName  string
		weight         float64
		colorIndex     float64
		expectedType   interface{}
		expectError    bool
		expectedErrMsg string
	}{
		{"sinusoidal", "sinusoidal", 1.0, 0.5, domain.Sinusoidal{}, false, ""},
		{"spherical", "spherical", 1.0, 0.5, domain.Spherical{}, false, ""},
		{"swirl", "swirl", 1.0, 0.5, domain.Swirl{}, false, ""},
		{"horseshoe", "horseshoe", 1.0, 0.5, domain.Horseshoe{}, false, ""},
		{"handkerchief", "handkerchief", 1.0, 0.5, domain.Handkerchief{}, false, ""},
		{"heart", "heart", 1.0, 0.5, domain.Heart{}, false, ""},
		{"disk", "disk", 1.0, 0.5, domain.Disk{}, false, ""},
		{"spiral", "spiral", 1.0, 0.5, domain.Spiral{}, false, ""},
		{"linear", "linear", 1.0, 0.5, domain.Linear{}, false, ""},
		{"unknown transform", "unknown", 1.0, 0.5, nil, true, "unknown transform: unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			transform, err := GetTransformByName(tc.transformName, tc.weight, tc.colorIndex)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, transform)
				assert.Equal(t, tc.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, transform)
				assert.IsType(t, tc.expectedType, transform)
				assert.Equal(t, tc.weight, transform.Weight())
				assert.Equal(t, tc.colorIndex, transform.ColorIndex())
			}
		})
	}
}

func TestGetPaletteByName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		paletteName  string
		expectedFunc func() [256]color.RGBA
	}{
		{"rainbow", "rainbow", domain.RainbowPalette},
		{"fire", "fire", domain.FirePalette},
		{"ocean", "ocean", domain.OceanPalette},
		{"default case", "unknown", domain.RainbowPalette},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			palette := GetPaletteByName(tc.paletteName)
			expectedPalette := tc.expectedFunc()
			assert.Equal(t, expectedPalette, palette)
		})
	}
}
