package compute

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type Enum8Key uint8

func (i Enum8Key) Less(than indexes.Key) bool { return i < than.(Enum8Key) }

type enum8Comparator interface {
	CompareValue(value uint8) (bool, error)
}

type enum8IndexComputation struct {
	getter record.Enum8Getter
}

func (idx enum8IndexComputation) ForRecord(item record.Record) indexes.Key {
	return Enum8Key(idx.getter.Get(item).Value())
}

func (idx enum8IndexComputation) ForValue(value interface{}) indexes.Key {
	return Enum8Key(value.(record.Enum8).Value())
}

func (idx enum8IndexComputation) Check(indexKey indexes.Key, comparator where.FieldComparator) (bool, error) {
	return comparator.(enum8Comparator).CompareValue(uint8(indexKey.(Enum8Key)))
}

func CreateEnum8IndexComputation(getter record.Enum8Getter) indexes.IndexComputer {
	return enum8IndexComputation{getter: getter}
}
