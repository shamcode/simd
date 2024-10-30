package hash

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/indexes/compute"
	"github.com/shamcode/simd/record"
)

func NewBoolHashIndex[R record.Record](
	getter record.BoolGetter[R],
	unique bool,
) indexes.Index[R] {
	return NewIndex(
		getter.Field,
		compute.CreateBoolIndexComputation(getter),
		CreateHashTable(),
		unique,
	)
}

func NewEnumHashIndex[R record.Record, T record.LessComparable](
	getter record.EnumGetter[R, T],
	unique bool,
) indexes.Index[R] {
	return NewIndex(
		getter.Field,
		compute.CreateEnum8IndexComputation(getter),
		CreateHashTable(),
		unique,
	)
}

func NewComparableHashIndex[R record.Record, T record.LessComparable](
	getter record.ComparableGetter[R, T],
	unique bool,
) indexes.Index[R] {
	return NewIndex(
		getter.Field,
		compute.CreateIndexComputation(getter),
		CreateHashTable(),
		unique,
	)
}
