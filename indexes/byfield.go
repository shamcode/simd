package indexes

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
)

type ByField interface {
	Add(index *Index)
	Insert(item record.Record)
	Delete(item record.Record)
	Update(oldItem, item record.Record)
	SelectForCondition(condition where.Condition) (
		indexExists bool,
		count int,
		ids []storage.LockableIDStorage,
		err error,
	)
}

var _ ByField = byField{}

type byField map[string][]*Index

func (ibf byField) Add(index *Index) {
	ibf[index.Field] = append(ibf[index.Field], index)
}

func (ibf byField) Insert(item record.Record) {
	for _, indexesForField := range ibf {
		for _, idx := range indexesForField {
			key := idx.Compute.ForRecord(item)
			records := idx.Storage.Get(key)
			if nil == records {
				records = storage.NewIDStorage()
				idx.Storage.Set(key, records)
			}
			records.Add(item.GetID())
		}
	}
}

func (ibf byField) Delete(item record.Record) {
	for _, indexesForField := range ibf {
		for _, idx := range indexesForField {
			records := idx.Storage.Get(idx.Compute.ForRecord(item))
			if nil != records {
				records.Delete(item.GetID())
			}
		}
	}

}

func (ibf byField) Update(oldItem, item record.Record) {
	for _, indexesForField := range ibf {
		for _, idx := range indexesForField {

			oldValue := idx.Compute.ForRecord(oldItem)
			newValue := idx.Compute.ForRecord(item)

			// TODO: if key is pointer, then compare invalid. Need add check for optional interface{ Equal(key interface{}} bool }
			if newValue == oldValue {

				// Field index not changed, ignore
				continue
			}

			// Remove old item from index
			oldRecords := idx.Storage.Get(oldValue)
			if nil != oldRecords {
				oldRecords.Delete(item.GetID())
			}

			records := idx.Storage.Get(newValue)
			if nil == records {

				// It's first item in index, create index storage
				records = storage.NewIDStorage()
				idx.Storage.Set(newValue, records)
			}

			// Add new item to index
			records.Add(item.GetID())
		}
	}
}

func (ibf byField) SelectForCondition(condition where.Condition) (
	indexExists bool,
	count int,
	ids []storage.LockableIDStorage,
	err error,
) {
	var indexes []*Index
	indexes, indexExists = ibf[condition.Cmp.GetField()]
	if !indexExists || 0 == len(indexes) {
		return
	}
	first := true
	for _, index := range indexes {
		countForIndex, idsForIndex, errForIndex := ibf.selectFromIndex(index, condition)
		if nil != errForIndex {
			err = errForIndex
			return
		}
		if first || countForIndex < count {
			first = false
			count = countForIndex
			ids = idsForIndex
		}
	}
	return
}

func (ibf byField) selectFromIndex(index *Index, condition where.Condition) (count int, ids []storage.LockableIDStorage, err error) {
	if !condition.WithNot {
		switch condition.Cmp.GetType() {
		case where.EQ: // A == '1'
			count, ids = ibf.selectForEqual(index, condition)
			return
		case where.InArray: // A IN ('1', '2', '3')
			count, ids = ibf.selectForInArray(index, condition)
			return
		}
	}
	count, ids, err = ibf.selectForOther(index, condition)
	return
}

func (ibf byField) selectForEqual(index *Index, condition where.Condition) (count int, ids []storage.LockableIDStorage) {
	itemsByValue := index.Storage.Get(index.Compute.ForValue(condition.Cmp.ValueAt(0)))
	if nil != itemsByValue {
		count = itemsByValue.Count()
		ids = append(ids, itemsByValue)
	}
	return
}

func (ibf byField) selectForInArray(index *Index, condition where.Condition) (count int, ids []storage.LockableIDStorage) {
	for i := 0; i < condition.Cmp.ValuesCount(); i++ {
		itemsByValue := index.Storage.Get(index.Compute.ForValue(condition.Cmp.ValueAt(i)))
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

func (ibf byField) selectForOther(index *Index, condition where.Condition) (count int, ids []storage.LockableIDStorage, err error) {
	keys := index.Storage.Keys()
	for _, key := range keys {
		resultForValue, errorForValue := index.Compute.Check(key, condition.Cmp)
		if nil != errorForValue {
			err = errorForValue
			return
		}
		if condition.WithNot != resultForValue {
			count += index.Storage.Count(key)
			ids = append(ids, index.Storage.Get(key))
		}
	}
	return
}

func NewByField() ByField {
	return make(byField)
}
