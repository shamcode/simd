package comparators

import (
	"fmt"
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

func (fc MapFieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc MapFieldComparator) CompareValue(value record.Map) (bool, error) {
	switch fc.Cmp {
	case where.MapHasValue:
		cmp, ok := fc.Value[0].(record.MapValueComparator)
		if !ok {
			return false, fmt.Errorf(
				"%w: %d, field = %s, value type = %T, expected type = record.MapValueComparator",
				ErrFailCastType,
				fc.Cmp,
				fc.GetField(),
				fc.Value[0],
			)
		}
		return value.HasValue(cmp)
	case where.MapHasKey:
		return value.HasKey(fc.Value[0]), nil
	default:
		return false, fmt.Errorf("%w: %d, field = %s", ErrNotImplementComparator, fc.Cmp, fc.GetField())
	}
}

func (fc MapFieldComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc MapFieldComparator) Values() []interface{} {
	values := make([]interface{}, len(fc.Value))
	for i, v := range fc.Value {
		values[i] = v
	}
	return values
}
