package executor

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type Selector[R record.Record] interface {
	PreselectForExecutor(conditions where.Conditions[R]) ([]R, error)
}
