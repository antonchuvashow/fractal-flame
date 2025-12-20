// internal/infrastructure/config/loader.go
package config

import (
	"encoding/json"
	"fmt"
	"os"

	"hw4-fractal-flame/internal/domain"
	"hw4-fractal-flame/internal/infrastructure"
)

// Load gets JSON file and creates domain.Config .
func Load(path string) (*domain.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var dto DTO

	err = json.Unmarshal(data, &dto)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return MapToDomain(&dto)
}

// MapToDomain converts DTO into Domain model.
func MapToDomain(dto *DTO) (*domain.Config, error) {
	transforms, err := mapTransforms(dto.Transforms)
	if err != nil {
		return nil, fmt.Errorf("failed to map transforms: %w", err)
	}

	affines := mapAffines(dto.Affine)

	return &domain.Config{
		Size: domain.Size{
			Width:  dto.Size.Width,
			Height: dto.Size.Height,
		},
		Iterations:      dto.Iterations,
		Output:          dto.Output,
		Threads:         dto.Threads,
		Seed:            dto.Seed,
		Transforms:      transforms,
		SymmetryLevel:   dto.Symmetry,
		GammaCorrection: dto.GammaCorrection,
		Gamma:           dto.Gamma,
		Brightness:      dto.Brightness,
		Palette:         mapPalette(dto.Palette),
		AffineTransform: affines,
	}, nil
}

func mapTransforms(dtos []TransformDTO) ([]domain.TransformEntry, error) {
	result := make([]domain.TransformEntry, 0, len(dtos))

	for _, dto := range dtos {
		transform, err := infrastructure.GetTransformByName(dto.Name, dto.Weight, dto.ColorIndex)
		if err != nil {
			return nil, fmt.Errorf("unknown transform %q: %w", dto.Name, err)
		}

		result = append(result, domain.TransformEntry{
			Transform:  transform,
			Name:       dto.Name,
			Weight:     dto.Weight,
			ColorIndex: dto.ColorIndex,
		})
	}

	return result, nil
}

func mapAffines(dtos []AffineDTO) []domain.Affine {
	result := make([]domain.Affine, 0, len(dtos))

	for _, dto := range dtos {
		entry := domain.Affine{
			A: dto.A, B: dto.B, C: dto.C,
			D: dto.D, E: dto.E, F: dto.F,
		}

		result = append(result, entry)
	}

	return result
}

func mapPalette(name string) domain.Palette {
	return domain.Palette{
		Name:   name,
		Colors: infrastructure.GetPaletteByName(name),
	}
}
