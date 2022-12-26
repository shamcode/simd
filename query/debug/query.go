package debug

import "github.com/shamcode/simd/query"

type debugQuery struct {
	query.Query
	queryDump string
}

func (dq *debugQuery) String() string {
	return dq.queryDump
}
