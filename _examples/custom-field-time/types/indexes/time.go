package indexes

import (
	"github.com/shamcode/simd/_examples/custom-field-time/types"
	"github.com/shamcode/simd/_examples/custom-field-time/types/comparators"
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/indexes/compute"
	"github.com/shamcode/simd/indexes/hash"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"time"
)

var _ indexes.IndexComputer = timeIndexComputation{}

type timeIndexComputation struct {
	getter *types.TimeGetter
}

func (idx timeIndexComputation) ForRecord(item record.Record) indexes.Key {
	return compute.Int64Key(idx.getter.Get(item).UnixNano())
}

func (idx timeIndexComputation) ForValue(item interface{}) indexes.Key {
	return compute.Int64Key(item.(time.Time).UnixNano())
}

func (idx timeIndexComputation) Check(indexKey indexes.Key, comparator where.FieldComparator) (bool, error) {
	return comparator.(comparators.TimeFieldComparator).CompareValue(time.Unix(0, int64(indexKey.(compute.Int64Key))))
}
func NewTimeHashIndex(getter *types.TimeGetter, unique bool) indexes.Index {
	return hash.NewIndex(
		getter.Field,
		timeIndexComputation{getter: getter},
		hash.CreateHashTable(),
		unique,
	)
}
