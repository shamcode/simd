package comparators

import (
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type Int32FieldComparator struct {
	BaseFieldComparator
	Getter *record.Int32Getter
	Value  []int32
}

func (fc *Int32FieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc *Int32FieldComparator) CompareValue(value int32) bool {
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

func (fc *Int32FieldComparator) Compare(item interface{}) bool {
	return fc.CompareValue(fc.Getter.Get(item))
}
