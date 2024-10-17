package comparators

import (
	"fmt"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type (
	NotImplementComparatorError struct {
		Field record.Field
		Cmp   where.ComparatorType
	}
	FailCastTypeError struct {
		Field        record.Field
		Cmp          where.ComparatorType
		ExpectedType string
		ReceivedType interface{}
	}
)

func (e NotImplementComparatorError) Error() string {
	return fmt.Sprintf("not implemented ComparatorType: %d, field = %s", e.Cmp, e.Field.String())
}

func (e NotImplementComparatorError) Is(err error) bool {
	_, ok := err.(NotImplementComparatorError)
	return ok
}

func NewNotImplementComparatorError(field record.Field, cmp where.ComparatorType) error {
	return NotImplementComparatorError{
		Field: field,
		Cmp:   cmp,
	}
}

func (e FailCastTypeError) Error() string {
	return fmt.Sprintf(
		"can't cast type: field = %s, ComparatorType = %d, value type = %T, expected type = %s",
		e.Field.String(),
		e.Cmp,
		e.ReceivedType,
		e.ExpectedType,
	)
}

func (e FailCastTypeError) Is(err error) bool {
	_, ok := err.(FailCastTypeError)
	return ok
}

func NewFailCastTypeError(field record.Field, cmp where.ComparatorType, receivedType interface{}, expectedType string) error {
	return FailCastTypeError{
		Field:        field,
		Cmp:          cmp,
		ExpectedType: expectedType,
		ReceivedType: receivedType,
	}
}
