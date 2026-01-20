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
	DumpQuery(query query.Query[R], onlyTotal bool)
}

type debugExecutor[R record.Record] struct {
	executor executor.QueryExecutor[R]
	dump     func(string)
}

func (e *debugExecutor[R]) FetchTotal(
	ctx context.Context,
	query query.Query[R],
) (int, error) {
	e.DumpQuery(query, true)
	return e.executor.FetchTotal(ctx, query)
}

func (e *debugExecutor[R]) FetchAll(
	ctx context.Context,
	query query.Query[R],
) (executor.Iterator[R], error) {
	e.DumpQuery(query, false)
	return e.executor.FetchAll(ctx, query)
}

func (e *debugExecutor[R]) FetchAllAndTotal(
	ctx context.Context,
	query query.Query[R],
) (executor.Iterator[R], int, error) {
	e.DumpQuery(query, false)
	return e.executor.FetchAllAndTotal(ctx, query)
}

func (e *debugExecutor[R]) DumpQuery(query query.Query[R], onlyTotal bool) {
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

	e.dump(result.String())
}

func WrapQueryExecutor[R record.Record](
	executor executor.QueryExecutor[R],
	dump func(string),
) executor.QueryExecutor[R] {
	return &debugExecutor[R]{
		executor: executor,
		dump:     dump,
	}
}
