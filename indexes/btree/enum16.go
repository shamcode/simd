package btree

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

var _ indexes.IndexComputer = int64IndexComputation{}

type _uint16 uint16

func (i _uint16) Less(than Key) bool { return i < than.(_uint16) }

type enum16IndexComputation struct {
	getter *record.Enum16Getter
}

func (idx enum16IndexComputation) ForRecord(item record.Record) interface{} {
	return _uint16(idx.getter.Get(item).Value())
}

func (idx enum16IndexComputation) ForValue(value interface{}) interface{} {
	return _uint16(value.(record.Enum16).Value())
}

func (idx enum16IndexComputation) Check(indexKey interface{}, comparator where.FieldComparator) (bool, error) {
	return comparator.(comparators.Enum16FieldComparator).CompareValue(indexKey.(uint16))
}

func NewEnum16BTreeIndex(getter *record.Enum16Getter, maxChildren int, uniq bool) indexes.Index {
	return NewIndex(
		getter.Field,
		enum16IndexComputation{getter: getter},
		NewTree(maxChildren, uniq),
		uniq,
	)
}
