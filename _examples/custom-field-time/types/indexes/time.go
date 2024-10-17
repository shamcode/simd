package indexes

import (
	"time"

	"github.com/shamcode/simd/_examples/custom-field-time/types"
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/indexes/btree"
	"github.com/shamcode/simd/indexes/compute"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type timeComparator interface {
	CompareValue(value time.Time) (bool, error)
}

type timeIndexComputation struct {
	getter types.TimeGetter
}

func (idx timeIndexComputation) ForRecord(item record.Record) indexes.Key {
	return compute.Int64Key(idx.getter.Get(item).UnixNano())
}

func (idx timeIndexComputation) ForValue(item interface{}) indexes.Key {
	return compute.Int64Key(item.(time.Time).UnixNano())
}

func (idx timeIndexComputation) Check(
	indexKey indexes.Key,
	comparator where.FieldComparator,
) (bool, error) {
	return comparator.(timeComparator).CompareValue(time.Unix(0, int64(indexKey.(compute.Int64Key)))) //nolint:wrapcheck
}

func NewTimeBTreeIndex(getter types.TimeGetter, maxChildren int, unique bool) indexes.Index {
	return btree.NewIndex(
		getter.Field,
		timeIndexComputation{getter: getter},
		btree.NewTree(maxChildren, unique),
		unique,
	)
}
