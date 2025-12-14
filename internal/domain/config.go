package domain

import (
	"image/color"
)

// Config is a pure domain entity.
type Config struct {
	Size            Size
	Iterations      int
	Output          string
	Threads         int
	Seed            float64
	Transforms      []TransformEntry
	SymmetryLevel   int
	Gamma           float64
	Brightness      float64
	Palette         Palette
	AffineTransform []AffineEntry
}

// Size describes the size of the output image.
type Size struct {
	Width  int
	Height int
}

// TransformEntry contains information about a transformation.
type TransformEntry struct {
	Transform  Transform
	Name       string
	Weight     float64
	ColorIndex float64
}

// AffineEntry contains information about an affine transformation.
type AffineEntry struct {
	Transform Transform
	A, B, C   float64
	D, E, F   float64
}

// Palette contains information about a color palette.
type Palette struct {
	Name   string
	Colors [256]color.RGBA
}
