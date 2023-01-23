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
	storage indexes.Storage
}

func (index *index) Field() string {
	return index.field
}

func (index *index) Compute() indexes.IndexComputer {
	return index.compute
}

func (index *index) Weight(condition where.Condition) (canApplyIndex bool, weight indexes.IndexWeight) {
	cmp := condition.Cmp.GetType()
	if !condition.WithNot && (where.EQ == cmp || where.InArray == cmp) {

		// Hash index optimal for A == 1 and A in (1, 2, 3)
		return true, indexes.IndexWeightLow
	}

	// For other condition index can apply, but not optimal
	return true, indexes.IndexWeightHigh
}

func (index *index) Select(condition where.Condition) (count int, ids []storage.LockableIDStorage, err error) {
	if !condition.WithNot {
		switch condition.Cmp.GetType() {
		case where.EQ:
			count, ids = index.selectForEqual(condition)
			return
		case where.InArray:
			count, ids = index.selectForInArray(condition)
			return
		}
	}
	return index.selectForOther(condition)
}

func (index *index) selectForEqual(condition where.Condition) (count int, ids []storage.LockableIDStorage) {
	itemsByValue := index.storage.Get(index.compute.ForValue(condition.Cmp.ValueAt(0)))
	if nil != itemsByValue {
		count = itemsByValue.Count()
		ids = []storage.LockableIDStorage{itemsByValue}
	}
	return
}

func (index *index) selectForInArray(condition where.Condition) (count int, ids []storage.LockableIDStorage) {
	for i := 0; i < condition.Cmp.ValuesCount(); i++ {
		itemsByValue := index.storage.Get(index.compute.ForValue(condition.Cmp.ValueAt(i)))
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

func (index *index) selectForOther(condition where.Condition) (count int, ids []storage.LockableIDStorage, err error) {
	keys := index.storage.Keys()
	for _, key := range keys {
		resultForValue, errorForValue := index.compute.Check(key, condition.Cmp)
		if nil != errorForValue {
			err = errorForValue
			return
		}
		if condition.WithNot != resultForValue {
			idsForKey := index.storage.Get(key)
			count += idsForKey.Count()
			ids = append(ids, idsForKey)
		}
	}
	return
}

func (index *index) Storage() indexes.Storage {
	return index.storage
}

func NewIndex(field string, compute indexes.IndexComputer, storage indexes.Storage) indexes.Index {
	return &index{
		field:   field,
		compute: compute,
		storage: storage,
	}
}
