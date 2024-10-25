package compute

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type EnumKey[T record.LessComparable] struct {
	Value T
}

func (i EnumKey[T]) Less(than indexes.Key) bool { return i.Value < than.(EnumKey[T]).Value }

type enumComparator[T record.LessComparable] interface {
	CompareValue(value T) (bool, error)
}

type enum8IndexComputation[R record.Record, V record.LessComparable] struct {
	getter record.EnumGetter[R, V]
}

func (idx enum8IndexComputation[R, V]) ForRecord(item R) indexes.Key {
	return EnumKey[V]{
		Value: idx.getter.Get(item).Value(),
	}
}

func (idx enum8IndexComputation[R, V]) ForValue(value interface{}) indexes.Key {
	return EnumKey[V]{
		Value: value.(record.Enum[V]).Value(),
	}
}

func (idx enum8IndexComputation[R, V]) Check(
	indexKey indexes.Key,
	comparator where.FieldComparator[R],
) (bool, error) {
	return comparator.(enumComparator[V]).CompareValue(indexKey.(EnumKey[V]).Value)
}

func CreateEnum8IndexComputation[
	R record.Record,
	V record.LessComparable,
](getter record.EnumGetter[R, V]) indexes.IndexComputer[R] {
	return enum8IndexComputation[R, V]{getter: getter}
}
