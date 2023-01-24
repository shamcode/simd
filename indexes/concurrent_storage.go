package indexes

import (
	"github.com/shamcode/simd/storage"
	"sync"
)

var _ ConcurrentStorage = (*concurrentStorage)(nil)

type concurrentStorage struct {
	sync.RWMutex
	original Storage
}

func (idx *concurrentStorage) Get(key interface{}) *storage.IDStorage {
	idx.RLock()
	defer idx.RUnlock()
	return idx.original.Get(key)
}

func (idx *concurrentStorage) GetOrCreate(key interface{}) *storage.IDStorage {
	idx.RLock()
	idStorage := idx.original.Get(key)
	idx.RUnlock()
	if nil != idStorage {
		return idStorage
	}
	idx.Lock()
	idStorage = idx.original.Get(key)
	if nil == idStorage {

		// Prevent override in race
		idStorage = storage.NewIDStorage()
		idx.original.Set(key, idStorage)
	}
	idx.Unlock()
	return idStorage
}

func (idx *concurrentStorage) Keys() []interface{} {
	idx.RLock()
	defer idx.RUnlock()
	return idx.original.Keys()
}

func CreateConcurrentStorage(original Storage) ConcurrentStorage {
	return &concurrentStorage{
		original: original,
	}
}
