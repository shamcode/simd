package storage

import (
	"sync"

	"github.com/shamcode/simd/record"
)

var _ RecordsByID = (*recordsByID)(nil)

type recordsByID struct {
	sync.RWMutex
	data map[int64]record.Record
	ids  MapIDStorage
}

func (r *recordsByID) GetIDStorage() IDIterator {
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
	r.ids[id] = struct{}{}
	r.Unlock()
}

func (r *recordsByID) Delete(id int64) {
	r.Lock()
	delete(r.data, id)
	delete(r.ids, id)
	r.Unlock()
}

func (r *recordsByID) Count() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.data)
}

func (r *recordsByID) GetData(stores []IDIterator, totalCount int, idsUnique bool) []record.Record {
	if idsUnique {
		return r.selectByUniqIDsStore(stores, totalCount)
	}
	if len(stores) == 1 {
		// Optimization for one store case
		return r.selectByStore(stores[0], totalCount)
	}
	return r.selectUniq(stores, totalCount)
}

func (r *recordsByID) selectByStore(store IDIterator, totalCount int) []record.Record {
	items := make([]record.Record, totalCount)
	var i int
	r.RLock()
	store.Iterate(func(id int64) {
		items[i] = r.data[id]
		i++
	})
	r.RUnlock()
	return items
}

func (r *recordsByID) selectByUniqIDsStore(stores []IDIterator, totalCount int) []record.Record {
	items := make([]record.Record, 0, totalCount)
	r.RLock()
	for _, store := range stores {
		if id := store.(UniqueIDStorage).ID(); id != 0 {
			items = append(items, r.data[id])
		}
	}
	r.RUnlock()
	return items
}

func (r *recordsByID) selectUniq(stores []IDIterator, totalCount int) []record.Record {
	var items []record.Record
	added := make(map[int64]struct{}, totalCount)
	r.RLock()
	for _, store := range stores {
		store.Iterate(func(id int64) {
			if _, ok := added[id]; !ok {
				items = append(items, r.data[id])
				added[id] = struct{}{}
			}
		})
	}
	r.RUnlock()
	return items
}

func (r *recordsByID) GetAllData() []record.Record {
	r.RLock()
	items := make([]record.Record, len(r.data))
	var i int
	for _, item := range r.data {
		items[i] = item
		i++
	}
	r.RUnlock()
	return items
}

func CreateRecordsByID() RecordsByID {
	return &recordsByID{ //nolint:exhaustruct
		ids:  CreateMapIDStorage(),
		data: make(map[int64]record.Record),
	}
}
