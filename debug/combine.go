package debug

import (
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
)

var _ query.Builder = (*combineBuilder)(nil)

type combineBuilder struct {
	debug QueryBuilderWithDumper
	base  query.Builder
}

func (q *combineBuilder) Limit(limitItems int) {
	q.debug.Limit(limitItems)
	q.base.Limit(limitItems)
}

func (q *combineBuilder) Offset(startOffset int) {
	q.debug.Offset(startOffset)
	q.base.Offset(startOffset)
}

func (q *combineBuilder) Not() {
	q.debug.Not()
	q.base.Not()
}

func (q *combineBuilder) Or() {
	q.debug.Or()
	q.base.Or()
}

func (q *combineBuilder) OpenBracket() {
	q.debug.OpenBracket()
	q.base.OpenBracket()
}

func (q *combineBuilder) CloseBracket() {
	q.debug.CloseBracket()
	q.base.CloseBracket()
}

func (q *combineBuilder) AddWhere(cmp where.FieldComparator) {
	q.debug.AddWhere(cmp)
	q.base.AddWhere(cmp)
}

func (q *combineBuilder) Sort(sortBy sort.By) {
	q.debug.Sort(sortBy)
	q.base.Sort(sortBy)
}

func (q *combineBuilder) OnIteration(cb func(item record.Record)) {
	q.debug.OnIteration(cb)
	q.base.OnIteration(cb)
}

func (q *combineBuilder) Append(options ...query.BuilderOption) {
	q.debug.Append(options...)
	q.base.Append(options...)
}

func (q *combineBuilder) MakeCopy() query.Builder {
	return &combineBuilder{
		debug: q.debug.MakeCopy().(QueryBuilderWithDumper),
		base:  q.base.MakeCopy(),
	}
}

func (q *combineBuilder) Query() query.Query {
	return NewQueryWithDumper(
		q.base.Query(),
		q.debug.Dump(),
	)
}

func WrapQueryBuilder(qb query.Builder, options ...query.BuilderOption) query.Builder {
	return &combineBuilder{
		debug: CreateDebugQueryBuilder(options...),
		base:  qb,
	}
}

type QueryBuilderConstructor func(options ...query.BuilderOption) query.Builder

func WrapCreateQueryBuilder(constructor QueryBuilderConstructor) QueryBuilderConstructor {
	return func(options ...query.BuilderOption) query.Builder {
		return WrapQueryBuilder(
			constructor(options...),
			options...,
		)
	}
}
