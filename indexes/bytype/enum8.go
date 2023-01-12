package bytype

import (
	"github.com/shamcode/simd/indexes/storage"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

var _ IndexComputer = enum8IndexComputation{}

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

var _ Storage = (*enum8IndexStorage)(nil)

type enum8IndexStorage struct {
	byValue map[uint8]*storage.IDStorage
}

func (idx *enum8IndexStorage) Get(key interface{}) *storage.IDStorage {
	return idx.byValue[key.(uint8)]
}

func (idx *enum8IndexStorage) Set(key interface{}, records *storage.IDStorage) {
	idx.byValue[key.(uint8)] = records
}

func (idx *enum8IndexStorage) Count(key interface{}) int {
	return idx.byValue[key.(uint8)].Count()
}

func (idx *enum8IndexStorage) Keys() []interface{} {
	i := 0
	keys := make([]interface{}, len(idx.byValue))
	for key := range idx.byValue {
		keys[i] = key
		i += 1
	}
	return keys
}

func NewEnum8Index(getter *record.Enum8Getter) *Index {
	return &Index{
		Field:   getter.Field,
		Compute: enum8IndexComputation{getter: getter},
		Storage: WrapToThreadSafeStorage(&enum8IndexStorage{
			byValue: make(map[uint8]*storage.IDStorage),
		}),
	}
}
