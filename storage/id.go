package storage

import (
	"sync"
)

type LockableIDStorage interface {
	RLock()
	RUnlock()
	ThreadUnsafeData() map[int64]struct{}
}

type IDStorage struct {
	sync.RWMutex
	data map[int64]struct{}
}

func (r *IDStorage) ThreadUnsafeData() map[int64]struct{} {
	return r.data
}

func (r *IDStorage) Count() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.data)
}

func (r *IDStorage) Add(id int64) {
	r.Lock()
	r.data[id] = struct{}{}
	r.Unlock()
}

func (r *IDStorage) Delete(id int64) {
	r.Lock()
	delete(r.data, id)
	r.Unlock()
}

func NewIDStorage() *IDStorage {
	return &IDStorage{
		data: make(map[int64]struct{}),
	}
}
