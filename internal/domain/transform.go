package domain

import (
	"math"
)

type Transform interface {
	Apply(p Point) Point
	Weight() float64
	ColorIndex() float64
}

// Composite applies a series of transformations in sequence.
type Composite struct {
	Transforms []Transform
}

// Apply iterates through and applies each transformation.
func (ct Composite) Apply(p Point) Point {
	for _, t := range ct.Transforms {
		p = t.Apply(p)
	}
	return p
}

// Weight returns the average weight of the contained transformations.
func (ct Composite) Weight() float64 {
	if len(ct.Transforms) == 0 {
		return 0
	}
	var totalWeight float64
	for _, t := range ct.Transforms {
		totalWeight += t.Weight()
	}
	return totalWeight / float64(len(ct.Transforms))
}

// ColorIndex returns the color index of the last transformation.
func (ct Composite) ColorIndex() float64 {
	if len(ct.Transforms) == 0 {
		return 0
	}
	// Return the color of the last transform in the chain
	return ct.Transforms[len(ct.Transforms)-1].ColorIndex()
}

type Affine struct {
	A, B, C, D, E, F, W, Color float64
}

func (a Affine) ColorIndex() float64 {
	return a.Color
}

func (a Affine) Weight() float64 {
	return a.W
}

func (a Affine) Apply(p Point) Point {
	return Point{
		X:          a.A*p.X + a.B*p.Y + a.C,
		Y:          a.D*p.X + a.E*p.Y + a.F,
		ColorIndex: p.ColorIndex,
	}
}

type Linear struct {
	W, Color float64
}

func (l Linear) ColorIndex() float64 {
	return l.Color
}

func (l Linear) Weight() float64 {
	return l.W
}

func (Linear) Apply(p Point) Point {
	return p
}

type Sinusoidal struct {
	W, Color float64
}

func (s Sinusoidal) ColorIndex() float64 {
	return s.Color
}

func (s Sinusoidal) Weight() float64 {
	return s.W
}

func (Sinusoidal) Apply(p Point) Point {
	return Point{
		X:          math.Sin(p.X),
		Y:          math.Sin(p.Y),
		ColorIndex: p.ColorIndex,
	}
}

type Spherical struct {
	W, Color float64
}

func (s Spherical) ColorIndex() float64 {
	return s.Color
}

func (s Spherical) Weight() float64 {
	return s.W
}

func (Spherical) Apply(p Point) Point {
	r2 := p.X*p.X + p.Y*p.Y
	if r2 == 0 {
		return p
	}
	inv := 1.0 / r2
	return Point{
		X:          p.X * inv,
		Y:          p.Y * inv,
		ColorIndex: p.ColorIndex,
	}
}

type Swirl struct {
	W, Color float64
}

func (s Swirl) ColorIndex() float64 {
	return s.Color
}

func (s Swirl) Weight() float64 {
	return s.W
}

func (Swirl) Apply(p Point) Point {
	r2 := p.X*p.X + p.Y*p.Y
	sinr := math.Sin(r2)
	cosr := math.Cos(r2)
	return Point{
		X:          p.X*sinr - p.Y*cosr,
		Y:          p.X*cosr + p.Y*sinr,
		ColorIndex: p.ColorIndex,
	}
}

type Horseshoe struct {
	W, Color float64
}

func (h Horseshoe) ColorIndex() float64 {
	return h.Color
}

func (h Horseshoe) Weight() float64 {
	return h.W
}

func (Horseshoe) Apply(p Point) Point {
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)
	if r == 0 {
		return p
	}
	return Point{
		X:          (p.X - p.Y) * (p.X + p.Y) / r,
		Y:          2 * p.X * p.Y / r,
		ColorIndex: p.ColorIndex,
	}
}

type Polar struct {
	W, Color float64
}

func (p Polar) ColorIndex() float64 {
	return p.Color
}

func (p Polar) Weight() float64 {
	return p.W
}

func (Polar) Apply(p Point) Point {
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)
	theta := math.Atan2(p.Y, p.X)
	return Point{
		X:          theta / math.Pi,
		Y:          r - 1.0,
		ColorIndex: p.ColorIndex,
	}
}

type Handkerchief struct {
	W, Color float64
}

func (h Handkerchief) ColorIndex() float64 {
	return h.Color
}

func (h Handkerchief) Weight() float64 {
	return h.W
}

func (Handkerchief) Apply(p Point) Point {
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)
	theta := math.Atan2(p.Y, p.X)
	return Point{
		X:          r * math.Sin(theta+r),
		Y:          r * math.Cos(theta-r),
		ColorIndex: p.ColorIndex,
	}
}

type Heart struct {
	W, Color float64
}

func (h Heart) ColorIndex() float64 {
	return h.Color
}

func (h Heart) Weight() float64 {
	return h.W
}

func (Heart) Apply(p Point) Point {
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)
	theta := math.Atan2(p.Y, p.X)
	return Point{
		X:          r * math.Sin(theta*r),
		Y:          -r * math.Cos(theta*r),
		ColorIndex: p.ColorIndex,
	}
}

type Spiral struct {
	W, Color float64
}

func (s Spiral) ColorIndex() float64 {
	return s.Color
}

func (s Spiral) Weight() float64 {
	return s.W
}

func (Spiral) Apply(p Point) Point {
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)
	theta := math.Atan2(p.Y, p.X)
	sin := math.Sin(theta)
	cos := math.Cos(theta)
	inv := 1.0 / (r + 1.0)
	return Point{
		X:          inv * (cos + sin),
		Y:          inv * (sin - cos),
		ColorIndex: p.ColorIndex,
	}
}

type Disk struct {
	W, Color float64
}

func (d Disk) ColorIndex() float64 {
	return d.Color
}

func (d Disk) Weight() float64 {
	return d.W
}

func (Disk) Apply(p Point) Point {
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)
	theta := math.Atan2(p.Y, p.X)
	s := theta / math.Pi
	return Point{
		X:          s * math.Sin(math.Pi*r),
		Y:          s * math.Cos(math.Pi*r),
		ColorIndex: p.ColorIndex,
	}
}
