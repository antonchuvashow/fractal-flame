package domain

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposite_Apply(t *testing.T) {
	p := Point{X: 1, Y: 1, ColorIndex: 0}
	transforms := []Transform{
		Affine{A: 1, B: 0, C: 1, D: 0, E: 1, F: 1, W: 1, Color: 0.5}, // p becomes (2, 2)
		Sinusoidal{W: 1, Color: 0.5},                                 // p becomes (sin(2), sin(2))
	}
	ct := Composite{Transforms: transforms}
	expected := Point{X: math.Sin(2), Y: math.Sin(2), ColorIndex: 0}

	result := ct.Apply(p)

	assert.InDelta(t, expected.X, result.X, 1e-9)
	assert.InDelta(t, expected.Y, result.Y, 1e-9)
}

func TestComposite_Weight(t *testing.T) {
	testCases := []struct {
		name       string
		transforms []Transform
		expected   float64
	}{
		{"no transforms", []Transform{}, 0},
		{"single transform", []Transform{Affine{W: 5}}, 5},
		{"multiple transforms", []Transform{Affine{W: 2}, Sinusoidal{W: 4}}, 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ct := Composite{Transforms: tc.transforms}
			assert.Equal(t, tc.expected, ct.Weight())
		})
	}
}

func TestComposite_ColorIndex(t *testing.T) {
	testCases := []struct {
		name       string
		transforms []Transform
		expected   float64
	}{
		{"no transforms", []Transform{}, 0},
		{"single transform", []Transform{Affine{Color: 0.7}}, 0.7},
		{"multiple transforms", []Transform{Affine{Color: 0.2}, Sinusoidal{Color: 0.8}}, 0.8},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ct := Composite{Transforms: tc.transforms}
			assert.Equal(t, tc.expected, ct.ColorIndex())
		})
	}
}

func TestTransformations(t *testing.T) {
	p := Point{X: 0.5, Y: -0.5, ColorIndex: 0.3}

	testCases := []struct {
		name      string
		transform Transform
		expected  Point
	}{
		{
			name:      "Affine",
			transform: Affine{A: 1, B: 2, C: 3, D: 4, E: 5, F: 6, W: 1, Color: 0.5},
			expected:  Point{X: 1*p.X + 2*p.Y + 3, Y: 4*p.X + 5*p.Y + 6, ColorIndex: p.ColorIndex},
		},
		{
			name:      "Linear",
			transform: Linear{W: 1, Color: 0.5},
			expected:  p,
		},
		{
			name:      "Sinusoidal",
			transform: Sinusoidal{W: 1, Color: 0.5},
			expected:  Point{X: math.Sin(p.X), Y: math.Sin(p.Y), ColorIndex: p.ColorIndex},
		},
		{
			name:      "Spherical",
			transform: Spherical{W: 1, Color: 0.5},
			expected:  Point{X: p.X / (p.X*p.X + p.Y*p.Y), Y: p.Y / (p.X*p.X + p.Y*p.Y), ColorIndex: p.ColorIndex},
		},
		{
			name:      "Swirl",
			transform: Swirl{W: 1, Color: 0.5},
			expected: Point{
				X:          p.X*math.Sin(p.X*p.X+p.Y*p.Y) - p.Y*math.Cos(p.X*p.X+p.Y*p.Y),
				Y:          p.X*math.Cos(p.X*p.X+p.Y*p.Y) + p.Y*math.Sin(p.X*p.X+p.Y*p.Y),
				ColorIndex: p.ColorIndex,
			},
		},
		{
			name:      "Horseshoe",
			transform: Horseshoe{W: 1, Color: 0.5},
			expected: Point{
				X:          (p.X - p.Y) * (p.X + p.Y) / math.Sqrt(p.X*p.X+p.Y*p.Y),
				Y:          2 * p.X * p.Y / math.Sqrt(p.X*p.X+p.Y*p.Y),
				ColorIndex: p.ColorIndex,
			},
		},
		{
			name:      "Polar",
			transform: Polar{W: 1, Color: 0.5},
			expected: Point{
				X:          math.Atan2(p.Y, p.X) / math.Pi,
				Y:          math.Sqrt(p.X*p.X+p.Y*p.Y) - 1.0,
				ColorIndex: p.ColorIndex,
			},
		},
		{
			name:      "Handkerchief",
			transform: Handkerchief{W: 1, Color: 0.5},
			expected: Point{
				X:          math.Sqrt(p.X*p.X+p.Y*p.Y) * math.Sin(math.Atan2(p.Y, p.X)+math.Sqrt(p.X*p.X+p.Y*p.Y)),
				Y:          math.Sqrt(p.X*p.X+p.Y*p.Y) * math.Cos(math.Atan2(p.Y, p.X)-math.Sqrt(p.X*p.X+p.Y*p.Y)),
				ColorIndex: p.ColorIndex,
			},
		},
		{
			name:      "Heart",
			transform: Heart{W: 1, Color: 0.5},
			expected: Point{
				X:          math.Sqrt(p.X*p.X+p.Y*p.Y) * math.Sin(math.Atan2(p.Y, p.X)*math.Sqrt(p.X*p.X+p.Y*p.Y)),
				Y:          -math.Sqrt(p.X*p.X+p.Y*p.Y) * math.Cos(math.Atan2(p.Y, p.X)*math.Sqrt(p.X*p.X+p.Y*p.Y)),
				ColorIndex: p.ColorIndex,
			},
		},
		{
			name:      "Spiral",
			transform: Spiral{W: 1, Color: 0.5},
			expected: Point{
				X:          (1.0 / (math.Sqrt(p.X*p.X+p.Y*p.Y) + 1.0)) * (math.Cos(math.Atan2(p.Y, p.X)) + math.Sin(math.Atan2(p.Y, p.X))),
				Y:          (1.0 / (math.Sqrt(p.X*p.X+p.Y*p.Y) + 1.0)) * (math.Sin(math.Atan2(p.Y, p.X)) - math.Cos(math.Atan2(p.Y, p.X))),
				ColorIndex: p.ColorIndex,
			},
		},
		{
			name:      "Disk",
			transform: Disk{W: 1, Color: 0.5},
			expected: Point{
				X:          (math.Atan2(p.Y, p.X) / math.Pi) * math.Sin(math.Pi*math.Sqrt(p.X*p.X+p.Y*p.Y)),
				Y:          (math.Atan2(p.Y, p.X) / math.Pi) * math.Cos(math.Pi*math.Sqrt(p.X*p.X+p.Y*p.Y)),
				ColorIndex: p.ColorIndex,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name+"_Apply", func(t *testing.T) {
			result := tc.transform.Apply(p)
			assert.InDelta(t, tc.expected.X, result.X, 1e-9)
			assert.InDelta(t, tc.expected.Y, result.Y, 1e-9)
			assert.Equal(t, tc.expected.ColorIndex, result.ColorIndex)
		})

		t.Run(tc.name+"_Weight", func(t *testing.T) {
			assert.Equal(t, 1.0, tc.transform.Weight())
		})

		t.Run(tc.name+"_ColorIndex", func(t *testing.T) {
			assert.Equal(t, 0.5, tc.transform.ColorIndex())
		})
	}
}
