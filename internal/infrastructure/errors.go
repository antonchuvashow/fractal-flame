package infrastructure

import "fmt"

type ErrUnknownTransform struct {
	transform string
}

func (e *ErrUnknownTransform) Error() string {
	return fmt.Sprintf("unknown transform: %s", e.transform)
}
