package fields

import (
	"github.com/shamcode/simd/indexes/storage"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

var _ IndexComputer = int64IndexComputation{}

type int64IndexComputation struct {
	getter *record.Int64Getter
}

func (idx int64IndexComputation) ForItem(item interface{}) interface{} {
	return idx.getter.Get(item)
}

func (idx int64IndexComputation) ForComparatorAllValues(comparator where.FieldComparator, cb func(interface{})) {
	for _, item := range comparator.(comparators.Int64FieldComparator).Value {
		cb(item)
	}
}

func (idx int64IndexComputation) ForComparatorFirstValue(comparator where.FieldComparator) interface{} {
	return comparator.(comparators.Int64FieldComparator).Value[0]
}

func (idx int64IndexComputation) Compare(value interface{}, comparator where.FieldComparator) (bool, error) {
	return comparator.(comparators.Int64FieldComparator).CompareValue(value.(int64))
}

var _ Storage = (*int64IndexStorage)(nil)

type int64IndexStorage struct {
	byValue map[int64]*storage.IDStorage
}

func (idx *int64IndexStorage) Get(key interface{}) *storage.IDStorage {
	return idx.byValue[key.(int64)]
}

func (idx *int64IndexStorage) Set(key interface{}, records *storage.IDStorage) {
	idx.byValue[key.(int64)] = records
}

func (idx *int64IndexStorage) Count(key interface{}) int {
	return idx.byValue[key.(int64)].Count()
}

func (idx *int64IndexStorage) Keys() []interface{} {
	i := 0
	keys := make([]interface{}, len(idx.byValue))
	for key := range idx.byValue {
		keys[i] = key
		i += 1
	}
	return keys
}

func NewInt64Index(getter *record.Int64Getter) *Index {
	return &Index{
		Field:   getter.Field,
		Compute: int64IndexComputation{getter: getter},
		Storage: WrapToThreadSafeStorage(&int64IndexStorage{
			byValue: make(map[int64]*storage.IDStorage),
		}),
	}
}
