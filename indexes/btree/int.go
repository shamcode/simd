package btree

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

var _ indexes.IndexComputer = intIndexComputation{}

type _int int

func (i _int) Less(than Key) bool { return i < than.(_int) }

type intIndexComputation struct {
	getter *record.IntGetter
}

func (idx intIndexComputation) ForRecord(item record.Record) interface{} {
	return _int(idx.getter.Get(item))
}

func (idx intIndexComputation) ForValue(value interface{}) interface{} {
	return _int(value.(int))
}

func (idx intIndexComputation) Check(indexKey interface{}, comparator where.FieldComparator) (bool, error) {
	return comparator.(comparators.IntFieldComparator).CompareValue(indexKey.(int))
}

func NewIntBTreeIndex(getter *record.IntGetter, maxChildren int, uniq bool) indexes.Index {
	return NewIndex(
		getter.Field,
		intIndexComputation{getter: getter},
		NewTree(maxChildren, uniq),
		uniq,
	)
}
