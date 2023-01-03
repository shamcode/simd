package bytype

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type IndexComputer interface {
	ForRecord(item record.Record) interface{}
	ForComparatorFirstValue(comparator where.FieldComparator) interface{}
	EachComparatorValues(comparator where.FieldComparator, cb func(index interface{}))
	Compare(value interface{}, comparator where.FieldComparator) (bool, error)
}

type Index struct {
	Field   string
	Compute IndexComputer
	Storage Storage
}
