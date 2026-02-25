package query

import (
	"regexp"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

type AddWhereOption[R record.Record] struct {
	Cmp   where.FieldComparator[R]
	Error error
}

func (o AddWhereOption[R]) Apply(b BuilderGeneric[R]) {
	if o.Error == nil {
		b.AddWhere(o.Cmp)
	} else {
		b.Error(o.Error)
	}
}

func WhereAny[R record.Record](
	getter record.GetterInterface[R, any],
	condition where.ComparatorType,
	values ...any,
) AddWhereOption[R] {
	return AddWhereOption[R]{
		Cmp: comparators.EqualComparator[R, any]{
			Cmp:    condition,
			Getter: getter,
			Value:  values,
		},
		Error: nil,
	}
}

func Where[R record.Record, T record.LessComparable](
	getter record.ComparableGetter[R, T],
	condition where.ComparatorType,
	value ...T,
) AddWhereOption[R] {
	switch castedGetter := any(getter).(type) {
	case record.ComparableGetter[R, string]:
		castedValue, err := Cast[[]T, []string](value)
		if err != nil {
			return AddWhereOption[R]{
				Cmp:   nil,
				Error: GetterError{Field: getter.Field, Err: err},
			}
		}

		return AddWhereOption[R]{
			Cmp: comparators.NewStringFieldComparator[R](
				condition,
				castedGetter,
				castedValue...,
			),
			Error: nil,
		}

	default:
		return AddWhereOption[R]{
			Cmp:   comparators.NewComparableFieldComparator[R, T](condition, getter, value...),
			Error: nil,
		}
	}
}

func WhereStringRegexp[R record.Record](
	getter record.ComparableGetter[R, string],
	value *regexp.Regexp,
) AddWhereOption[R] {
	return AddWhereOption[R]{
		Cmp:   comparators.NewStringFieldRegexpComparator[R](where.Regexp, getter, value),
		Error: nil,
	}
}

func WhereBool[R record.Record](
	getter record.BoolGetter[R],
	condition where.ComparatorType,
	value ...bool,
) AddWhereOption[R] {
	return AddWhereOption[R]{
		Cmp: comparators.EqualComparator[R, bool]{
			Cmp:    condition,
			Getter: record.Getter[R, bool](getter),
			Value:  value,
		},
		Error: nil,
	}
}

func WhereMap[R record.Record, K comparable, V any](
	getter record.MapGetter[R, K, V],
	condition where.ComparatorType,
	value ...any,
) AddWhereOption[R] {
	return AddWhereOption[R]{
		Cmp:   comparators.NewMapFieldComparator[R, K, V](condition, getter, value...),
		Error: nil,
	}
}

func WhereSet[R record.Record, T comparable](
	getter record.SetGetter[R, T],
	condition where.ComparatorType,
	value ...T,
) AddWhereOption[R] {
	return AddWhereOption[R]{
		Cmp:   comparators.NewSetFieldComparator[R, T](condition, getter, value...),
		Error: nil,
	}
}
