package indexes

import (
	"github.com/shamcode/simd/storage"
	"sync"
)

type Storage interface {
	Get(key interface{}) *storage.IDStorage
	Set(key interface{}, records *storage.IDStorage)
	Count(key interface{}) int
	Keys() []interface{}
}

var _ Storage = (*rwStorage)(nil)

// rwStorage is a thread safe wrapper for Storage
type rwStorage struct {
	sync.RWMutex
	original Storage
}

func (idx *rwStorage) Get(key interface{}) *storage.IDStorage {
	idx.RLock()
	defer idx.RUnlock()
	return idx.original.Get(key)
}

func (idx *rwStorage) Set(key interface{}, records *storage.IDStorage) {
	idx.Lock()
	idx.original.Set(key, records)
	idx.Unlock()
}

func (idx *rwStorage) Count(key interface{}) int {
	idx.RLock()
	defer idx.RUnlock()
	return idx.original.Count(key)
}

func (idx *rwStorage) Keys() []interface{} {
	idx.RLock()
	defer idx.RUnlock()
	return idx.original.Keys()
}

func WrapToThreadSafeStorage(original Storage) Storage {
	return &rwStorage{
		original: original,
	}
}
