package storage

import (
	"sync"
	"sync/atomic"
)

type LockableIDStorage interface {
	RLock()
	RUnlock()
	ThreadUnsafeData() map[int64]struct{}
}

type IDStorage struct {
	sync.RWMutex
	data  map[int64]struct{}
	count int64
}

func (r *IDStorage) ThreadUnsafeData() map[int64]struct{} {
	return r.data
}

func (r *IDStorage) Count() int {
	return int(atomic.LoadInt64(&r.count))
}

func (r *IDStorage) Add(id int64) {
	r.Lock()
	r.data[id] = struct{}{}
	atomic.StoreInt64(&r.count, int64(len(r.data)))
	r.Unlock()
}

func (r *IDStorage) Delete(id int64) {
	r.Lock()
	delete(r.data, id)
	atomic.StoreInt64(&r.count, int64(len(r.data)))
	r.Unlock()
}

func NewIDStorage() *IDStorage {
	return &IDStorage{
		data: make(map[int64]struct{}),
	}
}
