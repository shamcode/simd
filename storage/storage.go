package storage

import "github.com/shamcode/simd/record"

type IDIterator interface {
	Iterate(f func(id int64))
}

type IDStorage interface {
	IDIterator
	Count() int
	Add(id int64)
	Delete(id int64)
}

type UniqueIDStorage interface {
	IDStorage
	ID() int64
}

type RecordsByID[R record.Record] interface {
	GetIDStorage() IDIterator
	Get(id int64) (R, bool)
	Set(id int64, item R)
	Delete(id int64)
	Count() int
	GetData(stores []IDIterator, totalCount int, idsUnique bool) []R
	GetAllData() []R
}
