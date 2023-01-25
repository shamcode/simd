package btree

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

var _ indexes.IndexComputer = int32IndexComputation{}

type _int32 int32

func (i _int32) Less(than Key) bool { return i < than.(_int32) }

type int32IndexComputation struct {
	getter *record.Int32Getter
}

func (idx int32IndexComputation) ForRecord(item record.Record) interface{} {
	return _int32(idx.getter.Get(item))
}

func (idx int32IndexComputation) ForValue(value interface{}) interface{} {
	return _int32(value.(int32))
}

func (idx int32IndexComputation) Check(indexKey interface{}, comparator where.FieldComparator) (bool, error) {
	return comparator.(comparators.Int32FieldComparator).CompareValue(indexKey.(int32))
}

func NewInt32BTreeIndex(getter *record.Int32Getter, maxChildren int, uniq bool) indexes.Index {
	return NewIndex(
		getter.Field,
		int32IndexComputation{getter: getter},
		NewTree(maxChildren, uniq),
		uniq,
	)
}
