package indexes

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
)

type ByField interface {
	Add(index Index)
	Insert(item record.Record)
	Delete(item record.Record)
	Update(oldItem, item record.Record)
	SelectForCondition(condition where.Condition) (
		indexExists bool,
		count int,
		ids []storage.IDIterator,
		idsUnique bool,
		err error,
	)
}

var _ ByField = byField{}

type byField map[uint8][]Index

func (ibf byField) Add(index Index) {
	i := index.Field().Index()
	ibf[i] = append(ibf[i], index)
}

func (ibf byField) Insert(item record.Record) {
	for _, indexesForField := range ibf {
		for _, idx := range indexesForField {
			key := idx.Compute().ForRecord(item)
			idx.ConcurrentStorage().GetOrCreate(key).Add(item.GetID())
		}
	}
}

func (ibf byField) Delete(item record.Record) {
	for _, indexesForField := range ibf {
		for _, idx := range indexesForField {
			records := idx.ConcurrentStorage().Get(idx.Compute().ForRecord(item))
			if nil != records {
				records.Delete(item.GetID())
			}
		}
	}
}

func (ibf byField) Update(oldItem, item record.Record) {
	for _, indexesForField := range ibf {
		for _, idx := range indexesForField {
			oldValue := idx.Compute().ForRecord(oldItem)
			newValue := idx.Compute().ForRecord(item)

			// TODO: if key is pointer, then compare invalid. Need add check for optional interface{ Equal(key interface{}} bool }
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

func (ibf byField) SelectForCondition(condition where.Condition) (
	indexExists bool,
	count int,
	ids []storage.IDIterator,
	idsUnique bool,
	err error,
) {
	var indexes []Index
	indexes, indexExists = ibf[condition.Cmp.GetField().Index()]
	if !indexExists || 0 == len(indexes) {
		return
	}
	first := true
	var minWeight IndexWeight
	var indexForApply Index
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
	return
}

func CreateByField() ByField {
	return make(byField)
}
