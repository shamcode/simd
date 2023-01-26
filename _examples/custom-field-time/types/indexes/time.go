package indexes

import (
	"github.com/shamcode/simd/_examples/custom-field-time/types"
	"github.com/shamcode/simd/_examples/custom-field-time/types/comparators"
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/indexes/hash"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
	"time"
)

var _ indexes.IndexComputer = timeIndexComputation{}

type timeIndexComputation struct {
	getter *types.TimeGetter
}

func (idx timeIndexComputation) ForRecord(item record.Record) interface{} {
	return idx.getter.Get(item).UnixNano()
}

func (idx timeIndexComputation) ForValue(item interface{}) interface{} {
	return item.(time.Time).UnixNano()
}

func (idx timeIndexComputation) Check(indexKey interface{}, comparator where.FieldComparator) (bool, error) {
	return comparator.(comparators.TimeFieldComparator).CompareValue(time.Unix(0, indexKey.(int64)))
}

var _ hash.HashTable = (*timeHashIndexStorage)(nil)

type timeHashIndexStorage struct {
	byValue map[int64]storage.IDStorage
}

func (idx *timeHashIndexStorage) Get(key interface{}) storage.IDStorage {
	return idx.byValue[key.(int64)]
}

func (idx *timeHashIndexStorage) Set(key interface{}, records storage.IDStorage) {
	idx.byValue[key.(int64)] = records
}

func (idx *timeHashIndexStorage) Keys() []interface{} {
	keys := make([]interface{}, 0, len(idx.byValue))
	for key := range idx.byValue {
		keys = append(keys, key)
	}
	return keys
}

func NewTimeHashIndex(getter *types.TimeGetter, unique bool) indexes.Index {
	return hash.NewIndex(
		getter.Field,
		timeIndexComputation{getter: getter},
		&timeHashIndexStorage{
			byValue: make(map[int64]storage.IDStorage),
		},
		unique,
	)
}
