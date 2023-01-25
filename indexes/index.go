package indexes

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/storage"
	"github.com/shamcode/simd/where"
)

type IndexComputer interface {
	ForRecord(item record.Record) interface{}
	ForValue(value interface{}) interface{}
	Check(indexKey interface{}, comparator where.FieldComparator) (bool, error)
}

type Index interface {
	Field() string
	Unique() bool
	Compute() IndexComputer
	Weight(condition where.Condition) (canApplyIndex bool, weight IndexWeight)
	Select(condition where.Condition) (count int, ids []storage.LockableIDStorage, err error)
	ConcurrentStorage() ConcurrentStorage
}

// Storage is base interface for indexes
type Storage interface {
	Get(key interface{}) storage.IDStorage
	Set(key interface{}, records storage.IDStorage)
}

// ConcurrentStorage wrapped Storage for concurrent safe access
type ConcurrentStorage interface {
	RLock()
	RUnlock()
	Unwrap() Storage
	Get(key interface{}) storage.IDStorage
	GetOrCreate(key interface{}) storage.IDStorage
}
