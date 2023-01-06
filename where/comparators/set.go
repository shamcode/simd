package comparators

import (
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type SetFieldComparator struct {
	Cmp    where.ComparatorType
	Getter *record.SetGetter
	Value  []interface{}
}

func (fc SetFieldComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc SetFieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc SetFieldComparator) CompareValue(value record.Set) (bool, error) {
	switch fc.Cmp {
	case where.SetHas:
		return value.Has(fc.Value[0]), nil
	default:
		return false, fmt.Errorf("%w: %d, field = %s", ErrNotImplementComparator, fc.Cmp, fc.GetField())
	}
}

func (fc SetFieldComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc SetFieldComparator) Values() []interface{} {
	values := make([]interface{}, len(fc.Value))
	for i, v := range fc.Value {
		values[i] = v
	}
	return values
}
