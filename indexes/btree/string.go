package btree

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

var _ indexes.IndexComputer = stringIndexComputation{}

type stringComparator interface {
	CompareValue(value string) (bool, error)
}

type _string string

func (i _string) Less(than Key) bool { return i < than.(_string) }

type stringIndexComputation struct {
	getter *record.StringGetter
}

func (idx stringIndexComputation) ForRecord(item record.Record) interface{} {
	return _string(idx.getter.Get(item))
}

func (idx stringIndexComputation) ForValue(value interface{}) interface{} {
	return _string(value.(string))
}

func (idx stringIndexComputation) Check(indexKey interface{}, comparator where.FieldComparator) (bool, error) {
	return comparator.(stringComparator).CompareValue(indexKey.(string))
}

func NewStringBTreeIndex(getter *record.StringGetter, maxChildren int, uniq bool) indexes.Index {
	return NewIndex(
		getter.Field,
		stringIndexComputation{getter: getter},
		NewTree(maxChildren, uniq),
		uniq,
	)
}
