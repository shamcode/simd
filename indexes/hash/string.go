package hash

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
)

type stringComparator interface {
	CompareValue(value string) (bool, error)
}

var _ indexes.IndexComputer = stringIndexComputation{}

type stringIndexComputation struct {
	getter *record.StringGetter
}

func (idx stringIndexComputation) ForRecord(item record.Record) interface{} {
	return idx.getter.Get(item)
}

func (idx stringIndexComputation) ForValue(value interface{}) interface{} {
	return value.(string)
}

func (idx stringIndexComputation) Check(indexKey interface{}, comparator where.FieldComparator) (bool, error) {
	return comparator.(stringComparator).CompareValue(indexKey.(string))
}

var _ HashTable = (*stringHashIndexStorage)(nil)

type stringHashIndexStorage struct {
	byValue map[string]storage.IDStorage
}

func (idx *stringHashIndexStorage) Get(key interface{}) storage.IDStorage {
	return idx.byValue[key.(string)]
}

func (idx *stringHashIndexStorage) Set(key interface{}, records storage.IDStorage) {
	idx.byValue[key.(string)] = records
}

func (idx *stringHashIndexStorage) Keys() []interface{} {
	i := 0
	keys := make([]interface{}, len(idx.byValue))
	for key := range idx.byValue {
		keys[i] = key
		i += 1
	}
	return keys
}

func NewStringHashIndex(getter *record.StringGetter, unique bool) indexes.Index {
	return NewIndex(
		getter.Field,
		stringIndexComputation{getter: getter},
		&stringHashIndexStorage{
			byValue: make(map[string]storage.IDStorage),
		},
		unique,
	)
}
