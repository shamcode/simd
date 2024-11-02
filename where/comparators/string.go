package comparators

import (
	"regexp"
	"strings"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type (
	// StringFieldComparator is a comparator for string field.
	StringFieldComparator[R record.Record] struct {
		ComparableFieldComparator[R, string]
	}

	// StringFieldRegexpComparator is a special comparator for handling Regexp.
	StringFieldRegexpComparator[R record.Record] struct {
		Cmp    where.ComparatorType
		Getter record.ComparableGetter[R, string]
		Value  *regexp.Regexp
	}
)

func (fc StringFieldComparator[R]) CompareValue(value string) (bool, error) {
	switch fc.Cmp { //nolint:exhaustive
	case where.Like:
		return strings.Contains(value, fc.Value[0]), nil
	default:
		return fc.ComparableFieldComparator.CompareValue(value)
	}
}

func (fc StringFieldComparator[R]) Compare(item R) (bool, error) {
	return fc.CompareValue(fc.Getter.GetForRecord(item))
}

func (fc StringFieldRegexpComparator[R]) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc StringFieldRegexpComparator[R]) GetField() record.Field {
	return fc.Getter.Field
}

func (fc StringFieldRegexpComparator[R]) CompareValue(value string) (bool, error) {
	switch fc.Cmp { //nolint:exhaustive
	case where.Regexp:
		return fc.Value.MatchString(value), nil
	default:
		return false, NewNotImplementComparatorError(fc.GetField(), fc.Cmp)
	}
}

func (fc StringFieldRegexpComparator[R]) Compare(item R) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc StringFieldRegexpComparator[R]) ValuesCount() int {
	return 1
}

func (fc StringFieldRegexpComparator[R]) ValueAt(index int) interface{} {
	if index == 0 {
		return fc.Value
	}
	return nil
}

func NewStringFieldComparator[R record.Record](
	cmp where.ComparatorType,
	getter record.ComparableGetter[R, string],
	values ...string,
) StringFieldComparator[R] {
	return StringFieldComparator[R]{
		ComparableFieldComparator: ComparableFieldComparator[R, string]{
			EqualComparator: EqualComparator[R, string]{
				Cmp:    cmp,
				Getter: getter,
				Value:  values,
			},
		},
	}
}

func NewStringFieldRegexpComparator[R record.Record](
	cmp where.ComparatorType,
	getter record.ComparableGetter[R, string],
	value *regexp.Regexp,
) StringFieldRegexpComparator[R] {
	return StringFieldRegexpComparator[R]{
		Cmp:    cmp,
		Getter: getter,
		Value:  value,
	}
}
