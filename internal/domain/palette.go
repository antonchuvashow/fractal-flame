package domain

import (
	"image/color"
	"math"
)

func FirePalette() [256]color.RGBA {
	palette := [256]color.RGBA{}
	for i := range 256 {
		t := float64(i) / 255.0
		palette[i] = color.RGBA{
			R: uint8(min(255, int(t*3*255))),
			G: uint8(max(0, min(255, int((t-0.33)*3*255)))),
			B: uint8(max(0, min(255, int((t-0.66)*3*255)))),
			A: 255,
		}
	}
	return palette
}

func RainbowPalette() [256]color.RGBA {
	palette := [256]color.RGBA{}
	for i := range 256 {
		t := float64(i) / 255.0
		h := t * 6.0
		c := 1.0
		x := 1.0 - math.Abs(math.Mod(h, 2.0)-1.0)

		var r, g, b float64
		switch int(h) {
		case 0:
			r, g, b = c, x, 0
		case 1:
			r, g, b = x, c, 0
		case 2:
			r, g, b = 0, c, x
		case 3:
			r, g, b = 0, x, c
		case 4:
			r, g, b = x, 0, c
		default:
			r, g, b = c, 0, x
		}

		palette[i] = color.RGBA{
			R: uint8(r * 255),
			G: uint8(g * 255),
			B: uint8(b * 255),
			A: 255,
		}
	}
	return palette
}

func OceanPalette() [256]color.RGBA {
	var palette [256]color.RGBA

	for i := range 256 {
		t := float64(i) / 255.0

		var r, g, b float64

		if t < 0.25 {
			lt := t / 0.25
			r = lerp(5, 15, lt)
			g = lerp(10, 30, lt)
			b = lerp(40, 80, lt)
		} else if t < 0.5 {
			lt := (t - 0.25) / 0.25
			r = lerp(15, 30, lt)
			g = lerp(30, 100, lt)
			b = lerp(80, 170, lt)
		} else if t < 0.75 {
			lt := (t - 0.5) / 0.25
			r = lerp(30, 70, lt)
			g = lerp(100, 180, lt)
			b = lerp(170, 200, lt)
		} else {
			lt := (t - 0.75) / 0.25
			r = lerp(70, 180, lt)
			g = lerp(180, 230, lt)
			b = lerp(200, 255, lt)
		}

		palette[i] = color.RGBA{
			R: clampByte(r),
			G: clampByte(g),
			B: clampByte(b),
			A: 255,
		}
	}

	return palette
}

// lerp - linear interpolation
func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func clampByte(v float64) uint8 {
	if v < 0 {
		return 0
	}

	if v > 255 {
		return 255
	}

	return uint8(v)
}
