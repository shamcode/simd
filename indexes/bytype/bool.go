package bytype

import (
	"github.com/shamcode/simd/indexes/storage"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

var _ IndexComputer = boolIndexComputation{}

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

var _ Storage = (*boolIndexStorage)(nil)

type boolIndexStorage struct {
	byValue map[bool]*storage.IDStorage
}

func (idx *boolIndexStorage) Get(key interface{}) *storage.IDStorage {
	return idx.byValue[key.(bool)]
}

func (idx *boolIndexStorage) Set(key interface{}, records *storage.IDStorage) {
	idx.byValue[key.(bool)] = records
}

func (idx *boolIndexStorage) Count(key interface{}) int {
	return idx.byValue[key.(bool)].Count()
}

func (idx *boolIndexStorage) Keys() []interface{} {
	i := 0
	keys := make([]interface{}, len(idx.byValue))
	for key := range idx.byValue {
		keys[i] = key
		i += 1
	}
	return keys
}

func NewBoolIndex(getter *record.BoolGetter) *Index {
	return &Index{
		Field:   getter.Field,
		Compute: boolIndexComputation{getter: getter},
		Storage: WrapToThreadSafeStorage(&boolIndexStorage{
			byValue: make(map[bool]*storage.IDStorage),
		}),
	}
}
