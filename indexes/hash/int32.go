package hash

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

var _ indexes.IndexComputer = int32IndexComputation{}

type int32IndexComputation struct {
	getter *record.Int32Getter
}

func (idx int32IndexComputation) ForRecord(item record.Record) interface{} {
	return idx.getter.Get(item)
}

func (idx int32IndexComputation) ForValue(value interface{}) interface{} {
	return value.(int32)
}

func (idx int32IndexComputation) Check(indexKey interface{}, comparator where.FieldComparator) (bool, error) {
	return comparator.(comparators.Int32FieldComparator).CompareValue(indexKey.(int32))
}

var _ indexes.Storage = (*int32HashIndexStorage)(nil)

type int32HashIndexStorage struct {
	byValue map[int32]*storage.IDStorage
}

func (idx *int32HashIndexStorage) Get(key interface{}) *storage.IDStorage {
	return idx.byValue[key.(int32)]
}

func (idx *int32HashIndexStorage) Set(key interface{}, records *storage.IDStorage) {
	idx.byValue[key.(int32)] = records
}

func (idx *int32HashIndexStorage) Keys() []interface{} {
	i := 0
	keys := make([]interface{}, len(idx.byValue))
	for key := range idx.byValue {
		keys[i] = key
		i += 1
	}
	return keys
}

func NewInt32HashIndex(getter *record.Int32Getter) indexes.Index {
	return NewIndex(
		getter.Field,
		int32IndexComputation{getter: getter},
		indexes.CreateConcurrentStorage(&int32HashIndexStorage{
			byValue: make(map[int32]*storage.IDStorage),
		}),
	)
}
