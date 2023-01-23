package hash

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

var _ indexes.IndexComputer = enum16IndexComputation{}

type enum16IndexComputation struct {
	getter *record.Enum16Getter
}

func (idx enum16IndexComputation) ForRecord(item record.Record) interface{} {
	return idx.getter.Get(item).Value()
}

func (idx enum16IndexComputation) ForValue(value interface{}) interface{} {
	return value.(record.Enum16).Value()
}

func (idx enum16IndexComputation) Check(indexKey interface{}, comparator where.FieldComparator) (bool, error) {
	return comparator.(comparators.Enum16FieldComparator).CompareValue(indexKey.(uint16))
}

var _ indexes.Storage = (*enum16HashIndexStorage)(nil)

type enum16HashIndexStorage struct {
	byValue map[uint16]*storage.IDStorage
}

func (idx *enum16HashIndexStorage) Get(key interface{}) *storage.IDStorage {
	return idx.byValue[key.(uint16)]
}

func (idx *enum16HashIndexStorage) Set(key interface{}, records *storage.IDStorage) {
	idx.byValue[key.(uint16)] = records
}

func (idx *enum16HashIndexStorage) Keys() []interface{} {
	i := 0
	keys := make([]interface{}, len(idx.byValue))
	for key := range idx.byValue {
		keys[i] = key
		i += 1
	}
	return keys
}

func NewEnum16HashIndex(getter *record.Enum16Getter) indexes.Index {
	return NewIndex(
		getter.Field,
		enum16IndexComputation{getter: getter},
		indexes.WrapToThreadSafeStorage(&enum16HashIndexStorage{
			byValue: make(map[uint16]*storage.IDStorage),
		}),
	)
}
