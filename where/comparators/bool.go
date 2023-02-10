package comparators

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type BoolFieldComparator struct {
	Cmp    where.ComparatorType
	Getter record.BoolGetter
	Value  []bool
}

func (fc BoolFieldComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc BoolFieldComparator) GetField() record.Field {
	return fc.Getter.Field
}

func (fc BoolFieldComparator) CompareValue(value bool) (bool, error) {
	switch fc.Cmp {
	case where.EQ:
		return value == fc.Value[0], nil
	default:
		return false, NewNotImplementComparatorError(fc.GetField(), fc.Cmp)
	}
}

func (fc BoolFieldComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc BoolFieldComparator) ValuesCount() int {
	return len(fc.Value)
}

func (fc BoolFieldComparator) ValueAt(index int) interface{} {
	return fc.Value[index]
}
