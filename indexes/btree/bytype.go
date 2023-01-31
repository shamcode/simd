package btree

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/indexes/compute"
	"github.com/shamcode/simd/record"
)

func NewEnum8BTreeIndex(getter *record.Enum8Getter, maxChildren int, uniq bool) indexes.Index {
	return NewIndex(
		getter.Field,
		compute.CreateEnum8IndexComputation(getter),
		NewTree(maxChildren, uniq),
		uniq,
	)
}

func NewEnum16BTreeIndex(getter *record.Enum16Getter, maxChildren int, uniq bool) indexes.Index {
	return NewIndex(
		getter.Field,
		compute.CreateEnum16IndexComputation(getter),
		NewTree(maxChildren, uniq),
		uniq,
	)
}

func NewIntBTreeIndex(getter *record.IntGetter, maxChildren int, uniq bool) indexes.Index {
	return NewIndex(
		getter.Field,
		compute.CreateIntIndexComputation(getter),
		NewTree(maxChildren, uniq),
		uniq,
	)
}

func NewInt32BTreeIndex(getter *record.Int32Getter, maxChildren int, uniq bool) indexes.Index {
	return NewIndex(
		getter.Field,
		compute.CreateInt32IndexComputation(getter),
		NewTree(maxChildren, uniq),
		uniq,
	)
}

func NewInt64BTreeIndex(getter *record.Int64Getter, maxChildren int, uniq bool) indexes.Index {
	return NewIndex(
		getter.Field,
		compute.CreateInt64IndexComputation(getter),
		NewTree(maxChildren, uniq),
		uniq,
	)
}

func NewStringBTreeIndex(getter *record.StringGetter, maxChildren int, uniq bool) indexes.Index {
	return NewIndex(
		getter.Field,
		compute.CreateStringIndexComputation(getter),
		NewTree(maxChildren, uniq),
		uniq,
	)
}
