package debug

import (
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
)

type QueryWithDumper[R record.Record] interface {
	query.Query[R]
	String() string
}

type debugQuery[R record.Record] struct {
	query.Query[R]

	queryDump string
}

func (dq *debugQuery[R]) String() string {
	return dq.queryDump
}

func NewQueryWithDumper[R record.Record](
	query query.Query[R],
	dumpString string,
) QueryWithDumper[R] {
	return &debugQuery[R]{
		Query:     query,
		queryDump: dumpString,
	}
}
