package comparators

import (
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type IntFieldComparator struct {
	BaseFieldComparator
	Getter *record.IntGetter
	Value  []int
}

func (fc *IntFieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc *IntFieldComparator) CompareValue(value int) bool {
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

func (fc *IntFieldComparator) Compare(item interface{}) bool {
	return fc.CompareValue(fc.Getter.Get(item))
}
