package fields

import (
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type InterfaceFieldComparator struct {
	BaseFieldComparator
	Getter *record.InterfaceGetter
	Value  []interface{}
}

func (fc *InterfaceFieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc *InterfaceFieldComparator) CompareValue(value interface{}) bool {
	switch fc.Cmp {
	case where.EQ:
		return value == fc.Value[0]
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

func (fc *InterfaceFieldComparator) Compare(item interface{}) bool {
	return fc.CompareValue(fc.Getter.Get(item))
}
