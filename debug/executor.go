package debug

import (
	"context"
	"strings"

	"github.com/shamcode/simd/record"

	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/query"
)

type QueryExecutorWithDump[R record.Record] interface {
	executor.QueryExecutor[R]
	DumpQuery(ctx context.Context, query query.Query[R], onlyTotal bool)
}

type debugExecutor[R record.Record] struct {
	executor executor.QueryExecutor[R]
	dump     func(ctx context.Context, q string)
}

func (e *debugExecutor[R]) FetchTotal(
	ctx context.Context,
	query query.Query[R],
) (int, error) {
	e.DumpQuery(ctx, query, true)
	return e.executor.FetchTotal(ctx, query)
}

func (e *debugExecutor[R]) FetchAll(
	ctx context.Context,
	query query.Query[R],
) (executor.Iterator[R], error) {
	e.DumpQuery(ctx, query, false)
	return e.executor.FetchAll(ctx, query)
}

func (e *debugExecutor[R]) FetchAllAndTotal(
	ctx context.Context,
	query query.Query[R],
) (executor.Iterator[R], int, error) {
	e.DumpQuery(ctx, query, false)
	return e.executor.FetchAllAndTotal(ctx, query)
}

func (e *debugExecutor[R]) DumpQuery(ctx context.Context, query query.Query[R], onlyTotal bool) {
	var result strings.Builder
	result.WriteString("SELECT ")

	if !onlyTotal {
		result.WriteString("*, ")
	}

	result.WriteString("COUNT(*)")

	if dq, ok := any(query).(QueryWithDumper[R]); ok {
		result.WriteString(dq.String())
	} else {
		result.WriteString(" <Query dont implement QueryWithDumper interface, check QueryBuilder>")
	}

	e.dump(ctx, result.String())
}

func WrapQueryExecutor[R record.Record](
	executor executor.QueryExecutor[R],
	dump func(ctx context.Context, q string),
) executor.QueryExecutor[R] {
	return &debugExecutor[R]{
		executor: executor,
		dump:     dump,
	}
}
