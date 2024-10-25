package storage

import (
	"sync"

	"github.com/shamcode/simd/record"
)

type recordsByID[R record.Record] struct {
	sync.RWMutex
	data map[int64]R
	ids  MapIDStorage
}

func (r *recordsByID[R]) GetIDStorage() IDIterator {
	return r.ids
}

func (r *recordsByID[R]) Get(id int64) (R, bool) {
	r.RLock()
	defer r.RUnlock()
	value, ok := r.data[id]
	return value, ok
}

func (r *recordsByID[R]) Set(id int64, item R) {
	r.Lock()
	r.data[id] = item
	r.ids[id] = struct{}{}
	r.Unlock()
}

func (r *recordsByID[R]) Delete(id int64) {
	r.Lock()
	delete(r.data, id)
	delete(r.ids, id)
	r.Unlock()
}

func (r *recordsByID[R]) Count() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.data)
}

func (r *recordsByID[R]) GetData(stores []IDIterator, totalCount int, idsUnique bool) []R {
	if idsUnique {
		return r.selectByUniqIDsStore(stores, totalCount)
	}
	if len(stores) == 1 {
		// Optimization for one store case
		return r.selectByStore(stores[0], totalCount)
	}
	return r.selectUniq(stores, totalCount)
}

func (r *recordsByID[R]) selectByStore(store IDIterator, totalCount int) []R {
	items := make([]R, totalCount)
	var i int
	r.RLock()
	store.Iterate(func(id int64) {
		items[i] = r.data[id]
		i++
	})
	r.RUnlock()
	return items
}

func (r *recordsByID[R]) selectByUniqIDsStore(stores []IDIterator, totalCount int) []R {
	items := make([]R, 0, totalCount)
	r.RLock()
	for _, store := range stores {
		if id := store.(UniqueIDStorage).ID(); id != 0 {
			items = append(items, r.data[id])
		}
	}
	r.RUnlock()
	return items
}

func (r *recordsByID[R]) selectUniq(stores []IDIterator, totalCount int) []R {
	var items []R
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

func (r *recordsByID[R]) GetAllData() []R {
	r.RLock()
	items := make([]R, len(r.data))
	var i int
	for _, item := range r.data {
		items[i] = item
		i++
	}
	r.RUnlock()
	return items
}

func CreateRecordsByID[R record.Record]() RecordsByID[R] {
	return &recordsByID[R]{ //nolint:exhaustruct
		ids:  CreateMapIDStorage(),
		data: make(map[int64]R),
	}
}
