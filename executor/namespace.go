package executor

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type Selector interface {
	PreselectForExecutor(conditions where.Conditions) ([]record.Record, error)
}
