package compute

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

type BoolKey bool

func (i BoolKey) Less(than indexes.Key) bool { return !bool(i) && bool(than.(BoolKey)) }

var _ indexes.IndexComputer = boolIndexComputation{}

type boolIndexComputation struct {
	getter *record.BoolGetter
}

func (idx boolIndexComputation) ForRecord(item record.Record) indexes.Key {
	return BoolKey(idx.getter.Get(item))
}

func (idx boolIndexComputation) ForValue(value interface{}) indexes.Key {
	return BoolKey(value.(bool))
}

func (idx boolIndexComputation) Check(indexKey indexes.Key, comparator where.FieldComparator) (bool, error) {
	return comparator.(comparators.BoolFieldComparator).CompareValue(bool(indexKey.(BoolKey)))
}

func CreateBoolIndexComputation(getter *record.BoolGetter) indexes.IndexComputer {
	return boolIndexComputation{getter: getter}
}
