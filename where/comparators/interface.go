package comparators

import (
	"slices"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type InterfaceFieldComparator[R record.Record] struct {
	Cmp    where.ComparatorType
	Getter record.InterfaceGetter[R]
	Value  []any
}

func (fc InterfaceFieldComparator[R]) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc InterfaceFieldComparator[R]) GetField() record.Field {
	return fc.Getter.Field
}

func (fc InterfaceFieldComparator[R]) CompareValue(value any) (bool, error) {
	switch fc.Cmp { //nolint:exhaustive
	case where.EQ:
		return value == fc.Value[0], nil
	case where.InArray:
		return slices.Index(fc.Value, value) != -1, nil
	default:
		return false, NewNotImplementComparatorError(fc.GetField(), fc.Cmp)
	}
}

func (fc InterfaceFieldComparator[R]) Compare(item R) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc InterfaceFieldComparator[R]) ValuesCount() int {
	return len(fc.Value)
}

func (fc InterfaceFieldComparator[R]) ValueAt(index int) interface{} {
	return fc.Value[index]
}
