package hash

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

var _ indexes.IndexComputer = int64IndexComputation{}

type int64IndexComputation struct {
	getter *record.Int64Getter
}

func (idx int64IndexComputation) ForRecord(item record.Record) interface{} {
	return idx.getter.Get(item)
}

func (idx int64IndexComputation) ForValue(value interface{}) interface{} {
	return value.(int64)
}

func (idx int64IndexComputation) Check(indexKey interface{}, comparator where.FieldComparator) (bool, error) {
	return comparator.(comparators.Int64FieldComparator).CompareValue(indexKey.(int64))
}

var _ indexes.Storage = (*int64HashIndexStorage)(nil)

type int64HashIndexStorage struct {
	byValue map[int64]*storage.IDStorage
}

func (idx *int64HashIndexStorage) Get(key interface{}) *storage.IDStorage {
	return idx.byValue[key.(int64)]
}

func (idx *int64HashIndexStorage) Set(key interface{}, records *storage.IDStorage) {
	idx.byValue[key.(int64)] = records
}

func (idx *int64HashIndexStorage) Keys() []interface{} {
	i := 0
	keys := make([]interface{}, len(idx.byValue))
	for key := range idx.byValue {
		keys[i] = key
		i += 1
	}
	return keys
}

func NewInt64HashIndex(getter *record.Int64Getter) indexes.Index {
	return NewIndex(
		getter.Field,
		int64IndexComputation{getter: getter},
		indexes.WrapToThreadSafeStorage(&int64HashIndexStorage{
			byValue: make(map[int64]*storage.IDStorage),
		}),
	)
}
