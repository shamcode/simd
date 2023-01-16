package fields

import (
	"github.com/shamcode/simd/record"
	"time"
)

type TimeGetter struct {
	Field string
	Get   func(item record.Record) time.Time
}
