package btree

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

var _ indexes.IndexComputer = int64IndexComputation{}

type _int64 int64

func (i _int64) Less(than Key) bool { return i < than.(_int64) }

type int64IndexComputation struct {
	getter *record.Int64Getter
}

func (idx int64IndexComputation) ForRecord(item record.Record) interface{} {
	return _int64(idx.getter.Get(item))
}

func (idx int64IndexComputation) ForValue(value interface{}) interface{} {
	return _int64(value.(int64))
}

func (idx int64IndexComputation) Check(indexKey interface{}, comparator where.FieldComparator) (bool, error) {
	return comparator.(comparators.Int64FieldComparator).CompareValue(indexKey.(int64))
}

func NewInt64BTreeIndex(getter *record.Int64Getter, maxChildren int, uniq bool) indexes.Index {
	return NewIndex(
		getter.Field,
		int64IndexComputation{getter: getter},
		NewTree(maxChildren, uniq),
		uniq,
	)
}
