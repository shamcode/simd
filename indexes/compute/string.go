package compute

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type StringKey string

func (i StringKey) Less(than indexes.Key) bool { return i < than.(StringKey) }

type stringComparator interface {
	CompareValue(value string) (bool, error)
}

var _ indexes.IndexComputer = stringIndexComputation{}

type stringIndexComputation struct {
	getter *record.StringGetter
}

func (idx stringIndexComputation) ForRecord(item record.Record) indexes.Key {
	return StringKey(idx.getter.Get(item))
}

func (idx stringIndexComputation) ForValue(value interface{}) indexes.Key {
	return StringKey(value.(string))
}

func (idx stringIndexComputation) Check(indexKey indexes.Key, comparator where.FieldComparator) (bool, error) {
	return comparator.(stringComparator).CompareValue(string(indexKey.(StringKey)))
}

func CreateStringIndexComputation(getter *record.StringGetter) indexes.IndexComputer {
	return stringIndexComputation{getter: getter}
}
