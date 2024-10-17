package indexes

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
)

type Key interface {
	Less(than Key) bool
}

type IndexComputer interface {
	ForRecord(item record.Record) Key
	ForValue(value interface{}) Key
	Check(indexKey Key, comparator where.FieldComparator) (bool, error)
}

type Index interface {
	Field() record.Field
	Unique() bool
	Compute() IndexComputer
	Weight(condition where.Condition) (canApplyIndex bool, weight IndexWeight)
	Select(condition where.Condition) (count int, ids []storage.IDIterator, err error)
	ConcurrentStorage() ConcurrentStorage
}

// Storage is base interface for indexes.
type Storage interface {
	Get(key Key) storage.IDStorage
	Set(key Key, records storage.IDStorage)
}

// ConcurrentStorage wrapped Storage for concurrent safe access.
type ConcurrentStorage interface {
	RLock()
	RUnlock()
	Unwrap() Storage
	Get(key Key) storage.IDStorage
	GetOrCreate(key Key) storage.IDStorage
}
