package compute

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type Int64Key int64

func (i Int64Key) Less(than indexes.Key) bool { return i < than.(Int64Key) }

type int64Comparator interface {
	CompareValue(value int64) (bool, error)
}

type int64IndexComputation struct {
	getter record.Int64Getter
}

func (idx int64IndexComputation) ForRecord(item record.Record) indexes.Key {
	return Int64Key(idx.getter.Get(item))
}

func (idx int64IndexComputation) ForValue(value interface{}) indexes.Key {
	return Int64Key(value.(int64))
}

func (idx int64IndexComputation) Check(indexKey indexes.Key, comparator where.FieldComparator) (bool, error) {
	return comparator.(int64Comparator).CompareValue(int64(indexKey.(Int64Key)))
}

func CreateInt64IndexComputation(getter record.Int64Getter) indexes.IndexComputer {
	return int64IndexComputation{getter: getter}
}
