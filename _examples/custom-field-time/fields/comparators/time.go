package comparators

import (
	"github.com/shamcode/simd/_examples/custom-field-time/fields"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
	"time"
)

type TimeFieldComparator struct {
	Cmp    where.ComparatorType
	Getter *fields.TimeGetter
	Value  []time.Time
}

func (fc TimeFieldComparator) GetType() where.ComparatorType {
	return fc.Cmp
}

func (fc TimeFieldComparator) GetField() string {
	return fc.Getter.Field
}

func (fc TimeFieldComparator) CompareValue(value time.Time) (bool, error) {
	switch fc.Cmp {
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
		return false, comparators.NewErrNotImplementComparator(fc.GetField(), fc.Cmp)
	}
}

func (fc TimeFieldComparator) Compare(item record.Record) (bool, error) {
	return fc.CompareValue(fc.Getter.Get(item))
}

func (fc TimeFieldComparator) ValuesCount() int {
	return len(fc.Value)
}

func (fc TimeFieldComparator) ValueAt(index int) interface{} {
	return fc.Value[index]
}
