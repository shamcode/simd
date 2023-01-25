package hash

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

var _ indexes.IndexComputer = enum8IndexComputation{}

type enum8IndexComputation struct {
	getter *record.Enum8Getter
}

func (idx enum8IndexComputation) ForRecord(item record.Record) interface{} {
	return idx.getter.Get(item).Value()
}

func (idx enum8IndexComputation) ForValue(value interface{}) interface{} {
	return value.(record.Enum8).Value()
}

func (idx enum8IndexComputation) Check(indexKey interface{}, comparator where.FieldComparator) (bool, error) {
	return comparator.(comparators.Enum8FieldComparator).CompareValue(indexKey.(uint8))
}

var _ HashTable = (*enum8HashIndexStorage)(nil)

type enum8HashIndexStorage struct {
	byValue map[uint8]storage.IDStorage
}

func (idx *enum8HashIndexStorage) Get(key interface{}) storage.IDStorage {
	return idx.byValue[key.(uint8)]
}

func (idx *enum8HashIndexStorage) Set(key interface{}, records storage.IDStorage) {
	idx.byValue[key.(uint8)] = records
}

func (idx *enum8HashIndexStorage) Keys() []interface{} {
	i := 0
	keys := make([]interface{}, len(idx.byValue))
	for key := range idx.byValue {
		keys[i] = key
		i += 1
	}
	return keys
}

func NewEnum8HashIndex(getter *record.Enum8Getter, unique bool) indexes.Index {
	return NewIndex(
		getter.Field,
		enum8IndexComputation{getter: getter},
		&enum8HashIndexStorage{
			byValue: make(map[uint8]storage.IDStorage),
		},
		unique,
	)
}
