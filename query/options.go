package query

import (
	"regexp"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

type WhereOption[R record.Record] struct {
	Cmp   where.FieldComparator[R]
	Error error
}

func FieldAny[R record.Record](
	getter record.GetterInterface[R, any],
	condition where.ComparatorType,
	values ...any,
) WhereOption[R] {
	return WhereOption[R]{
		Cmp: comparators.EqualComparator[R, any]{
			Cmp:    condition,
			Getter: getter,
			Value:  values,
		},
		Error: nil,
	}
}

func Field[R record.Record, T record.LessComparable](
	getter record.ComparableGetter[R, T],
	condition where.ComparatorType,
	value ...T,
) WhereOption[R] {
	switch castedGetter := any(getter).(type) {
	case record.ComparableGetter[R, string]:
		castedValue, err := Cast[[]T, []string](value)
		if err != nil {
			return WhereOption[R]{
				Cmp:   nil,
				Error: GetterError{Field: getter.Field, Err: err},
			}
		}

		return WhereOption[R]{
			Cmp: comparators.NewStringFieldComparator[R](
				condition,
				castedGetter,
				castedValue...,
			),
			Error: nil,
		}

	default:
		return WhereOption[R]{
			Cmp:   comparators.NewComparableFieldComparator[R, T](condition, getter, value...),
			Error: nil,
		}
	}
}

func FieldStringRegexp[R record.Record](
	getter record.ComparableGetter[R, string],
	value *regexp.Regexp,
) WhereOption[R] {
	return WhereOption[R]{
		Cmp:   comparators.NewStringFieldRegexpComparator[R](where.Regexp, getter, value),
		Error: nil,
	}
}

func FieldBool[R record.Record](
	getter record.BoolGetter[R],
	condition where.ComparatorType,
	value ...bool,
) WhereOption[R] {
	return WhereOption[R]{
		Cmp: comparators.EqualComparator[R, bool]{
			Cmp:    condition,
			Getter: record.Getter[R, bool](getter),
			Value:  value,
		},
		Error: nil,
	}
}

func FieldMap[R record.Record, K comparable, V any](
	getter record.MapGetter[R, K, V],
	condition where.ComparatorType,
	value ...any,
) WhereOption[R] {
	return WhereOption[R]{
		Cmp:   comparators.NewMapFieldComparator[R, K, V](condition, getter, value...),
		Error: nil,
	}
}

func FieldSet[R record.Record, T comparable](
	getter record.SetGetter[R, T],
	condition where.ComparatorType,
	value ...T,
) WhereOption[R] {
	return WhereOption[R]{
		Cmp:   comparators.NewSetFieldComparator[R, T](condition, getter, value...),
		Error: nil,
	}
}
