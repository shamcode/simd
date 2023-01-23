package storage

import (
	"github.com/shamcode/simd/record"
	"sync"
)

type RecordsByID interface {
	GetIDStorage() LockableIDStorage
	Get(id int64) record.Record
	Set(id int64, item record.Record)
	Delete(id int64)
	Count() int
	GetData(stores []LockableIDStorage, totalCount int) []record.Record
	GetAllData() []record.Record
}

var _ RecordsByID = (*recordsByID)(nil)

type recordsByID struct {
	sync.RWMutex
	data map[int64]record.Record
	ids  *innerIDStorage
}

func (r *recordsByID) GetIDStorage() LockableIDStorage {
	return r.ids
}

func (r *recordsByID) Get(id int64) record.Record {
	r.RLock()
	defer r.RUnlock()
	return r.data[id]

}

func (r *recordsByID) Set(id int64, item record.Record) {
	r.Lock()
	r.data[id] = item
	r.ids.data[id] = struct{}{}
	r.Unlock()
}

func (r *recordsByID) Delete(id int64) {
	r.Lock()
	delete(r.data, id)
	delete(r.ids.data, id)
	r.Unlock()

}

func (r *recordsByID) Count() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.data)
}

func (r *recordsByID) GetData(stores []LockableIDStorage, totalCount int) []record.Record {
	if 1 == len(stores) {
		// Optimization for one store case
		return r.selectByStore(stores[0], totalCount)
	}
	return r.selectUniq(stores, totalCount)
}

func (r *recordsByID) selectByStore(store LockableIDStorage, totalCount int) []record.Record {
	items := make([]record.Record, 0, totalCount)
	r.RLock()
	store.RLock()
	for id := range store.ThreadUnsafeData() {
		items = append(items, r.data[id])
	}
	store.RUnlock()
	r.RUnlock()
	return items
}

func (r *recordsByID) selectUniq(stores []LockableIDStorage, totalCount int) []record.Record {
	var items []record.Record
	added := make(map[int64]struct{}, totalCount)
	r.RLock()
	for _, store := range stores {
		store.RLock()
		for id := range store.ThreadUnsafeData() {
			if _, ok := added[id]; !ok {
				items = append(items, r.data[id])
				added[id] = struct{}{}
			}
		}
		store.RUnlock()
	}
	r.RUnlock()
	return items
}

func (r *recordsByID) GetAllData() []record.Record {
	r.RLock()
	items := make([]record.Record, 0, len(r.data))
	for _, item := range r.data {
		items = append(items, item)
	}
	r.RUnlock()
	return items
}

func CreateRecordsByID() RecordsByID {
	return &recordsByID{
		ids:  createInnerIDStorage(),
		data: make(map[int64]record.Record),
	}
}
