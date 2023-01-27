package indexes

import (
	"github.com/shamcode/simd/storage"
	"sync"
)

var _ ConcurrentStorage = (*concurrentStorage)(nil)

type concurrentStorage struct {
	sync.RWMutex
	original Storage
	uniq     bool
}

func (idx *concurrentStorage) Get(key interface{}) storage.IDStorage {
	idx.RLock()
	defer idx.RUnlock()
	return idx.original.Get(key)
}

func (idx *concurrentStorage) GetOrCreate(key interface{}) storage.IDStorage {
	idx.RLock()
	idStorage := idx.original.Get(key)
	idx.RUnlock()
	if nil != idStorage {
		return idStorage
	}
	idx.Lock()
	idStorage = idx.original.Get(key)
	if nil == idStorage { // Prevent override in race
		if idx.uniq {
			idStorage = storage.CreateUniqueIDStorage()
		} else {
			idStorage = storage.CreateSetIDStorage()
		}
		idx.original.Set(key, idStorage)
	}
	idx.Unlock()
	return idStorage
}

func (idx *concurrentStorage) Unwrap() Storage {
	return idx.original
}

func CreateConcurrentStorage(original Storage, storeOnlyUniqID bool) ConcurrentStorage {
	return &concurrentStorage{
		original: original,
		uniq:     storeOnlyUniqID,
	}
}
