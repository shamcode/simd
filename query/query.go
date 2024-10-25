package query

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
)

type Query[R record.Record] interface {
	Conditions() where.Conditions[R]
	Sorting() []sort.ByWithOrder[R]
	Limit() (count int, set bool)
	Offset() int
	OnIterationCallback() *func(item R)
	Error() error
}

type query[R record.Record] struct {
	offset              int
	limit               int
	withLimit           bool
	conditions          where.Conditions[R]
	sorting             []sort.ByWithOrder[R]
	onIterationCallback *func(item R)
	error               error
}

func (q query[R]) Conditions() where.Conditions[R] {
	return q.conditions
}

func (q query[R]) Sorting() []sort.ByWithOrder[R] {
	return q.sorting
}

func (q query[R]) Limit() (int, bool) {
	return q.limit, q.withLimit
}

func (q query[R]) Offset() int {
	return q.offset
}

func (q query[R]) OnIterationCallback() *func(item R) {
	return q.onIterationCallback
}

func (q query[R]) Error() error {
	return q.error
}
