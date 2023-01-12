package bytype

import (
	"github.com/shamcode/simd/indexes/storage"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

var _ IndexComputer = enum16IndexComputation{}

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

var _ Storage = (*enum16IndexStorage)(nil)

type enum16IndexStorage struct {
	byValue map[uint16]*storage.IDStorage
}

func (idx *enum16IndexStorage) Get(key interface{}) *storage.IDStorage {
	return idx.byValue[key.(uint16)]
}

func (idx *enum16IndexStorage) Set(key interface{}, records *storage.IDStorage) {
	idx.byValue[key.(uint16)] = records
}

func (idx *enum16IndexStorage) Count(key interface{}) int {
	return idx.byValue[key.(uint16)].Count()
}

func (idx *enum16IndexStorage) Keys() []interface{} {
	i := 0
	keys := make([]interface{}, len(idx.byValue))
	for key := range idx.byValue {
		keys[i] = key
		i += 1
	}
	return keys
}

func NewEnum16Index(getter *record.Enum16Getter) *Index {
	return &Index{
		Field:   getter.Field,
		Compute: enum16IndexComputation{getter: getter},
		Storage: WrapToThreadSafeStorage(&enum16IndexStorage{
			byValue: make(map[uint16]*storage.IDStorage),
		}),
	}
}
