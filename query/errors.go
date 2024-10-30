package query

import (
	"errors"
	"fmt"

	"github.com/shamcode/simd/record"
)

var (
	ErrOrBeforeAnyConditions   = errors.New(".Or() before any condition not supported, add any condition before .Or()")
	ErrNotOpenBracket          = errors.New(".Not().OpenBracket() not supported")
	ErrCloseBracketWithoutOpen = errors.New("close bracket without open")
	ErrInvalidBracketBalance   = errors.New("invalid bracket balance: has not closed bracket")
)

type (
	GetterError struct {
		Field record.Field
		Err   error
	}
	CastError[A, B any] struct {
		Expected A
		Actual   B
	}
)

func (e GetterError) Error() string {
	return e.Field.String() + ": " + e.Err.Error()
}

func (e GetterError) Unwrap() error {
	return e.Err
}

func (e CastError[A, B]) Error() string {
	return fmt.Sprintf("cannot cast %T to %T", e.Actual, e.Expected)
}
