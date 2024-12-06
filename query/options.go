package query

import (
	"regexp"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

type LimitOption int

func (o LimitOption) Apply(b any) { b.(Builder).Limit(int(o)) }

func Limit(limitItems int) BuilderOption {
	return LimitOption(limitItems)
}

type OffsetOption int

func (o OffsetOption) Apply(b any) { b.(Builder).Offset(int(o)) }
func Offset(startOffset int) BuilderOption {
	return OffsetOption(startOffset)
}

type OrOption struct{}

func (OrOption) Apply(b any) { b.(Builder).Or() }
func Or() OrOption {
	return OrOption{}
}

type NotOption struct{}

func (NotOption) Apply(b any) { b.(Builder).Not() }
func Not() BuilderOption {
	return NotOption{}
}

type OpenBracketOption struct{}

func (OpenBracketOption) Apply(b any) { b.(Builder).OpenBracket() }
func OpenBracket() BuilderOption {
	return OpenBracketOption{}
}

type CloseBracketOption struct{}

func (CloseBracketOption) Apply(b any) { b.(Builder).CloseBracket() }
func CloseBracket() BuilderOption {
	return CloseBracketOption{}
}

type ErrorOption struct {
	err error
}

func (o ErrorOption) Apply(b any) { b.(Builder).Error(o.err) }

func Error(err error) BuilderOption {
	return ErrorOption{err: err}
}

type SortOption[R record.Record] struct {
	by sort.ByWithOrder[R]
}

func (o SortOption[R]) Apply(b any) { b.(BuilderGeneric[R]).Sort(o.by) }
func Sort[R record.Record](by sort.ByWithOrder[R]) BuilderOption {
	return SortOption[R]{by: by}
}

type AddWhereOption[R record.Record] struct {
	Cmp where.FieldComparator[R]
}

func (o AddWhereOption[R]) Apply(b any) { b.(BuilderGeneric[R]).AddWhere(o.Cmp) }

func WhereAny[R record.Record](
	getter record.GetterInterface[R, any],
	condition where.ComparatorType,
	values ...any,
) BuilderOption {
	return AddWhereOption[R]{
		Cmp: comparators.EqualComparator[R, any]{
			Cmp:    condition,
			Getter: getter,
			Value:  values,
		},
	}
}

func Where[R record.Record, T record.LessComparable](
	getter record.ComparableGetter[R, T],
	condition where.ComparatorType,
	value ...T,
) BuilderOption {
	switch castedGetter := any(getter).(type) {
	case record.ComparableGetter[R, string]:
		castedValue, err := Cast[[]T, []string](value)
		if err != nil {
			return Error(GetterError{Field: getter.Field, Err: err})
		}

		return AddWhereOption[R]{
			Cmp: comparators.NewStringFieldComparator[R](
				condition,
				castedGetter,
				castedValue...,
			),
		}

	default:
		return AddWhereOption[R]{
			Cmp: comparators.NewComparableFieldComparator[R, T](condition, getter, value...),
		}
	}
}

func WhereStringRegexp[R record.Record](
	getter record.ComparableGetter[R, string],
	value *regexp.Regexp,
) BuilderOption {
	return AddWhereOption[R]{
		Cmp: comparators.NewStringFieldRegexpComparator[R](where.Regexp, getter, value),
	}
}

func WhereBool[R record.Record](
	getter record.BoolGetter[R],
	condition where.ComparatorType,
	value ...bool,
) BuilderOption {
	return AddWhereOption[R]{
		Cmp: comparators.EqualComparator[R, bool]{
			Cmp:    condition,
			Getter: record.Getter[R, bool](getter),
			Value:  value,
		},
	}
}

func WhereMap[R record.Record](
	getter record.MapGetter[R],
	condition where.ComparatorType,
	value ...any,
) BuilderOption {
	return AddWhereOption[R]{
		Cmp: comparators.NewMapFieldComparator[R](condition, getter, value...),
	}
}

func WhereSet[R record.Record](
	getter record.SetGetter[R],
	condition where.ComparatorType,
	value ...any,
) BuilderOption {
	return AddWhereOption[R]{
		Cmp: comparators.NewSetFieldComparator[R](condition, getter, value...),
	}
}

type OnIterationOption[R record.Record] func(item R)

func (o OnIterationOption[R]) Apply(b any) { b.(BuilderGeneric[R]).OnIteration(o) }

func OnIteration[R record.Record](cb func(item R)) OnIterationOption[R] {
	return cb
}
