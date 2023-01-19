package indexes

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type IndexComputer interface {
	ForRecord(item record.Record) interface{}
	ForValue(value interface{}) interface{}
	Check(indexKey interface{}, comparator where.FieldComparator) (bool, error)
}

type Index struct {
	Field   string
	Compute IndexComputer
	Storage Storage
}
