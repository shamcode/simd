package debug

import (
	"context"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/query"
	"strings"
)

type QueryExecutorWithDump interface {
	namespace.QueryExecutor
	DumpQuery(query query.Query, onlyTotal bool)
}

var _ namespace.QueryExecutor = (*debugExecutor)(nil)

type debugExecutor struct {
	executor namespace.QueryExecutor
	dump     func(string)
}

func (e *debugExecutor) FetchTotal(ctx context.Context, query query.Query) (int, error) {
	e.DumpQuery(query, true)
	return e.executor.FetchTotal(ctx, query)
}

func (e *debugExecutor) FetchAll(ctx context.Context, query query.Query) (namespace.Iterator, error) {
	e.DumpQuery(query, false)
	return e.executor.FetchAll(ctx, query)
}

func (e *debugExecutor) FetchAllAndTotal(ctx context.Context, query query.Query) (namespace.Iterator, int, error) {
	e.DumpQuery(query, false)
	return e.executor.FetchAllAndTotal(ctx, query)
}

func (e *debugExecutor) DumpQuery(query query.Query, onlyTotal bool) {
	var result strings.Builder
	result.WriteString("SELECT ")
	if !onlyTotal {
		result.WriteString("*, ")
	}
	result.WriteString("COUNT(*) ")
	if dq, ok := query.(QueryWithDumper); ok {
		result.WriteString(dq.String())
	} else {
		result.WriteString("<Query dont implement QueryWithDumper interface, check QueryBuilder>")
	}
	e.dump(result.String())
}

func WrapQueryExecutor(executor namespace.QueryExecutor, dump func(string)) namespace.QueryExecutor {
	return &debugExecutor{
		executor: executor,
		dump:     dump,
	}
}
