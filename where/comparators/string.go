package comparators

import (
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"regexp"
	"strings"
)

type StringFieldComparator struct {
	Cmp    where.ComparatorType
	Getter *record.StringGetter
	Value  []string
}

func (fc StringFieldComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc StringFieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc StringFieldComparator) CompareValue(value string) (bool, error) {
	switch fc.Cmp {
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
		return false, fmt.Errorf("%w: %d, field = %s", ErrNotImplementComparator, fc.Cmp, fc.GetField())
	}
}

func (fc StringFieldComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc StringFieldComparator) Values() []interface{} {
	values := make([]interface{}, len(fc.Value))
	for i, v := range fc.Value {
		values[i] = v
	}
	return values
}

// StringFieldRegexpComparator is a special comparator for handling Regexp
type StringFieldRegexpComparator struct {
	Cmp    where.ComparatorType
	Getter *record.StringGetter
	Value  *regexp.Regexp
}

func (fc StringFieldRegexpComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc StringFieldRegexpComparator) GetField() string {
	return fc.Getter.Field
}

func (fc StringFieldRegexpComparator) CompareValue(value string) (bool, error) {
	switch fc.Cmp {
	case where.Regexp:
		return fc.Value.MatchString(value), nil
	default:
		return false, fmt.Errorf("%w: %d, field = %s", ErrNotImplementComparator, fc.Cmp, fc.GetField())
	}
}

func (fc StringFieldRegexpComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc StringFieldRegexpComparator) Values() []interface{} {
	return []interface{}{fc.Value}
}
