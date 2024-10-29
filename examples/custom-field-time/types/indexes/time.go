package indexes

import (
	"time"

	"github.com/shamcode/simd/examples/custom-field-time/types"
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/indexes/btree"
	"github.com/shamcode/simd/indexes/compute"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type timeComparator interface {
	CompareValue(value time.Time) (bool, error)
}

type timeIndexComputation[R record.Record] struct {
	getter types.TimeGetter[R]
}

func (idx timeIndexComputation[R]) ForRecord(item R) indexes.Key {
	return compute.ComparableKey[int64]{
		Value: idx.getter.Get(item).UnixNano(),
	}
}

func (idx timeIndexComputation[R]) ForValue(item interface{}) indexes.Key {
	return compute.ComparableKey[int64]{
		Value: item.(time.Time).UnixNano(),
	}
}

func (idx timeIndexComputation[R]) Check(
	indexKey indexes.Key,
	comparator where.FieldComparator[R],
) (bool, error) {
	return comparator.(timeComparator).CompareValue(
		time.Unix(0, indexKey.(compute.ComparableKey[int64]).Value),
	) //nolint:wrapcheck
}

func NewTimeBTreeIndex[R record.Record](
	getter types.TimeGetter[R],
	maxChildren int,
	unique bool,
) indexes.Index[R] {
	return btree.NewIndex[R](
		getter.Field,
		timeIndexComputation[R]{getter: getter},
		btree.NewTree(maxChildren, unique),
		unique,
	)
}
