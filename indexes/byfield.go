package indexes

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
)

type ByField[R record.Record] interface {
	Add(index Index[R])
	Insert(item R)
	Delete(item R)
	Update(oldItem, item R)
	SelectForCondition(condition where.Condition[R]) (
		indexExists bool,
		count int,
		ids []storage.IDIterator,
		idsUnique bool,
		err error,
	)
}

type byField[R record.Record] map[uint8][]Index[R]

func (ibf byField[R]) Add(index Index[R]) {
	i := index.Field().Index()
	ibf[i] = append(ibf[i], index)
}

func (ibf byField[R]) Insert(item R) {
	for _, indexesForField := range ibf {
		for _, idx := range indexesForField {
			key := idx.Compute().ForRecord(item)
			idx.ConcurrentStorage().GetOrCreate(key).Add(item.GetID())
		}
	}
}

func (ibf byField[R]) Delete(item R) {
	for _, indexesForField := range ibf {
		for _, idx := range indexesForField {
			records := idx.ConcurrentStorage().Get(idx.Compute().ForRecord(item))
			if nil != records {
				records.Delete(item.GetID())
			}
		}
	}
}

func (ibf byField[R]) Update(oldItem, item R) {
	for _, indexesForField := range ibf {
		for _, idx := range indexesForField {
			oldValue := idx.Compute().ForRecord(oldItem)
			newValue := idx.Compute().ForRecord(item)

			//nolint:godox
			// TODO: if key is pointer, then compare invalid.
			// Need add check for optional interface{ Equal(key interface{}} bool }
			if newValue == oldValue {
				// Field index not changed, ignore
				continue
			}

			// Remove old item from index
			oldRecords := idx.ConcurrentStorage().Get(oldValue)
			if nil != oldRecords {
				oldRecords.Delete(item.GetID())
			}

			// Add new item to index
			idx.ConcurrentStorage().GetOrCreate(newValue).Add(item.GetID())
		}
	}
}

func (ibf byField[R]) SelectForCondition(condition where.Condition[R]) ( //nolint:nonamedreturns
	indexExists bool,
	count int,
	ids []storage.IDIterator,
	idsUnique bool,
	err error,
) {
	var indexes []Index[R]

	indexes, indexExists = ibf[condition.Cmp.GetField().Index()]
	if !indexExists || len(indexes) == 0 {
		return
	}

	first := true

	var (
		minWeight     IndexWeight
		indexForApply Index[R]
	)

	for _, index := range indexes {
		canApplyIndex, weight := index.Weight(condition)
		if !canApplyIndex {
			continue
		}

		if first || weight < minWeight {
			first = false
			minWeight = weight
			indexForApply = index
		}
	}

	if first {
		indexExists = false
	} else {
		count, ids, err = indexForApply.Select(condition)
		idsUnique = indexForApply.Unique()
	}

	return //nolint:nakedret
}

func CreateByField[R record.Record]() ByField[R] {
	return make(byField[R])
}
