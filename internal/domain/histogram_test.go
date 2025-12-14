package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHistogram(t *testing.T) {
	t.Parallel()

	t.Run("creates a histogram with correct dimensions", func(t *testing.T) {
		t.Parallel()

		width, height := 100, 200
		hist := NewHistogram(width, height)

		assert.Equal(t, width, hist.Width)
		assert.Equal(t, height, hist.Height)
		assert.Len(t, hist.Count, height)
		assert.Len(t, hist.ColorSum, height)
		assert.Len(t, hist.Count[0], width)
		assert.Len(t, hist.ColorSum[0], width)
	})

	t.Run("initializes with zero counts and colors", func(t *testing.T) {
		t.Parallel()

		hist := NewHistogram(10, 10)
		for y := range hist.Height {
			for x := range hist.Width {
				assert.Equal(t, uint32(0), hist.Count[y][x])
				assert.InDelta(t, 0.0, hist.ColorSum[y][x], 1e-9)
			}
		}
	})
}

func TestHistogram_Contains(t *testing.T) {
	t.Parallel()

	hist := NewHistogram(50, 50)

	testCases := []struct {
		name     string
		x, y     int
		expected bool
	}{
		{"inside bounds", 25, 25, true},
		{"edge case bottom-left", 0, 0, true},
		{"edge case top-right", 49, 49, true},
		{"outside x-axis", 50, 25, false},
		{"outside y-axis", 25, 50, false},
		{"negative coordinates", -1, -1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, hist.Contains(tc.x, tc.y))
		})
	}
}

func TestHistogram_Increment(t *testing.T) {
	hist := NewHistogram(10, 10)
	x, y := 5, 5
	colorIndex := 0.75

	hist.Increment(x, y, colorIndex)
	assert.Equal(t, uint32(1), hist.Count[y][x])
	assert.InDelta(t, colorIndex, hist.ColorSum[y][x], 1e-9)

	hist.Increment(x, y, colorIndex)
	assert.Equal(t, uint32(2), hist.Count[y][x])
	assert.InDelta(t, colorIndex*2, hist.ColorSum[y][x], 1e-9)
}

func TestHistogram_GetCount(t *testing.T) {
	hist := NewHistogram(10, 10)
	x, y := 3, 4
	hist.Count[y][x] = 5

	assert.Equal(t, uint32(5), hist.GetCount(x, y))
	assert.Equal(t, uint32(0), hist.GetCount(0, 0))
}

func TestHistogram_GetAverageColorIndex(t *testing.T) {
	t.Parallel()

	hist := NewHistogram(10, 10)
	x, y := 2, 2

	t.Run("returns zero for no hits", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, float64(0), hist.GetAverageColorIndex(x, y))
	})

	t.Run("calculates average color correctly", func(t *testing.T) {
		t.Parallel()

		hist.Increment(x, y, 0.5)
		hist.Increment(x, y, 1.0)
		// (0.5 + 1.0) / 2 = 0.75
		assert.InDelta(t, 0.75, hist.GetAverageColorIndex(x, y), 1e-9)
	})
}

func TestHistogram_MaxCount(t *testing.T) {
	t.Parallel()

	hist := NewHistogram(5, 5)

	t.Run("returns zero for an empty histogram", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, uint32(0), hist.MaxCount())
	})

	t.Run("finds the maximum count", func(t *testing.T) {
		t.Parallel()

		hist.Count[1][1] = 10
		hist.Count[2][3] = 25
		hist.Count[4][0] = 15
		assert.Equal(t, uint32(25), hist.MaxCount())
	})
}
