package types

import (
	"time"

	"github.com/shamcode/simd/record"
)

type TimeGetter[R record.Record] struct {
	record.Field
	Get func(item R) time.Time
}

// Less implement sort.By interface for sorting by fields.
func (getter TimeGetter[R]) Less(a, b R) bool {
	return getter.Get(a).Before(getter.Get(b))
}
