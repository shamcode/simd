package comparators

import (
	"slices"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type ComparableFieldComparator[R record.Record, T record.LessComparable] struct {
	Cmp    where.ComparatorType
	Getter record.ComparableGetter[R, T]
	Value  []T
}

func (fc ComparableFieldComparator[R, T]) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc ComparableFieldComparator[R, T]) GetField() record.Field {
	return fc.Getter.Field
}

func (fc ComparableFieldComparator[R, T]) CompareValue(value T) (bool, error) {
	switch fc.Cmp { //nolint:exhaustive
	case where.EQ:
		return value == fc.Value[0], nil
	case where.GT:
		return value > fc.Value[0], nil
	case where.LT:
		return value < fc.Value[0], nil
	case where.GE:
		return value >= fc.Value[0], nil
	case where.LE:
		return value <= fc.Value[0], nil
	case where.InArray:
		return slices.Index(fc.Value, value) != -1, nil
	default:
		return false, NewNotImplementComparatorError(fc.GetField(), fc.Cmp)
	}
}

func (fc ComparableFieldComparator[R, T]) Compare(item R) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc ComparableFieldComparator[R, T]) ValuesCount() int {
	return len(fc.Value)
}

func (fc ComparableFieldComparator[R, T]) ValueAt(index int) interface{} {
	return fc.Value[index]
}
