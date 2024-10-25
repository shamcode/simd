package compute

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type ComparableKey[T record.LessComparable] struct {
	Value T
}

func (i ComparableKey[T]) Less(than indexes.Key) bool {
	return i.Value < than.(ComparableKey[T]).Value
}

type lessComparableComparator[T record.LessComparable] interface {
	CompareValue(value T) (bool, error)
}

type comparableIndexComputation[R record.Record, V record.LessComparable] struct {
	getter record.ComparableGetter[R, V]
}

func (idx comparableIndexComputation[R, V]) ForRecord(item R) indexes.Key {
	return ComparableKey[V]{
		Value: idx.getter.Get(item),
	}
}

func (idx comparableIndexComputation[R, V]) ForValue(value interface{}) indexes.Key {
	return ComparableKey[V]{
		Value: value.(V),
	}
}

func (idx comparableIndexComputation[R, V]) Check(
	indexKey indexes.Key,
	comparator where.FieldComparator[R],
) (bool, error) {
	return comparator.(lessComparableComparator[V]).CompareValue(indexKey.(ComparableKey[V]).Value)
}

func CreateIndexComputation[
	R record.Record,
	T record.LessComparable,
](getter record.ComparableGetter[R, T]) indexes.IndexComputer[R] {
	return comparableIndexComputation[R, T]{getter: getter}
}
