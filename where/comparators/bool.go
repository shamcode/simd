package comparators

import (
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type BoolFieldComparator struct {
	Cmp    where.ComparatorType
	Getter *record.BoolGetter
	Value  []bool
}

func (fc BoolFieldComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc BoolFieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc BoolFieldComparator) CompareValue(value bool) (bool, error) {
	switch fc.Cmp {
	case where.EQ:
		return value == fc.Value[0], nil
	default:
		return false, fmt.Errorf("%w: %d, field = %s", ErrNotImplementComparator, fc.Cmp, fc.GetField())
	}
}

func (fc BoolFieldComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc BoolFieldComparator) Values() []interface{} {
	values := make([]interface{}, len(fc.Value))
	for i, v := range fc.Value {
		values[i] = v
	}
	return values
}
