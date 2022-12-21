package fields

import (
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type StringsSetFieldComparator struct {
	BaseFieldComparator
	Getter *record.StringsSetGetter
	Value  []string
}

func (fc *StringsSetFieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc *StringsSetFieldComparator) CompareValue(value record.StringsSet) bool {
	switch fc.Cmp {
	case where.SetHas:
		_, ok := value.Set[fc.Value[0]]
		return ok
	default:
		panic(fmt.Errorf("%w: %d, field = %s", errNotImplementComparator, fc.Cmp, fc.GetField()))
	}
}

func (fc *StringsSetFieldComparator) Compare(item interface{}) bool {
	return fc.CompareValue(fc.Getter.Get(item))
}
