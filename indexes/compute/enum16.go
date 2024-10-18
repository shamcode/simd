package compute

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type Enum16Key uint16

func (i Enum16Key) Less(than indexes.Key) bool { return i < than.(Enum16Key) }

type enum16Comparator interface {
	CompareValue(value uint16) (bool, error)
}

type enum16IndexComputation struct {
	getter record.Enum16Getter
}

func (idx enum16IndexComputation) ForRecord(item record.Record) indexes.Key {
	return Enum16Key(idx.getter.Get(item).Value())
}

func (idx enum16IndexComputation) ForValue(value interface{}) indexes.Key {
	return Enum16Key(value.(record.Enum16).Value())
}

func (idx enum16IndexComputation) Check(indexKey indexes.Key, comparator where.FieldComparator) (bool, error) {
	return comparator.(enum16Comparator).CompareValue(uint16(indexKey.(Enum16Key)))
}

func CreateEnum16IndexComputation(getter record.Enum16Getter) indexes.IndexComputer {
	return enum16IndexComputation{getter: getter}
}
