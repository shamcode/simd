package types

import (
	"time"

	"github.com/shamcode/simd/record"
)

type TimeGetter struct {
	record.Field
	Get func(item record.Record) time.Time
}

// Implement sort.By interface for sorting by fields.
func (getter TimeGetter) Less(a, b record.Record) bool {
	return getter.Get(a).Before(getter.Get(b))
}
