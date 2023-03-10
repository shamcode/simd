package comparators

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type Enum8FieldComparator struct {
	Cmp    where.ComparatorType
	Getter record.Enum8Getter
	Value  []record.Enum8
}

func (fc Enum8FieldComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc Enum8FieldComparator) GetField() record.Field {
	return fc.Getter.Field
}

func (fc Enum8FieldComparator) CompareValue(value uint8) (bool, error) {
	switch fc.Cmp {
	case where.EQ:
		return value == fc.Value[0].Value(), nil
	case where.InArray:
		for _, x := range fc.Value {
			if x.Value() == value {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, NewNotImplementComparatorError(fc.GetField(), fc.Cmp)
	}
}

func (fc Enum8FieldComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item).Value())
}

func (fc Enum8FieldComparator) ValuesCount() int {
	return len(fc.Value)
}

func (fc Enum8FieldComparator) ValueAt(index int) interface{} {
	return fc.Value[index]
}
