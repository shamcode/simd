package comparators

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type ComparableFieldComparator[R record.Record, T record.LessComparable] struct {
	EqualComparator[R, T]
}

func (fc ComparableFieldComparator[R, T]) CompareValue(value T) (bool, error) {
	switch fc.Cmp { //nolint:exhaustive
	case where.GT:
		return value > fc.Value[0], nil
	case where.LT:
		return value < fc.Value[0], nil
	case where.GE:
		return value >= fc.Value[0], nil
	case where.LE:
		return value <= fc.Value[0], nil
	default:
		return fc.EqualComparator.CompareValue(value)
	}
}

func (fc ComparableFieldComparator[R, T]) Compare(item R) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}
