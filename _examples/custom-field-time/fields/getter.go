package fields

import "time"

type TimeGetter struct {
	Field string
	Get   func(item interface{}) time.Time
}
