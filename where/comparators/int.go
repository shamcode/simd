package comparators

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type IntFieldComparator struct {
	Cmp    where.ComparatorType
	Getter record.IntGetter
	Value  []int
}

func (fc IntFieldComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc IntFieldComparator) GetField() record.Field {
	return fc.Getter.Field
}

func (fc IntFieldComparator) CompareValue(value int) (bool, error) {
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

func (fc IntFieldComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc IntFieldComparator) ValuesCount() int {
	return len(fc.Value)
}

func (fc IntFieldComparator) ValueAt(index int) interface{} {
	return fc.Value[index]
}
