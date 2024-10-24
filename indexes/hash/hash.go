package hash

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/storage"
)

type hashTable map[indexes.Key]storage.IDStorage

func (idx hashTable) Get(key indexes.Key) storage.IDStorage {
	return idx[key]
}

func (idx hashTable) Set(key indexes.Key, records storage.IDStorage) {
	idx[key] = records
}

func (idx hashTable) Keys() []indexes.Key {
	i := 0
	keys := make([]indexes.Key, len(idx))
	for key := range idx {
		keys[i] = key
		i += 1
	}
	return keys
}

func CreateHashTable() Storage {
	return make(hashTable)
}
