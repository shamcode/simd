package namespace

import (
	"context"
	"github.com/shamcode/simd/query"
	"strings"
)

type QueryWithDumper interface {
	String() string
}

type debugExecutor struct {
	executor QueryExecutor
	dump     func(string)
}

func (e *debugExecutor) FetchTotal(ctx context.Context, query query.Query) (int, error) {
	e.dumpQuery(query, true)
	return e.executor.FetchTotal(ctx, query)
}

func (e *debugExecutor) FetchAll(ctx context.Context, query query.Query) (Iterator, error) {
	e.dumpQuery(query, false)
	return e.executor.FetchAll(ctx, query)
}

func (e *debugExecutor) FetchAllAndTotal(ctx context.Context, query query.Query) (Iterator, int, error) {
	e.dumpQuery(query, false)
	return e.executor.FetchAllAndTotal(ctx, query)
}

func (e *debugExecutor) dumpQuery(query query.Query, onlyTotal bool) {
	var result strings.Builder
	result.WriteString("SELECT ")
	if !onlyTotal {
		result.WriteString("*, ")
	}
	result.WriteString("COUNT(*) ")
	if dq, ok := query.(QueryWithDumper); ok {
		result.WriteString(dq.String())
	} else {
		result.WriteString("<query-dont-implement QueryWithDumper interface>")
	}
	e.dump(result.String())
}

func WrapWithDebug(executor QueryExecutor, dump func(string)) QueryExecutor {
	return &debugExecutor{
		executor: executor,
		dump:     dump,
	}
}
