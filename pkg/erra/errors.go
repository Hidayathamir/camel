package erra

import (
	"fmt"
)

// Wrap return wrapped error outer:: inner
func Wrap(inner, outer error) error {
	return fmt.Errorf("%w:: %w", outer, inner)
}

// Wrapf return wrapped error outer:: inner
func Wrapf(inner error, outer string, a ...any) error {
	e := fmt.Errorf(outer, a...)
	return fmt.Errorf("%w:: %w", e, inner)
}
