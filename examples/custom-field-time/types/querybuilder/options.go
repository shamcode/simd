package querybuilder

import (
	"time"

	"github.com/shamcode/simd/record"

	"github.com/shamcode/simd/examples/custom-field-time/types"
	"github.com/shamcode/simd/examples/custom-field-time/types/comparators"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/where"
)

// FieldTime add condition for check field with time.Time type.
func FieldTime[R record.Record](
	getter types.TimeGetter[R],
	condition where.ComparatorType,
	value ...time.Time,
) query.WhereOption[R] {
	return query.WhereOption[R]{
		Cmp: comparators.TimeFieldComparator[R]{
			Cmp:    condition,
			Getter: getter,
			Value:  value,
		},
		Error: nil,
	}
}
