package storage

import "github.com/shamcode/simd/record"

type LockableIDStorage interface {
	RLock()
	RUnlock()
	ThreadUnsafeData() map[int64]struct{}
}

type IDStorage interface {
	LockableIDStorage
	Count() int
	Add(id int64)
	Delete(id int64)
}

type UniqueIDStorage interface {
	IDStorage
	ID() int64
}

type RecordsByID interface {
	GetIDStorage() LockableIDStorage
	Get(id int64) record.Record
	Set(id int64, item record.Record)
	Delete(id int64)
	Count() int
	GetData(stores []LockableIDStorage, totalCount int, idsUnique bool) []record.Record
	GetAllData() []record.Record
}
