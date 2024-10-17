package querybuilder

import (
	"time"

	"github.com/shamcode/simd/_examples/custom-field-time/types"
	"github.com/shamcode/simd/_examples/custom-field-time/types/comparators"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/where"
)

// WhereTime add condition for check field with time.Time type.
func WhereTime(getter types.TimeGetter, condition where.ComparatorType, value ...time.Time) query.BuilderOption {
	return query.AddWhereOption{
		Cmp: comparators.TimeFieldComparator{
			Cmp:    condition,
			Getter: getter,
			Value:  value,
		},
	}
}
