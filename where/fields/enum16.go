package fields

import (
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type Enum16FieldComparator struct {
	BaseFieldComparator
	Getter *record.Enum16Getter
	Value  []record.Enum16
}

func (fc *Enum16FieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc *Enum16FieldComparator) CompareValue(value uint16) bool {
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

func (fc *Enum16FieldComparator) Compare(item interface{}) bool {
	return fc.CompareValue(fc.Getter.Get(item).Value())
}
