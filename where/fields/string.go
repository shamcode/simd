package fields

import (
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"regexp"
	"strings"
)

type StringFieldComparator struct {
	BaseFieldComparator
	Getter *record.StringGetter
	Value  []string
}

func (fc *StringFieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc *StringFieldComparator) CompareValue(value string) bool {
	switch fc.Cmp {
	case where.EQ:
		return value == fc.Value[0]
	case where.GT:
		return value > fc.Value[0]
	case where.LT:
		return value < fc.Value[0]
	case where.GE:
		return value >= fc.Value[0]
	case where.LE:
		return value <= fc.Value[0]
	case where.Like:
		return strings.Contains(value, fc.Value[0])
	case where.InArray:
		for _, x := range fc.Value {
			if x == value {
				return true
			}
		}
		return false
	default:
		panic(fmt.Errorf("%w: %d, field = %s", errNotImplementComparator, fc.Cmp, fc.GetField()))
	}
}

func (fc *StringFieldComparator) Compare(item interface{}) bool {
	return fc.CompareValue(fc.Getter.Get(item))
}

type StringFieldRegexpComparator struct {
	BaseFieldComparator
	Getter *record.StringGetter
	Value  *regexp.Regexp
}

func (fc *StringFieldRegexpComparator) GetField() string {
	return fc.Getter.Field
}

func (fc *StringFieldRegexpComparator) CompareValue(value string) bool {
	switch fc.Cmp {
	case where.Regexp:
		return fc.Value.MatchString(value)
	default:
		panic(fmt.Errorf("%w: %d, field = %s", errNotImplementComparator, fc.Cmp, fc.GetField()))
	}
}

func (fc *StringFieldRegexpComparator) Compare(item interface{}) bool {
	return fc.CompareValue(fc.Getter.Get(item))
}
