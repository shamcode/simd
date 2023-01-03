package bytype

import (
	"github.com/shamcode/simd/indexes/storage"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

var _ IndexComputer = int32IndexComputation{}

type int32IndexComputation struct {
	getter *record.Int32Getter
}

func (idx int32IndexComputation) ForRecord(item record.Record) interface{} {
	return idx.getter.Get(item)
}

func (idx int32IndexComputation) EachComparatorValues(comparator where.FieldComparator, cb func(interface{})) {
	for _, item := range comparator.(comparators.Int32FieldComparator).Value {
		cb(item)
	}
}

func (idx int32IndexComputation) ForComparatorFirstValue(comparator where.FieldComparator) interface{} {
	return comparator.(comparators.Int32FieldComparator).Value[0]
}

func (idx int32IndexComputation) Compare(value interface{}, comparator where.FieldComparator) (bool, error) {
	return comparator.(comparators.Int32FieldComparator).CompareValue(value.(int32))
}

var _ Storage = (*int32IndexStorage)(nil)

type int32IndexStorage struct {
	byValue map[int32]*storage.IDStorage
}

func (idx *int32IndexStorage) Get(key interface{}) *storage.IDStorage {
	return idx.byValue[key.(int32)]
}

func (idx *int32IndexStorage) Set(key interface{}, records *storage.IDStorage) {
	idx.byValue[key.(int32)] = records
}

func (idx *int32IndexStorage) Count(key interface{}) int {
	return idx.byValue[key.(int32)].Count()
}

func (idx *int32IndexStorage) Keys() []interface{} {
	i := 0
	keys := make([]interface{}, len(idx.byValue))
	for key := range idx.byValue {
		keys[i] = key
		i += 1
	}
	return keys
}

func NewInt32Index(getter *record.Int32Getter) *Index {
	return &Index{
		Field:   getter.Field,
		Compute: int32IndexComputation{getter: getter},
		Storage: WrapToThreadSafeStorage(&int32IndexStorage{
			byValue: make(map[int32]*storage.IDStorage),
		}),
	}
}
