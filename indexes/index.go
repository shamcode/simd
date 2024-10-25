package indexes

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
)

type Key interface {
	Less(than Key) bool
}

type IndexComputer[R record.Record] interface {
	ForRecord(item R) Key
	ForValue(value interface{}) Key
	Check(indexKey Key, comparator where.FieldComparator[R]) (bool, error)
}

type Index[R record.Record] interface {
	Field() record.Field
	Unique() bool
	Compute() IndexComputer[R]
	Weight(condition where.Condition[R]) (canApplyIndex bool, weight IndexWeight)
	Select(condition where.Condition[R]) (count int, ids []storage.IDIterator, err error)
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
