package debug

import (
	"context"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/query"
	"strings"
)

type QueryExecutorWithDump interface {
	executor.QueryExecutor
	DumpQuery(query query.Query, onlyTotal bool)
}

var _ executor.QueryExecutor = (*debugExecutor)(nil)

type debugExecutor struct {
	executor executor.QueryExecutor
	dump     func(string)
}

func (e *debugExecutor) FetchTotal(ctx context.Context, query query.Query) (int, error) {
	e.DumpQuery(query, true)
	return e.executor.FetchTotal(ctx, query)
}

func (e *debugExecutor) FetchAll(ctx context.Context, query query.Query) (executor.Iterator, error) {
	e.DumpQuery(query, false)
	return e.executor.FetchAll(ctx, query)
}

func (e *debugExecutor) FetchAllAndTotal(ctx context.Context, query query.Query) (executor.Iterator, int, error) {
	e.DumpQuery(query, false)
	return e.executor.FetchAllAndTotal(ctx, query)
}

func (e *debugExecutor) DumpQuery(query query.Query, onlyTotal bool) {
	var result strings.Builder
	result.WriteString("SELECT ")
	if !onlyTotal {
		result.WriteString("*, ")
	}
	result.WriteString("COUNT(*)")
	if dq, ok := query.(QueryWithDumper); ok {
		result.WriteString(dq.String())
	} else {
		result.WriteString(" <Query dont implement QueryWithDumper interface, check QueryBuilder>")
	}
	e.dump(result.String())
}

func WrapQueryExecutor(executor executor.QueryExecutor, dump func(string)) executor.QueryExecutor {
	return &debugExecutor{
		executor: executor,
		dump:     dump,
	}
}
