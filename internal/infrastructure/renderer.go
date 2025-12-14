package infrastructure

import (
	"image"
	"image/color"
	"math"

	"hw4-fractal-flame/internal/domain"
)

type Renderer struct {
	Palette    [256]color.RGBA
	Gamma      float64
	Brightness float64
}

func NewRenderer(palette [256]color.RGBA, gamma, brightness float64) *Renderer {
	return &Renderer{
		Palette:    palette,
		Gamma:      gamma,
		Brightness: brightness,
	}
}

func (r *Renderer) Render(h *domain.Histogram) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, h.Width, h.Height))
	maxCount := h.MaxCount()

	if maxCount == 0 {
		return img
	}

	logMaxCount := math.Log10(float64(maxCount))
	invGamma := 1.0 / r.Gamma

	for y := range h.Height {
		for x := range h.Width {
			count := h.GetCount(x, y)

			if count == 0 {
				img.Set(x, y, color.RGBA{A: 255})
				continue
			}

			colorIndex := h.GetAverageColorIndex(x, y)
			baseColor := r.colorFromPalette(colorIndex)
			logDensity := math.Log10(float64(count)) / logMaxCount

			// gamma-correction
			alpha := math.Pow(logDensity, invGamma) * r.Brightness

			img.Set(x, y, color.RGBA{
				R: clamp(float64(baseColor.R) * alpha),
				G: clamp(float64(baseColor.G) * alpha),
				B: clamp(float64(baseColor.B) * alpha),
				A: 255,
			})
		}
	}

	return img
}

func (r *Renderer) colorFromPalette(index float64) color.RGBA {
	if index < 0 {
		index = 0
	}

	if index > 1 {
		index = 1
	}

	pos := int(index * float64(len(r.Palette)-1))

	return r.Palette[pos]
}

func clamp(v float64) uint8 {
	if v < 0 {
		return 0
	}

	if v > 255 {
		return 255
	}

	return uint8(v)
}
