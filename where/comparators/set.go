package comparators

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type SetFieldComparator[R record.Record, T comparable] struct {
	Cmp    where.ComparatorType
	Getter record.SetGetter[R, T]
	Value  []T
}

func (fc SetFieldComparator[R, T]) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc SetFieldComparator[R, T]) GetField() record.Field {
	return fc.Getter.Field
}

func (fc SetFieldComparator[R, T]) CompareValue(value record.Set[T]) (bool, error) {
	switch fc.Cmp { //nolint:exhaustive
	case where.SetHas:
		return value.Has(fc.Value[0]), nil
	default:
		return false, NewNotImplementComparatorError(fc.GetField(), fc.Cmp)
	}
}

func (fc SetFieldComparator[R, T]) Compare(item R) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc SetFieldComparator[R, T]) ValuesCount() int {
	return len(fc.Value)
}

func (fc SetFieldComparator[R, T]) ValueAt(index int) any {
	return fc.Value[index]
}

func NewSetFieldComparator[R record.Record, T comparable](
	cmp where.ComparatorType,
	getter record.SetGetter[R, T],
	value ...T,
) SetFieldComparator[R, T] {
	return SetFieldComparator[R, T]{
		Cmp:    cmp,
		Getter: getter,
		Value:  value,
	}
}
