package compute

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type BoolKey bool

func (i BoolKey) Less(than indexes.Key) bool { return !bool(i) && bool(than.(BoolKey)) }

type boolComparator interface {
	CompareValue(value bool) (bool, error)
}

type boolIndexComputation[R record.Record] struct {
	getter record.BoolGetter[R]
}

func (idx boolIndexComputation[R]) ForRecord(item R) indexes.Key {
	return BoolKey(idx.getter.Get(item))
}

func (idx boolIndexComputation[R]) ForValue(value interface{}) indexes.Key {
	return BoolKey(value.(bool))
}

func (idx boolIndexComputation[R]) Check(
	indexKey indexes.Key,
	comparator where.FieldComparator[R],
) (bool, error) {
	return comparator.(boolComparator).CompareValue(bool(indexKey.(BoolKey)))
}

func CreateBoolIndexComputation[R record.Record](getter record.BoolGetter[R]) indexes.IndexComputer[R] {
	return boolIndexComputation[R]{getter: getter}
}
