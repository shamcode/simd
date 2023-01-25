package storage

import (
	"sync"
	"sync/atomic"
)

var _ IDStorage = (*multipleIDStorage)(nil)

type multipleIDStorage struct {
	sync.RWMutex
	data  map[int64]struct{}
	count int64
}

func (r *multipleIDStorage) ThreadUnsafeData() map[int64]struct{} {
	return r.data
}

func (r *multipleIDStorage) Count() int {
	return int(atomic.LoadInt64(&r.count))
}

func (r *multipleIDStorage) Add(id int64) {
	r.Lock()
	r.data[id] = struct{}{}
	atomic.StoreInt64(&r.count, int64(len(r.data)))
	r.Unlock()
}

func (r *multipleIDStorage) Delete(id int64) {
	r.Lock()
	delete(r.data, id)
	atomic.StoreInt64(&r.count, int64(len(r.data)))
	r.Unlock()
}

func CreateMultipleIDStorage() IDStorage {
	return &multipleIDStorage{
		data: make(map[int64]struct{}),
	}
}
