package btree

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/indexes/compute"
	"github.com/shamcode/simd/record"
)

func NewComparableBTreeIndex[R record.Record, V record.LessComparable](
	getter record.GetterInterface[R, V],
	maxChildren int,
	uniq bool,
) indexes.Index[R] {
	return NewIndex[R](
		getter,
		compute.CreateIndexComputation(getter),
		NewTree(maxChildren, uniq),
		uniq,
	)
}
