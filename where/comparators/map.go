package comparators

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type MapFieldComparator struct {
	Cmp    where.ComparatorType
	Getter *record.MapGetter
	Value  []interface{}
}

func (fc MapFieldComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc MapFieldComparator) GetField() record.Field {
	return fc.Getter.Field
}

func (fc MapFieldComparator) CompareValue(value record.Map) (bool, error) {
	switch fc.Cmp {
	case where.MapHasValue:
		cmp, ok := fc.Value[0].(record.MapValueComparator)
		if !ok {
			return false, NewErrFailCastType(fc.GetField(), fc.Cmp, fc.Value[0], "record.MapValueComparator")
		}
		return value.HasValue(cmp)
	case where.MapHasKey:
		return value.HasKey(fc.Value[0]), nil
	default:
		return false, NewErrNotImplementComparator(fc.GetField(), fc.Cmp)
	}
}

func (fc MapFieldComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc MapFieldComparator) ValuesCount() int {
	return len(fc.Value)
}

func (fc MapFieldComparator) ValueAt(index int) interface{} {
	return fc.Value[index]
}
