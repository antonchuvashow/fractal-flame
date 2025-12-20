package application

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"hw4-fractal-flame/internal/domain"
)

// mockConfig creates a default configuration for testing purposes.
func mockConfig() *domain.Config {
	return &domain.Config{
		Transforms: []domain.TransformEntry{
			{Transform: domain.Sinusoidal{W: 1, Color: 0.5}},
		},
		AffineTransform: []domain.Affine{
			{A: 1, B: 0, C: 0, D: 0, E: 1, F: 0, W: 1, Color: 0.2},
		},
		Palette:       domain.Palette{Colors: domain.FirePalette()},
		SymmetryLevel: 1,
		Size:          domain.Size{Width: 10, Height: 10},
		Iterations:    100,
		Threads:       1,
		Seed:          0.123,
		Gamma:         2.2,
		Brightness:    1.0,
		Output:        "test.png",
	}
}

func TestNewApplication(t *testing.T) {
	t.Run("should create a new application with the given config", func(t *testing.T) {
		config := mockConfig()
		app := NewApplication(config)

		assert.NotNil(t, app)
		assert.Equal(t, config, app.Config)
	})
}

func TestApplication_Run(t *testing.T) {
	t.Run("should successfully generate an image file on the happy path", func(t *testing.T) {
		// Setup: Create a temporary directory for the output file
		tempDir := t.TempDir()
		config := mockConfig()
		config.Output = filepath.Join(tempDir, "output.png")

		app := NewApplication(config)

		// Execute
		err := app.Run()

		// Assert
		assert.NoError(t, err)
		_, err = os.Stat(config.Output)
		assert.False(t, os.IsNotExist(err), "Expected output file to be created")
	})

	t.Run("should return an error for an invalid output path", func(t *testing.T) {
		// Setup: an invalid path
		config := mockConfig()
		config.Output = "/non_existent_dir/output.png"

		app := NewApplication(config)

		// Execute
		err := app.Run()

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no such file or directory")
	})
}

func TestMergeHistograms(t *testing.T) {
	testCases := []struct {
		name          string
		histograms    []*domain.Histogram
		expectedNil   bool
		expectedTotal uint32
		expectedColor float64
	}{
		{
			name:        "should return nil when merging an empty list",
			histograms:  []*domain.Histogram{},
			expectedNil: true,
		},
		{
			name: "should return the same histogram if only one is provided",
			histograms: []*domain.Histogram{
				{Width: 1, Height: 1, Count: [][]uint32{{10}}, ColorSum: [][]float64{{5.0}}},
			},
			expectedTotal: 10,
			expectedColor: 5.0,
		},
		{
			name: "should correctly merge multiple histograms",
			histograms: []*domain.Histogram{
				{Width: 1, Height: 1, Count: [][]uint32{{10}}, ColorSum: [][]float64{{5.0}}},
				{Width: 1, Height: 1, Count: [][]uint32{{20}}, ColorSum: [][]float64{{10.0}}},
				{Width: 1, Height: 1, Count: [][]uint32{{30}}, ColorSum: [][]float64{{15.0}}},
			},
			expectedTotal: 60,
			expectedColor: 30.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MergeHistograms(tc.histograms)

			if tc.expectedNil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)
			assert.Equal(t, tc.expectedTotal, result.Count[0][0])
			assert.InDelta(t, tc.expectedColor, result.ColorSum[0][0], 1e-9)
		})
	}
}
