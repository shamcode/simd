package comparators

import (
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type Enum16FieldComparator struct {
	Cmp    where.ComparatorType
	Getter *record.Enum16Getter
	Value  []record.Enum16
}

func (fc Enum16FieldComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc Enum16FieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc Enum16FieldComparator) CompareValue(value uint16) (bool, error) {
	switch fc.Cmp {
	case where.EQ:
		return value == fc.Value[0].Value(), nil
	case where.InArray:
		for _, x := range fc.Value {
			if x.Value() == value {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("%w: %d, field = %s", ErrNotImplementComparator, fc.Cmp, fc.GetField())
	}
}

func (fc Enum16FieldComparator) Compare(item interface{}) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item).Value())
}
