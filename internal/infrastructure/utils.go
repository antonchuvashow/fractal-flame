package infrastructure

import (
	"image/color"

	"hw4-fractal-flame/internal/domain"
)

func GetTransformByName(name string, weight, colorIndex float64) (domain.Transform, error) {
	switch name {
	case "sinusoidal":
		return domain.Sinusoidal{W: weight, Color: colorIndex}, nil
	case "spherical":
		return domain.Spherical{W: weight, Color: colorIndex}, nil
	case "swirl":
		return domain.Swirl{W: weight, Color: colorIndex}, nil
	case "horseshoe":
		return domain.Horseshoe{W: weight, Color: colorIndex}, nil
	case "handkerchief":
		return domain.Handkerchief{W: weight, Color: colorIndex}, nil
	case "heart":
		return domain.Heart{W: weight, Color: colorIndex}, nil
	case "disk":
		return domain.Disk{W: weight, Color: colorIndex}, nil
	case "spiral":
		return domain.Spiral{W: weight, Color: colorIndex}, nil
	case "linear":
		return domain.Linear{W: weight, Color: colorIndex}, nil
	default:
		return nil, &ErrUnknownTransform{transform: name}
	}
}

func GetPaletteByName(name string) [256]color.RGBA {
	switch name {
	case "rainbow":
		return domain.RainbowPalette()
	case "fire":
		return domain.FirePalette()
	case "ocean":
		return domain.OceanPalette()
	default:
		return domain.RainbowPalette()
	}
}
