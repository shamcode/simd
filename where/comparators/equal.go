package comparators

import (
	"slices"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type EqualComparator[R record.Record, T comparable] struct {
	Cmp    where.ComparatorType
	Getter record.GetterInterface[R, T]
	Value  []T
}

func (fc EqualComparator[R, T]) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc EqualComparator[R, T]) GetField() record.Field {
	return fc.Getter
}

func (fc EqualComparator[R, T]) CompareValue(value T) (bool, error) {
	switch fc.Cmp { //nolint:exhaustive
	case where.EQ:
		return value == fc.Value[0], nil
	case where.InArray:
		return slices.Index(fc.Value, value) >= 0, nil
	default:
		return false, NewNotImplementComparatorError(fc.GetField(), fc.Cmp)
	}
}

func (fc EqualComparator[R, T]) Compare(item R) (bool, error) {
	return fc.CompareValue(fc.Getter.GetForRecord(item))
}

func (fc EqualComparator[R, T]) ValuesCount() int {
	return len(fc.Value)
}

func (fc EqualComparator[R, T]) ValueAt(index int) any {
	return fc.Value[index]
}

func NewEqualComparator[R record.Record, T comparable](
	cmp where.ComparatorType,
	getter record.GetterInterface[R, T],
	value ...T,
) EqualComparator[R, T] {
	return EqualComparator[R, T]{
		Cmp:    cmp,
		Getter: getter,
		Value:  value,
	}
}
