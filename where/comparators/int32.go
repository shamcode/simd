package comparators

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type Int32FieldComparator struct {
	Cmp    where.ComparatorType
	Getter *record.Int32Getter
	Value  []int32
}

func (fc Int32FieldComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc Int32FieldComparator) GetField() record.Field {
	return fc.Getter.Field
}

func (fc Int32FieldComparator) CompareValue(value int32) (bool, error) {
	switch fc.Cmp {
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

func (fc Int32FieldComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc Int32FieldComparator) ValuesCount() int {
	return len(fc.Value)
}

func (fc Int32FieldComparator) ValueAt(index int) interface{} {
	return fc.Value[index]
}
