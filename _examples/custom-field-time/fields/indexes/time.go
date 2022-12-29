package indexes

import (
	"github.com/shamcode/simd/_examples/custom-field-time/fields"
	"github.com/shamcode/simd/_examples/custom-field-time/fields/comparators"
	"github.com/shamcode/simd/indexes/bytype"
	"github.com/shamcode/simd/indexes/storage"
	"github.com/shamcode/simd/where"
	"time"
)

var _ bytype.IndexComputer = timeIndexComputation{}

type timeIndexComputation struct {
	getter *fields.TimeGetter
}

func (idx timeIndexComputation) ForItem(item interface{}) interface{} {
	return idx.getter.Get(item).Unix()
}

func (idx timeIndexComputation) ForComparatorAllValues(comparator where.FieldComparator, cb func(interface{})) {
	for _, item := range comparator.(comparators.TimeFieldComparator).Value {
		cb(item.Unix())
	}
}

func (idx timeIndexComputation) ForComparatorFirstValue(comparator where.FieldComparator) interface{} {
	return comparator.(comparators.TimeFieldComparator).Value[0].Unix()
}

func (idx timeIndexComputation) Compare(value interface{}, comparator where.FieldComparator) (bool, error) {
	return comparator.(comparators.TimeFieldComparator).CompareValue(time.Unix(value.(int64), 0))
}

var _ bytype.Storage = (*timeIndexStorage)(nil)

type timeIndexStorage struct {
	byValue map[int64]*storage.IDStorage
}

func (idx *timeIndexStorage) Get(key interface{}) *storage.IDStorage {
	return idx.byValue[key.(int64)]
}

func (idx *timeIndexStorage) Set(key interface{}, records *storage.IDStorage) {
	idx.byValue[key.(int64)] = records
}

func (idx *timeIndexStorage) Count(key interface{}) int {
	return idx.byValue[key.(int64)].Count()
}

func (idx *timeIndexStorage) Keys() []interface{} {
	i := 0
	keys := make([]interface{}, len(idx.byValue))
	for key := range idx.byValue {
		keys[i] = key
		i += 1
	}
	return keys
}

func NewTimeIndex(getter *fields.TimeGetter) *bytype.Index {
	return &bytype.Index{
		Field:   getter.Field,
		Compute: timeIndexComputation{getter: getter},
		Storage: bytype.WrapToThreadSafeStorage(&timeIndexStorage{
			byValue: make(map[int64]*storage.IDStorage),
		}),
	}
}
