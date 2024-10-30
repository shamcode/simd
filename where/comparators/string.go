package comparators

import (
	"regexp"
	"strings"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type StringFieldComparator[R record.Record] struct {
	ComparableFieldComparator[R, string]
}

func (fc StringFieldComparator[R]) CompareValue(value string) (bool, error) {
	switch fc.Cmp { //nolint:exhaustive
	case where.Like:
		return strings.Contains(value, fc.Value[0]), nil
	default:
		return fc.ComparableFieldComparator.CompareValue(value)
	}
}

func (fc StringFieldComparator[R]) Compare(item R) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

// StringFieldRegexpComparator is a special comparator for handling Regexp.
type StringFieldRegexpComparator[R record.Record] struct {
	Cmp    where.ComparatorType
	Getter record.ComparableGetter[R, string]
	Value  *regexp.Regexp
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
