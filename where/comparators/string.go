package comparators

import (
	"regexp"
	"strings"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type StringFieldComparator struct {
	Cmp    where.ComparatorType
	Getter record.StringGetter
	Value  []string
}

func (fc StringFieldComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc StringFieldComparator) GetField() record.Field {
	return fc.Getter.Field
}

func (fc StringFieldComparator) CompareValue(value string) (bool, error) { //nolint:cyclop
	switch fc.Cmp { //nolint:exhaustive
	case where.EQ:
		return value == fc.Value[0], nil
	case where.GT:
		return value > fc.Value[0], nil
	case where.LT:
		return value < fc.Value[0], nil
	case where.GE:
		return value >= fc.Value[0], nil
	case where.LE:
		return value <= fc.Value[0], nil
	case where.Like:
		return strings.Contains(value, fc.Value[0]), nil
	case where.InArray:
		for _, x := range fc.Value {
			if x == value {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, NewNotImplementComparatorError(fc.GetField(), fc.Cmp)
	}
}

func (fc StringFieldComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc StringFieldComparator) ValuesCount() int {
	return len(fc.Value)
}

func (fc StringFieldComparator) ValueAt(index int) interface{} {
	return fc.Value[index]
}

// StringFieldRegexpComparator is a special comparator for handling Regexp.
type StringFieldRegexpComparator struct {
	Cmp    where.ComparatorType
	Getter record.StringGetter
	Value  *regexp.Regexp
}

func (fc StringFieldRegexpComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc StringFieldRegexpComparator) GetField() record.Field {
	return fc.Getter.Field
}

func (fc StringFieldRegexpComparator) CompareValue(value string) (bool, error) {
	switch fc.Cmp { //nolint:exhaustive
	case where.Regexp:
		return fc.Value.MatchString(value), nil
	default:
		return false, NewNotImplementComparatorError(fc.GetField(), fc.Cmp)
	}
}

func (fc StringFieldRegexpComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc StringFieldRegexpComparator) ValuesCount() int {
	return 1
}

func (fc StringFieldRegexpComparator) ValueAt(index int) interface{} {
	if index == 0 {
		return fc.Value
	}
	return nil
}
