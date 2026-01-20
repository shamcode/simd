package comparators

import (
	"fmt"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type MapFieldComparator[R record.Record, K comparable, V any] struct {
	Cmp    where.ComparatorType
	Getter record.MapGetter[R, K, V]
	Value  []any
}

func (fc MapFieldComparator[R, K, V]) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc MapFieldComparator[R, K, V]) GetField() record.Field {
	return fc.Getter.Field
}

func (fc MapFieldComparator[R, K, V]) CompareValue(value record.Map[K, V]) (bool, error) {
	switch fc.Cmp { //nolint:exhaustive
	case where.MapHasValue:
		cmp, ok := fc.Value[0].(record.MapValueComparator[V])
		if !ok {
			return false, NewFailCastTypeError(fc.GetField(), fc.Cmp, fc.Value[0], "record.MapValueComparator")
		}

		return value.HasValue(cmp)
	case where.MapHasKey:
		val, ok := fc.Value[0].(K)
		if !ok {
			return false, NewFailCastTypeError(fc.GetField(), fc.Cmp, fc.Value[0], fmt.Sprintf("%T", val))
		}

		return value.HasKey(val), nil
	default:
		return false, NewNotImplementComparatorError(fc.GetField(), fc.Cmp)
	}
}

func (fc MapFieldComparator[R, K, V]) Compare(item R) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc MapFieldComparator[R, K, V]) ValuesCount() int {
	return len(fc.Value)
}

func (fc MapFieldComparator[R, K, V]) ValueAt(index int) any {
	return fc.Value[index]
}

func NewMapFieldComparator[R record.Record, K comparable, V any](
	cmp where.ComparatorType,
	getter record.MapGetter[R, K, V],
	value ...any,
) MapFieldComparator[R, K, V] {
	return MapFieldComparator[R, K, V]{
		Cmp:    cmp,
		Getter: getter,
		Value:  value,
	}
}
