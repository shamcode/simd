package fields

import (
	"github.com/shamcode/simd/indexes/storage"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/fields"
)

var _ IndexComputer = (*stringIndexComputation)(nil)

type stringIndexComputation struct {
	getter *record.StringGetter
}

func (idx *stringIndexComputation) ForItem(item interface{}) interface{} {
	return idx.getter.Get(item)
}

func (idx *stringIndexComputation) ForComparatorAllValues(comparator where.FieldComparator, cb func(interface{})) {
	for _, item := range comparator.(*fields.StringFieldComparator).Value {
		cb(item)
	}
}

func (idx *stringIndexComputation) ForComparatorFirstValue(comparator where.FieldComparator) interface{} {
	return comparator.(*fields.StringFieldComparator).Value[0]
}

func (idx *stringIndexComputation) Compare(value interface{}, comparator where.FieldComparator) bool {
	return comparator.(*fields.StringFieldComparator).CompareValue(value.(string))
}

var _ Storage = (*stringIndexStorage)(nil)

type stringIndexStorage struct {
	byValue map[string]*storage.IDStorage
}

func (idx *stringIndexStorage) Get(key interface{}) *storage.IDStorage {
	return idx.byValue[key.(string)]
}

func (idx *stringIndexStorage) Set(key interface{}, records *storage.IDStorage) {
	idx.byValue[key.(string)] = records
}

func (idx *stringIndexStorage) Count(key interface{}) int {
	return idx.byValue[key.(string)].Count()
}

func (idx *stringIndexStorage) Keys() []interface{} {
	i := 0
	keys := make([]interface{}, len(idx.byValue))
	for key := range idx.byValue {
		keys[i] = key
		i += 1
	}
	return keys
}

func NewStringIndex(getter *record.StringGetter) *Index {
	return &Index{
		Field:   getter.Field,
		Compute: &stringIndexComputation{getter: getter},
		Storage: WrapToThreadSafeStorage(&stringIndexStorage{
			byValue: make(map[string]*storage.IDStorage),
		}),
	}
}
