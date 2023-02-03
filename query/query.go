package query

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
)

type Query interface {
	Conditions() where.Conditions
	Sorting() []sort.ByWithOrder
	Limit() (count int, set bool)
	Offset() int
	OnIterationCallback() *func(item record.Record)
	Error() error
}

var _ Query = query{}

type query struct {
	offset              int
	limit               int
	withLimit           bool
	conditions          where.Conditions
	sorting             []sort.ByWithOrder
	onIterationCallback *func(item record.Record)
	error               error
}

func (q query) Conditions() where.Conditions {
	return q.conditions
}

func (q query) Sorting() []sort.ByWithOrder {
	return q.sorting
}

func (q query) Limit() (int, bool) {
	return q.limit, q.withLimit
}

func (q query) Offset() int {
	return q.offset
}

func (q query) OnIterationCallback() *func(item record.Record) {
	return q.onIterationCallback
}

func (q query) Error() error {
	return q.error
}
