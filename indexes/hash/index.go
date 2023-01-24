package hash

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
)

var _ indexes.Index = (*index)(nil)

type index struct {
	field   string
	compute indexes.IndexComputer
	storage indexes.ConcurrentStorage
}

func (idx *index) Field() string {
	return idx.field
}

func (idx *index) Compute() indexes.IndexComputer {
	return idx.compute
}

func (idx *index) Weight(condition where.Condition) (canApplyIndex bool, weight indexes.IndexWeight) {
	cmp := condition.Cmp.GetType()
	if !condition.WithNot && (where.EQ == cmp || where.InArray == cmp) {

		// Hash index optimal for A == 1 and A in (1, 2, 3)
		return true, indexes.IndexWeightLow
	}

	// For other condition index can apply, but not optimal
	return true, indexes.IndexWeightHigh
}

func (idx *index) Select(condition where.Condition) (count int, ids []storage.LockableIDStorage, err error) {
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

func (idx *index) selectForEqual(condition where.Condition) (count int, ids []storage.LockableIDStorage) {
	itemsByValue := idx.storage.Get(idx.compute.ForValue(condition.Cmp.ValueAt(0)))
	if nil != itemsByValue {
		count = itemsByValue.Count()
		ids = []storage.LockableIDStorage{itemsByValue}
	}
	return
}

func (idx *index) selectForInArray(condition where.Condition) (count int, ids []storage.LockableIDStorage) {
	for i := 0; i < condition.Cmp.ValuesCount(); i++ {
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

func (idx *index) selectForOther(condition where.Condition) (count int, ids []storage.LockableIDStorage, err error) {
	keys := idx.storage.Keys()
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

func (idx *index) Storage() indexes.ConcurrentStorage {
	return idx.storage
}

func NewIndex(field string, compute indexes.IndexComputer, storage indexes.ConcurrentStorage) indexes.Index {
	return &index{
		field:   field,
		compute: compute,
		storage: storage,
	}
}
