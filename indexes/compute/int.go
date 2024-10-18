package compute

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type IntKey int

func (i IntKey) Less(than indexes.Key) bool { return i < than.(IntKey) }

type intComparator interface {
	CompareValue(value int) (bool, error)
}

type intIndexComputation struct {
	getter record.IntGetter
}

func (idx intIndexComputation) ForRecord(item record.Record) indexes.Key {
	return IntKey(idx.getter.Get(item))
}

func (idx intIndexComputation) ForValue(value interface{}) indexes.Key {
	return IntKey(value.(int))
}

func (idx intIndexComputation) Check(indexKey indexes.Key, comparator where.FieldComparator) (bool, error) {
	return comparator.(intComparator).CompareValue(int(indexKey.(IntKey)))
}

func CreateIntIndexComputation(getter record.IntGetter) indexes.IndexComputer {
	return intIndexComputation{getter: getter}
}
