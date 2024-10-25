package comparators

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type SetFieldComparator[R record.Record] struct {
	Cmp    where.ComparatorType
	Getter record.SetGetter[R]
	Value  []interface{}
}

func (fc SetFieldComparator[R]) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc SetFieldComparator[R]) GetField() record.Field {
	return fc.Getter.Field
}

func (fc SetFieldComparator[R]) CompareValue(value record.Set) (bool, error) {
	switch fc.Cmp { //nolint:exhaustive
	case where.SetHas:
		return value.Has(fc.Value[0]), nil
	default:
		return false, NewNotImplementComparatorError(fc.GetField(), fc.Cmp)
	}
}

func (fc SetFieldComparator[R]) Compare(item R) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc SetFieldComparator[R]) ValuesCount() int {
	return len(fc.Value)
}

func (fc SetFieldComparator[R]) ValueAt(index int) interface{} {
	return fc.Value[index]
}
