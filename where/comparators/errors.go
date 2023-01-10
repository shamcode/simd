package comparators

import (
	"fmt"
	"github.com/shamcode/simd/where"
)

type (
	ErrNotImplementComparator struct {
		Field string
		Cmp   where.ComparatorType
	}
	ErrFailCastType struct {
		Field        string
		Cmp          where.ComparatorType
		ExpectedType string
		ReceivedType interface{}
	}
)

func (e ErrNotImplementComparator) Error() string {
	return fmt.Sprintf("not implemented ComparatorType: %d, field = %s", e.Cmp, e.Field)
}

func (e ErrNotImplementComparator) Is(err error) bool {
	_, ok := err.(ErrNotImplementComparator)
	return ok
}

func NewErrNotImplementComparator(field string, cmp where.ComparatorType) error {
	return ErrNotImplementComparator{
		Field: field,
		Cmp:   cmp,
	}
}

func (e ErrFailCastType) Error() string {
	return fmt.Sprintf(
		"can't cast type: field = %s, ComparatorType = %d, value type = %T, expected type = %s",
		e.Field,
		e.Cmp,
		e.ReceivedType,
		e.ExpectedType,
	)
}

func (e ErrFailCastType) Is(err error) bool {
	_, ok := err.(ErrFailCastType)
	return ok
}

func NewErrFailCastType(field string, cmp where.ComparatorType, receivedType interface{}, expectedType string) error {
	return ErrFailCastType{
		Field:        field,
		Cmp:          cmp,
		ExpectedType: expectedType,
		ReceivedType: receivedType,
	}
}
