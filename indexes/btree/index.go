package btree

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
)

var _ indexes.Index = (*index)(nil)

type BTree interface {
	indexes.Storage
	LessThan(key Key) (int, []storage.LockableIDStorage)
	LessOrEqual(key Key) (int, []storage.LockableIDStorage)
	GreaterThan(key Key) (int, []storage.LockableIDStorage)
	GreaterOrEqual(key Key) (int, []storage.LockableIDStorage)
	ForKey(key Key) (int, storage.LockableIDStorage)
	All(func(key Key, records storage.IDStorage))
}

type index struct {
	field   string
	unique  bool
	compute indexes.IndexComputer
	storage indexes.ConcurrentStorage
}

func (idx *index) BTree() BTree {
	return idx.storage.Unwrap().(BTree)
}

func (idx *index) Field() string {
	return idx.field
}

func (idx *index) Unique() bool {
	return idx.unique
}

func (idx *index) Compute() indexes.IndexComputer {
	return idx.compute
}

func (idx *index) Weight(condition where.Condition) (canApplyIndex bool, weight indexes.IndexWeight) {
	switch condition.Cmp.GetType() {
	case where.LT, where.LE, where.GT, where.GE:

		// B-tree optimal for <, <=, >, >=
		return true, indexes.IndexWeightLow
	case where.EQ, where.InArray:
		if condition.WithNot {
			return true, indexes.IndexWeightHigh
		} else {

			// Hash index more optimal for A == 1 and A in (1, 2, 3)
			return true, indexes.IndexWeightMedium
		}
	default:

		// For other condition index can apply, but not optimal
		return true, indexes.IndexWeightHigh
	}
}

func (idx *index) Select(condition where.Condition) (count int, ids []storage.LockableIDStorage, err error) {
	cmp := condition.Cmp.GetType()
	if condition.WithNot {
		switch cmp {
		case where.LT:
			cmp = where.GE
		case where.LE:
			cmp = where.GT
		case where.GT:
			cmp = where.LE
		case where.GE:
			cmp = where.LT
		}
	}
	switch cmp {
	case where.LT:
		idx.storage.RLock()
		count, ids = idx.BTree().LessThan(idx.compute.ForValue(condition.Cmp.ValueAt(0)).(Key))
		idx.storage.RUnlock()
		return
	case where.LE:
		idx.storage.RLock()
		count, ids = idx.BTree().LessOrEqual(idx.compute.ForValue(condition.Cmp.ValueAt(0)).(Key))
		idx.storage.RUnlock()
		return
	case where.GT:
		idx.storage.RLock()
		count, ids = idx.BTree().GreaterThan(idx.compute.ForValue(condition.Cmp.ValueAt(0)).(Key))
		idx.storage.RUnlock()
		return
	case where.GE:
		idx.storage.RLock()
		count, ids = idx.BTree().GreaterOrEqual(idx.compute.ForValue(condition.Cmp.ValueAt(0)).(Key))
		idx.storage.RUnlock()
		return
	}
	if !condition.WithNot {
		switch cmp {
		case where.EQ:
			idx.storage.RLock()
			countForKey, idsForKey := idx.BTree().ForKey(idx.compute.ForValue(condition.Cmp.ValueAt(0)).(Key))
			idx.storage.RUnlock()
			if countForKey > 0 {
				count = countForKey
				ids = []storage.LockableIDStorage{idsForKey}
			}
			return
		case where.InArray:
			idx.storage.RLock()
			for i := 0; i < condition.Cmp.ValuesCount(); i++ {
				countForValue, idsForValue := idx.BTree().ForKey(idx.compute.ForValue(condition.Cmp.ValueAt(i)).(Key))
				if countForValue > 0 {
					count += countForValue
					ids = append(ids, idsForValue)
				}
			}
			idx.storage.RUnlock()
			return
		}
	}
	idx.storage.RLock()
	idx.BTree().All(func(key Key, records storage.IDStorage) {
		resultForValue, errorForValue := idx.compute.Check(key, condition.Cmp)
		if nil != errorForValue {
			err = errorForValue
			return
		}
		if condition.WithNot != resultForValue {
			itemCount := records.Count()
			if itemCount > 0 {
				count += itemCount
				ids = append(ids, records)
			}
		}
	})
	idx.storage.RLock()
	return
}

func (idx *index) ConcurrentStorage() indexes.ConcurrentStorage {
	return idx.storage
}

func NewIndex(field string, compute indexes.IndexComputer, btree BTree, unique bool) indexes.Index {
	return &index{
		field:   field,
		unique:  unique,
		compute: compute,
		storage: indexes.CreateConcurrentStorage(btree, unique),
	}
}
