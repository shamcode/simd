package comparators

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type BoolFieldComparator[R record.Record] struct {
	Cmp    where.ComparatorType
	Getter record.BoolGetter[R]
	Value  []bool
}

func (fc BoolFieldComparator[R]) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc BoolFieldComparator[R]) GetField() record.Field {
	return fc.Getter.Field
}

func (fc BoolFieldComparator[R]) CompareValue(value bool) (bool, error) {
	switch fc.Cmp { //nolint:exhaustive
	case where.EQ:
		return value == fc.Value[0], nil
	default:
		return false, NewNotImplementComparatorError(fc.GetField(), fc.Cmp)
	}
}

func (fc BoolFieldComparator[R]) Compare(item R) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc BoolFieldComparator[R]) ValuesCount() int {
	return len(fc.Value)
}

func (fc BoolFieldComparator[R]) ValueAt(index int) interface{} {
	return fc.Value[index]
}
