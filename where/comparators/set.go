package comparators

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type SetFieldComparator struct {
	Cmp    where.ComparatorType
	Getter record.SetGetter
	Value  []interface{}
}

func (fc SetFieldComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc SetFieldComparator) GetField() record.Field {
	return fc.Getter.Field
}

func (fc SetFieldComparator) CompareValue(value record.Set) (bool, error) {
	switch fc.Cmp {
	case where.SetHas:
		return value.Has(fc.Value[0]), nil
	default:
		return false, NewErrNotImplementComparator(fc.GetField(), fc.Cmp)
	}
}

func (fc SetFieldComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc SetFieldComparator) ValuesCount() int {
	return len(fc.Value)
}

func (fc SetFieldComparator) ValueAt(index int) interface{} {
	return fc.Value[index]
}
