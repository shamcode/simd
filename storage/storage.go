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

type RecordsByID interface {
	GetIDStorage() IDIterator
	Get(id int64) record.Record
	Set(id int64, item record.Record)
	Delete(id int64)
	Count() int
	GetData(stores []IDIterator, totalCount int, idsUnique bool) []record.Record
	GetAllData() []record.Record
}
