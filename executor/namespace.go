package executor

import (
	"context"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type Selector[R record.Record] interface {
	PreselectForExecutor(ctx context.Context, conditions where.Conditions[R]) ([]R, error)
}
