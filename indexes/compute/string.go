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

type stringIndexComputation[R record.Record] struct {
	getter record.StringGetter[R]
}

func (idx stringIndexComputation[R]) ForRecord(item R) indexes.Key {
	return StringKey(idx.getter.Get(item))
}

func (idx stringIndexComputation[R]) ForValue(value interface{}) indexes.Key {
	return StringKey(value.(string))
}

func (idx stringIndexComputation[R]) Check(
	indexKey indexes.Key,
	comparator where.FieldComparator[R],
) (bool, error) {
	return comparator.(stringComparator).CompareValue(string(indexKey.(StringKey)))
}

func CreateStringIndexComputation[R record.Record](getter record.StringGetter[R]) indexes.IndexComputer[R] {
	return stringIndexComputation[R]{getter: getter}
}
