//nolint:exhaustive,nonamedreturns
package hash

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
)

type Storage interface {
	indexes.Storage
	Keys() []indexes.Key
}

type index[R record.Record] struct {
	field   record.Field
	unique  bool
	compute indexes.IndexComputer[R]
	storage indexes.ConcurrentStorage
}

func (idx index[R]) hashTable() Storage {
	return idx.storage.Unwrap().(Storage)
}

func (idx index[R]) Field() record.Field {
	return idx.field
}

func (idx index[R]) Unique() bool {
	return idx.unique
}

func (idx index[R]) Compute() indexes.IndexComputer[R] {
	return idx.compute
}

func (idx index[R]) Weight(condition where.Condition[R]) (canApplyIndex bool, weight indexes.IndexWeight) {
	cmp := condition.Cmp.GetType()
	if !condition.WithNot && (where.EQ == cmp || where.InArray == cmp) {
		// Hash index optimal for A == 1 and A in (1, 2, 3)
		return true, indexes.IndexWeightLow
	}

	// For other condition index can apply, but not optimal
	return true, indexes.IndexWeightHigh
}

func (idx index[R]) Select(condition where.Condition[R]) (count int, ids []storage.IDIterator, err error) {
	if !condition.WithNot {
		switch condition.Cmp.GetType() {
		case where.EQ:
			count, ids = idx.selectForEqual(condition)
			return
		case where.InArray:
			count, ids = idx.selectForInArray(condition)
			return
		}
	}
	return idx.selectForOther(condition)
}

func (idx index[R]) selectForEqual(condition where.Condition[R]) (count int, ids []storage.IDIterator) {
	itemsByValue := idx.storage.Get(idx.compute.ForValue(condition.Cmp.ValueAt(0)))
	if nil != itemsByValue {
		count = itemsByValue.Count()
		ids = []storage.IDIterator{itemsByValue}
	}
	return
}

func (idx index[R]) selectForInArray(condition where.Condition[R]) (count int, ids []storage.IDIterator) {
	for i := range condition.Cmp.ValuesCount() {
		itemsByValue := idx.storage.Get(idx.compute.ForValue(condition.Cmp.ValueAt(i)))
		if nil != itemsByValue {
			countForValue := itemsByValue.Count()
			if countForValue > 0 {
				count += countForValue
				ids = append(ids, itemsByValue)
			}
		}
	}
	return
}

func (idx index[R]) selectForOther(condition where.Condition[R]) (count int, ids []storage.IDIterator, err error) {
	idx.storage.RLock()
	keys := idx.hashTable().Keys()
	idx.storage.RUnlock()
	for _, key := range keys {
		resultForValue, errorForValue := idx.compute.Check(key, condition.Cmp)
		if nil != errorForValue {
			err = errorForValue
			return
		}
		if condition.WithNot != resultForValue {
			idsForKey := idx.storage.Get(key)
			count += idsForKey.Count()
			ids = append(ids, idsForKey)
		}
	}
	return
}

func (idx index[R]) ConcurrentStorage() indexes.ConcurrentStorage {
	return idx.storage
}

func NewIndex[R record.Record](
	field record.Field,
	compute indexes.IndexComputer[R],
	hashTable Storage,
	unique bool,
) indexes.Index[R] {
	return index[R]{
		field:   field,
		unique:  unique,
		compute: compute,
		storage: indexes.CreateConcurrentStorage(hashTable, unique),
	}
}
