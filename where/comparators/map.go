package comparators

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type MapFieldComparator[R record.Record] struct {
	Cmp    where.ComparatorType
	Getter record.MapGetter[R]
	Value  []any
}

func (fc MapFieldComparator[R]) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc MapFieldComparator[R]) GetField() record.Field {
	return fc.Getter.Field
}

func (fc MapFieldComparator[R]) CompareValue(value record.Map) (bool, error) {
	switch fc.Cmp { //nolint:exhaustive
	case where.MapHasValue:
		cmp, ok := fc.Value[0].(record.MapValueComparator)
		if !ok {
			return false, NewFailCastTypeError(fc.GetField(), fc.Cmp, fc.Value[0], "record.MapValueComparator")
		}
		return value.HasValue(cmp)
	case where.MapHasKey:
		return value.HasKey(fc.Value[0]), nil
	default:
		return false, NewNotImplementComparatorError(fc.GetField(), fc.Cmp)
	}
}

func (fc MapFieldComparator[R]) Compare(item R) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc MapFieldComparator[R]) ValuesCount() int {
	return len(fc.Value)
}

func (fc MapFieldComparator[R]) ValueAt(index int) any {
	return fc.Value[index]
}

func NewMapFieldComparator[R record.Record](
	cmp where.ComparatorType,
	getter record.MapGetter[R],
	value ...any,
) MapFieldComparator[R] {
	return MapFieldComparator[R]{
		Cmp:    cmp,
		Getter: getter,
		Value:  value,
	}
}
