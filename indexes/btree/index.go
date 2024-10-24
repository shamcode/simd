//nolint:exhaustive,nonamedreturns,funlen
package btree

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
)

var _ indexes.Index = (*index)(nil)

type Storage interface {
	indexes.Storage
	LessThan(key indexes.Key) (int, []storage.IDIterator)
	LessOrEqual(key indexes.Key) (int, []storage.IDIterator)
	GreaterThan(key indexes.Key) (int, []storage.IDIterator)
	GreaterOrEqual(key indexes.Key) (int, []storage.IDIterator)
	ForKey(key indexes.Key) (int, storage.IDIterator)
	All(callback func(key indexes.Key, records storage.IDStorage))
}

type index struct {
	field   record.Field
	unique  bool
	compute indexes.IndexComputer
	storage indexes.ConcurrentStorage
}

func (idx index) btree() Storage {
	return idx.storage.Unwrap().(Storage)
}

func (idx index) Field() record.Field {
	return idx.field
}

func (idx index) Unique() bool {
	return idx.unique
}

func (idx index) Compute() indexes.IndexComputer {
	return idx.compute
}

func (idx index) Weight(condition where.Condition) (canApplyIndex bool, weight indexes.IndexWeight) {
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

func (idx index) Select(condition where.Condition) ( //nolint:cyclop
	count int,
	ids []storage.IDIterator,
	err error,
) {
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
		count, ids = idx.btree().LessThan(idx.compute.ForValue(condition.Cmp.ValueAt(0)))
		idx.storage.RUnlock()
		return
	case where.LE:
		idx.storage.RLock()
		count, ids = idx.btree().LessOrEqual(idx.compute.ForValue(condition.Cmp.ValueAt(0)))
		idx.storage.RUnlock()
		return
	case where.GT:
		idx.storage.RLock()
		count, ids = idx.btree().GreaterThan(idx.compute.ForValue(condition.Cmp.ValueAt(0)))
		idx.storage.RUnlock()
		return
	case where.GE:
		idx.storage.RLock()
		count, ids = idx.btree().GreaterOrEqual(idx.compute.ForValue(condition.Cmp.ValueAt(0)))
		idx.storage.RUnlock()
		return
	}
	if !condition.WithNot {
		switch cmp {
		case where.EQ:
			idx.storage.RLock()
			countForKey, idsForKey := idx.btree().ForKey(idx.compute.ForValue(condition.Cmp.ValueAt(0)))
			idx.storage.RUnlock()
			if countForKey > 0 {
				count = countForKey
				ids = []storage.IDIterator{idsForKey}
			}
			return
		case where.InArray:
			idx.storage.RLock()
			for i := range condition.Cmp.ValuesCount() {
				countForValue, idsForValue := idx.btree().ForKey(idx.compute.ForValue(condition.Cmp.ValueAt(i)))
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
	idx.btree().All(func(key indexes.Key, records storage.IDStorage) {
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
	return //nolint:nakedret
}

func (idx index) ConcurrentStorage() indexes.ConcurrentStorage {
	return idx.storage
}

func NewIndex(field record.Field, compute indexes.IndexComputer, btree Storage, unique bool) indexes.Index {
	return index{
		field:   field,
		unique:  unique,
		compute: compute,
		storage: indexes.CreateConcurrentStorage(btree, unique),
	}
}
