package comparators

import (
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type (
	ErrNotImplementComparator struct {
		Field record.Field
		Cmp   where.ComparatorType
	}
	ErrFailCastType struct {
		Field        record.Field
		Cmp          where.ComparatorType
		ExpectedType string
		ReceivedType interface{}
	}
)

func (e ErrNotImplementComparator) Error() string {
	return fmt.Sprintf("not implemented ComparatorType: %d, field = %s", e.Cmp, e.Field.String())
}

func (e ErrNotImplementComparator) Is(err error) bool {
	_, ok := err.(ErrNotImplementComparator)
	return ok
}

func NewErrNotImplementComparator(field record.Field, cmp where.ComparatorType) error {
	return ErrNotImplementComparator{
		Field: field,
		Cmp:   cmp,
	}
}

func (e ErrFailCastType) Error() string {
	return fmt.Sprintf(
		"can't cast type: field = %s, ComparatorType = %d, value type = %T, expected type = %s",
		e.Field.String(),
		e.Cmp,
		e.ReceivedType,
		e.ExpectedType,
	)
}

func (e ErrFailCastType) Is(err error) bool {
	_, ok := err.(ErrFailCastType)
	return ok
}

func NewErrFailCastType(field record.Field, cmp where.ComparatorType, receivedType interface{}, expectedType string) error {
	return ErrFailCastType{
		Field:        field,
		Cmp:          cmp,
		ExpectedType: expectedType,
		ReceivedType: receivedType,
	}
}
