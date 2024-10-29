package comparators

import (
	"time"

	"github.com/shamcode/simd/examples/custom-field-time/types"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

type TimeFieldComparator[R record.Record] struct {
	Cmp    where.ComparatorType
	Getter types.TimeGetter[R]
	Value  []time.Time
}

func (fc TimeFieldComparator[R]) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc TimeFieldComparator[R]) GetField() record.Field {
	return fc.Getter.Field
}

func (fc TimeFieldComparator[R]) CompareValue(value time.Time) (bool, error) {
	switch fc.Cmp { //nolint:exhaustive
	case where.EQ:
		return value.Equal(fc.Value[0]), nil
	case where.GT:
		return value.After(fc.Value[0]), nil
	case where.LT:
		return value.Before(fc.Value[0]), nil
	case where.GE:
		return value.Equal(fc.Value[0]) || value.After(fc.Value[0]), nil
	case where.LE:
		return value.Equal(fc.Value[0]) || value.Before(fc.Value[0]), nil
	default:
		return false, comparators.NewNotImplementComparatorError(fc.GetField(), fc.Cmp)
	}
}

func (fc TimeFieldComparator[R]) Compare(item R) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc TimeFieldComparator[R]) ValuesCount() int {
	return len(fc.Value)
}

func (fc TimeFieldComparator[R]) ValueAt(index int) interface{} {
	return fc.Value[index]
}
