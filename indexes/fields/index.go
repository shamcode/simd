package fields

import (
	"github.com/shamcode/simd/indexes/storage"
	"github.com/shamcode/simd/where"
)

type Storage interface {
	Get(key interface{}) *storage.IDStorage
	Set(key interface{}, records *storage.IDStorage)
	Count(key interface{}) int
	Keys() []interface{}
}

type IndexComputer interface {
	ForItem(item interface{}) interface{}
	ForComparatorAllValues(comparator where.FieldComparator, cb func(index interface{}))
	ForComparatorFirstValue(comparator where.FieldComparator) interface{}
	Compare(value interface{}, comparator where.FieldComparator) bool
}

type Index struct {
	Field   string
	Compute IndexComputer
	Storage Storage
}
