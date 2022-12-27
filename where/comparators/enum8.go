package comparators

import (
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type Enum8FieldComparator struct {
	BaseFieldComparator
	Getter *record.Enum8Getter
	Value  []record.Enum8
}

func (fc *Enum8FieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc *Enum8FieldComparator) CompareValue(value uint8) bool {
	switch fc.Cmp {
	case where.EQ:
		return value == fc.Value[0].Value()
	case where.InArray:
		for _, x := range fc.Value {
			if x.Value() == value {
				return true
			}
		}
		return false
	default:
		panic(fmt.Errorf("%w: %d, field = %s", errNotImplementComparator, fc.Cmp, fc.GetField()))
	}
}

func (fc *Enum8FieldComparator) Compare(item interface{}) bool {
	return fc.CompareValue(fc.Getter.Get(item).Value())
}
