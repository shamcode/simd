package types

import (
	"github.com/shamcode/simd/record"
	"time"
)

type TimeGetter struct {
	Field record.Field
	Get   func(item record.Record) time.Time
}
