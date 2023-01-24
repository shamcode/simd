package hash

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

var _ indexes.IndexComputer = boolIndexComputation{}

type boolIndexComputation struct {
	getter *record.BoolGetter
}

func (idx boolIndexComputation) ForRecord(item record.Record) interface{} {
	return idx.getter.Get(item)
}

func (idx boolIndexComputation) ForValue(value interface{}) interface{} {
	return value.(bool)
}

func (idx boolIndexComputation) Check(indexKey interface{}, comparator where.FieldComparator) (bool, error) {
	return comparator.(comparators.BoolFieldComparator).CompareValue(indexKey.(bool))
}

var _ indexes.Storage = (*boolHashIndexStorage)(nil)

type boolHashIndexStorage struct {
	byValue map[bool]*storage.IDStorage
}

func (idx *boolHashIndexStorage) Get(key interface{}) *storage.IDStorage {
	return idx.byValue[key.(bool)]
}

func (idx *boolHashIndexStorage) Set(key interface{}, records *storage.IDStorage) {
	idx.byValue[key.(bool)] = records
}

func (idx *boolHashIndexStorage) Keys() []interface{} {
	i := 0
	keys := make([]interface{}, len(idx.byValue))
	for key := range idx.byValue {
		keys[i] = key
		i += 1
	}
	return keys
}

func NewBoolHashIndex(getter *record.BoolGetter) indexes.Index {
	return NewIndex(
		getter.Field,
		boolIndexComputation{getter: getter},
		indexes.CreateConcurrentStorage(&boolHashIndexStorage{
			byValue: make(map[bool]*storage.IDStorage),
		}),
	)
}
