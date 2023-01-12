package comparators

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type InterfaceFieldComparator struct {
	Cmp    where.ComparatorType
	Getter *record.InterfaceGetter
	Value  []interface{}
}

func (fc InterfaceFieldComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc InterfaceFieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc InterfaceFieldComparator) CompareValue(value interface{}) (bool, error) {
	switch fc.Cmp {
	case where.EQ:
		return value == fc.Value[0], nil
	case where.InArray:
		for _, x := range fc.Value {
			if x == value {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, NewErrNotImplementComparator(fc.GetField(), fc.Cmp)
	}
}

func (fc InterfaceFieldComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc InterfaceFieldComparator) ValuesCount() int {
	return len(fc.Value)
}

func (fc InterfaceFieldComparator) ValueAt(index int) interface{} {
	return fc.Value[index]
}
