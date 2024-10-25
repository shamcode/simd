package btree

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/indexes/compute"
	"github.com/shamcode/simd/record"
)

func NewEnumBTreeIndex[R record.Record, V record.LessComparable](
	getter record.EnumGetter[R, V],
	maxChildren int,
	uniq bool,
) indexes.Index[R] {
	return NewIndex[R](
		getter.Field,
		compute.CreateEnum8IndexComputation(getter),
		NewTree(maxChildren, uniq),
		uniq,
	)
}

func NewComparableBTreeIndex[R record.Record, V record.LessComparable](
	getter record.ComparableGetter[R, V],
	maxChildren int,
	uniq bool,
) indexes.Index[R] {
	return NewIndex[R](
		getter.Field,
		compute.CreateIndexComputation(getter),
		NewTree(maxChildren, uniq),
		uniq,
	)
}

func NewStringBTreeIndex[R record.Record](
	getter record.StringGetter[R],
	maxChildren int,
	uniq bool,
) indexes.Index[R] {
	return NewIndex(
		getter.Field,
		compute.CreateStringIndexComputation(getter),
		NewTree(maxChildren, uniq),
		uniq,
	)
}
