package hash

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
)

var _ indexes.Index = (*Index)(nil)

type Index struct {
	field   string
	compute indexes.IndexComputer
	storage indexes.Storage
}

func (index *Index) Field() string {
	return index.field
}

func (index *Index) Compute() indexes.IndexComputer {
	return index.compute
}

func (index *Index) Weight(condition where.Condition) (canApplyIndex bool, weight indexes.IndexWeight) {
	cmp := condition.Cmp.GetType()
	if !condition.WithNot && (where.EQ == cmp || where.InArray == cmp) {

		// Hash index optimal for A == 1 and A in (1, 2, 3)
		return true, indexes.IndexWeightLow
	}

	// For other condition index can apply, but not optimal
	return true, indexes.IndexWeightHigh
}

func (index *Index) Select(condition where.Condition) (count int, ids []storage.LockableIDStorage, err error) {
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

func (index *Index) selectForEqual(condition where.Condition) (count int, ids []storage.LockableIDStorage) {
	itemsByValue := index.storage.Get(index.compute.ForValue(condition.Cmp.ValueAt(0)))
	if nil != itemsByValue {
		count = itemsByValue.Count()
		ids = []storage.LockableIDStorage{itemsByValue}
	}
	return
}

func (index *Index) selectForInArray(condition where.Condition) (count int, ids []storage.LockableIDStorage) {
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

func (index *Index) selectForOther(condition where.Condition) (count int, ids []storage.LockableIDStorage, err error) {
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

func (index *Index) Storage() indexes.Storage {
	return index.storage
}

func NewIndex(field string, compute indexes.IndexComputer, storage indexes.Storage) indexes.Index {
	return &Index{
		field:   field,
		compute: compute,
		storage: storage,
	}
}
