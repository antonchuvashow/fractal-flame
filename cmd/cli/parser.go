package cli

import (
	"math/rand"
	"strconv"
	"strings"

	"hw4-fractal-flame/internal/domain"
	"hw4-fractal-flame/internal/infrastructure"
)

// ParseTransforms parse the string of the form "swirl:1.0,horseshoe:0.8,sinusoidal"
// "name:weight:color" or "name:weight" or "name".
func ParseTransforms(input string, rnd *rand.Rand) ([]domain.TransformEntry, error) {
	if input == "" {
		return nil, nil
	}

	parts := strings.Split(input, ",")
	transforms := make([]domain.TransformEntry, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		tc, err := parseTransformPart(part, rnd)
		if err != nil {
			return nil, err
		}

		transforms = append(transforms, tc)
	}

	return transforms, nil
}

// parseTransformPart parse one line transform: "swirl:1.0:0.5" or "swirl:1.0" or "swirl".
func parseTransformPart(part string, rnd *rand.Rand) (domain.TransformEntry, error) {
	tokens := strings.Split(part, ":")

	name := strings.TrimSpace(tokens[0])
	weight := 1.0               // default weight
	colorIndex := rnd.Float64() // random color by default

	var err error

	if len(tokens) >= 2 && tokens[1] != "" {
		weight, err = strconv.ParseFloat(tokens[1], 64)
		if err != nil {
			return domain.TransformEntry{}, &ErrInvalidWeight{Part: part}
		}
	}

	if len(tokens) >= 3 && tokens[2] != "" {
		colorIndex, err = strconv.ParseFloat(tokens[2], 64)
		if err != nil {
			return domain.TransformEntry{}, &ErrInvalidColor{Part: part}
		}
	}

	t, err := infrastructure.GetTransformByName(name, weight, colorIndex)
	if err != nil {
		return domain.TransformEntry{}, err
	}

	return domain.TransformEntry{
		Transform:  t,
		Name:       name,
		Weight:     weight,
		ColorIndex: colorIndex,
	}, nil
}

// ParseAffineParams parse the affine transforation string
// "a,b,c,d,e,f:weight:color;a,b,c,d,e,f:weight:color" or "a,b,c,d,e,f".
func ParseAffineParams(input string, rnd *rand.Rand) ([]domain.AffineEntry, error) {
	if input == "" {
		return nil, nil
	}

	parts := strings.Split(input, ";")
	affines := make([]domain.AffineEntry, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		ac, err := parseAffinePart(part, rnd)
		if err != nil {
			return nil, err
		}

		affines = append(affines, ac)
	}

	return affines, nil
}

// parseAffinePart parse one line of affine transform "0.5,0.1,0,0,0.5,0:1.0:0.3".
func parseAffinePart(part string, rnd *rand.Rand) (domain.AffineEntry, error) {
	const numberOfCoeffs = 6

	sections := strings.Split(part, ":")

	coeffs := strings.Split(sections[0], ",")
	if len(coeffs) != numberOfCoeffs {
		return domain.AffineEntry{}, &ErrInvalidAffine{Part: part}
	}

	values := make([]float64, numberOfCoeffs)
	for i, c := range coeffs {
		v, err := strconv.ParseFloat(strings.TrimSpace(c), 64)
		if err != nil {
			return domain.AffineEntry{}, &ErrInvalidCoefficient{Part: part, Coefficient: c}
		}

		values[i] = v
	}

	weight := 1.0
	colorIndex := rnd.Float64()

	var err error

	if len(sections) >= 2 && sections[1] != "" {
		weight, err = strconv.ParseFloat(sections[1], 64)
		if err != nil {
			return domain.AffineEntry{}, &ErrInvalidWeight{Part: part}
		}
	}

	if len(sections) >= 3 && sections[2] != "" {
		colorIndex, err = strconv.ParseFloat(sections[2], 64)
		if err != nil {
			return domain.AffineEntry{}, &ErrInvalidColor{Part: part}
		}
	}

	ac := domain.AffineEntry{
		A: values[0],
		B: values[1],
		C: values[2],
		D: values[3],
		E: values[4],
		F: values[5],
		Transform: domain.Affine{
			A:     values[0],
			B:     values[1],
			C:     values[2],
			D:     values[3],
			E:     values[4],
			F:     values[5],
			W:     weight,
			Color: colorIndex,
		},
	}

	return ac, nil
}
