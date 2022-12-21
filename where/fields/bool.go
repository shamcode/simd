package fields

import (
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type BoolFieldComparator struct {
	BaseFieldComparator
	Getter *record.BoolGetter
	Value  []bool
}

func (fc *BoolFieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc *BoolFieldComparator) CompareValue(value bool) bool {
	switch fc.Cmp {
	case where.EQ:
		return value == fc.Value[0]
	default:
		panic(fmt.Errorf("%w: %d, field = %s", errNotImplementComparator, fc.Cmp, fc.GetField()))
	}
}

func (fc *BoolFieldComparator) Compare(item interface{}) bool {
	return fc.CompareValue(fc.Getter.Get(item))
}
