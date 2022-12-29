package debug

import "github.com/shamcode/simd/query"

type QueryWithDumper interface {
	query.Query
	String() string
}

var _ QueryWithDumper = (*debugQuery)(nil)

type debugQuery struct {
	query.Query
	queryDump string
}

func (dq *debugQuery) String() string {
	return dq.queryDump
}

func NewQueryWithDumper(query query.Query, dumpString string) QueryWithDumper {
	return &debugQuery{
		Query:     query,
		queryDump: dumpString,
	}
}
