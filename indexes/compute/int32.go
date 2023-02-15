package compute

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type Int32Key int32

func (i Int32Key) Less(than indexes.Key) bool { return i < than.(Int32Key) }

type int32Comparator interface {
	CompareValue(value int32) (bool, error)
}

var _ indexes.IndexComputer = int32IndexComputation{}

type int32IndexComputation struct {
	getter record.Int32Getter
}

func (idx int32IndexComputation) ForRecord(item record.Record) indexes.Key {
	return Int32Key(idx.getter.Get(item))
}

func (idx int32IndexComputation) ForValue(value interface{}) indexes.Key {
	return Int32Key(value.(int32))
}

func (idx int32IndexComputation) Check(indexKey indexes.Key, comparator where.FieldComparator) (bool, error) {
	return comparator.(int32Comparator).CompareValue(int32(indexKey.(Int32Key)))
}

func CreateInt32IndexComputation(getter record.Int32Getter) indexes.IndexComputer {
	return int32IndexComputation{getter: getter}
}
