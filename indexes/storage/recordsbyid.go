package storage

import (
	"github.com/shamcode/simd/record"
	"sync"
)

type RecordsByID struct {
	sync.RWMutex
	data map[int64]record.Record
	ids  *innerIDStorage
}

func (r *RecordsByID) GetIDStorage() LockableIDStorage {
	return r.ids
}

func (r *RecordsByID) Get(id int64) record.Record {
	r.RLock()
	defer r.RUnlock()
	return r.data[id]

}

func (r *RecordsByID) Set(id int64, item record.Record) {
	r.Lock()
	r.data[id] = item
	r.ids.data[id] = struct{}{}
	r.Unlock()
}

func (r *RecordsByID) Delete(id int64) {
	r.Lock()
	delete(r.data, id)
	delete(r.ids.data, id)
	r.Unlock()

}

func (r *RecordsByID) Count() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.data)
}

func (r *RecordsByID) GetData(stores []LockableIDStorage, totalCount int) []record.Record {
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

func (r *RecordsByID) GetAllData() []record.Record {
	r.RLock()
	items := make([]record.Record, 0, len(r.data))
	for _, item := range r.data {
		items = append(items, item)
	}
	r.RUnlock()
	return items
}

func NewRecordsByID() *RecordsByID {
	return &RecordsByID{
		ids:  newInnerIDStorage(),
		data: make(map[int64]record.Record),
	}
}
