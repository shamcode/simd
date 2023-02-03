package comparators

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type Int64FieldComparator struct {
	Cmp    where.ComparatorType
	Getter record.Int64Getter
	Value  []int64
}

func (fc Int64FieldComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc Int64FieldComparator) GetField() record.Field {
	return fc.Getter.Field
}

func (fc Int64FieldComparator) CompareValue(value int64) (bool, error) {
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

func (fc Int64FieldComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc Int64FieldComparator) ValuesCount() int {
	return len(fc.Value)
}

func (fc Int64FieldComparator) ValueAt(index int) interface{} {
	return fc.Value[index]
}
