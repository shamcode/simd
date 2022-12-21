package fields

import (
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type MapFieldComparator struct {
	BaseFieldComparator
	Getter *record.MapGetter
	Value  []interface{}
}

func (fc *MapFieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc *MapFieldComparator) CompareValue(value record.Map) bool {
	switch fc.Cmp {
	case where.MapHasValue:
		return value.HasValue(fc.Value[0].(record.Comparator))
	case where.MapHasKey:
		return value.HasKey(fc.Value[0])
	default:
		panic(fmt.Errorf("%w: %d, field = %s", errNotImplementComparator, fc.Cmp, fc.GetField()))
	}
}

func (fc *MapFieldComparator) Compare(item interface{}) bool {
	return fc.CompareValue(fc.Getter.Get(item))
}
