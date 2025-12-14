package cli

import "fmt"

type ErrInvalidWeight struct {
	Part string
}

func (e *ErrInvalidWeight) Error() string {
	return fmt.Sprintf("invalid weight in %q", e.Part)
}

type ErrInvalidColor struct {
	Part string
}

func (e *ErrInvalidColor) Error() string {
	return fmt.Sprintf("invalid color in %q", e.Part)
}

type ErrInvalidAffine struct {
	Part string
}

func (e *ErrInvalidAffine) Error() string {
	return fmt.Sprintf("invalid affine in %q", e.Part)
}

type ErrInvalidCoefficient struct {
	Coefficient string
	Part        string
}

func (e *ErrInvalidCoefficient) Error() string {
	return fmt.Sprintf("invalid coefficient %q in %q", e.Coefficient, e.Part)
}
