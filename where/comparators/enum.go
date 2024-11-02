package comparators

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type EnumFieldComparator[R record.Record, T record.LessComparable] struct {
	Cmp    where.ComparatorType
	Getter record.EnumGetter[R, T]
	Value  []record.Enum[T]
}

func (fc EnumFieldComparator[R, T]) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc EnumFieldComparator[R, T]) GetField() record.Field {
	return fc.Getter
}

func (fc EnumFieldComparator[R, T]) CompareValue(value T) (bool, error) {
	switch fc.Cmp { //nolint:exhaustive
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

func (fc EnumFieldComparator[R, T]) Compare(item R) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item).Value())
}

func (fc EnumFieldComparator[R, T]) ValuesCount() int {
	return len(fc.Value)
}

func (fc EnumFieldComparator[R, T]) ValueAt(index int) interface{} {
	return fc.Value[index]
}

func NewEnumFieldComparator[R record.Record, T record.LessComparable](
	cmp where.ComparatorType,
	getter record.EnumGetter[R, T],
	values ...record.Enum[T],
) EnumFieldComparator[R, T] {
	return EnumFieldComparator[R, T]{
		Cmp:    cmp,
		Getter: getter,
		Value:  values,
	}
}
