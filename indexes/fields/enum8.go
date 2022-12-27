package fields

import (
	"github.com/shamcode/simd/indexes/storage"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

var _ IndexComputer = (*enum8IndexComputation)(nil)

type enum8IndexComputation struct {
	getter *record.Enum8Getter
}

func (idx *enum8IndexComputation) ForItem(item interface{}) interface{} {
	return idx.getter.Get(item).Value()
}

func (idx *enum8IndexComputation) ForComparatorAllValues(comparator where.FieldComparator, cb func(interface{})) {
	for _, item := range comparator.(*comparators.Enum8FieldComparator).Value {
		cb(item.Value())
	}
}

func (idx *enum8IndexComputation) ForComparatorFirstValue(comparator where.FieldComparator) interface{} {
	return comparator.(*comparators.Enum8FieldComparator).Value[0].Value()
}

func (idx *enum8IndexComputation) Compare(value interface{}, comparator where.FieldComparator) bool {
	return comparator.(*comparators.Enum8FieldComparator).CompareValue(value.(uint8))
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
		Compute: &enum8IndexComputation{getter: getter},
		Storage: WrapToThreadSafeStorage(&enum8IndexStorage{
			byValue: make(map[uint8]*storage.IDStorage),
		}),
	}
}
