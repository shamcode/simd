package fields

import (
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type SetFieldComparator struct {
	BaseFieldComparator
	Getter *record.SetGetter
	Value  []interface{}
}

func (fc *SetFieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc *SetFieldComparator) CompareValue(value record.Set) bool {
	switch fc.Cmp {
	case where.SetHas:
		return value.Has(fc.Value[0])
	default:
		panic(fmt.Errorf("%w: %d, field = %s", errNotImplementComparator, fc.Cmp, fc.GetField()))
	}
}

func (fc *SetFieldComparator) Compare(item interface{}) bool {
	return fc.CompareValue(fc.Getter.Get(item))
}
