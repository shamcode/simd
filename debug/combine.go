package debug

import (
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
)

type combineBuilder[R record.Record] struct {
	debug QueryBuilderWithDumper[R]
	base  query.BuilderGeneric[R]
}

func (q *combineBuilder[R]) Limit(limitItems int) {
	q.debug.Limit(limitItems)
	q.base.Limit(limitItems)
}

func (q *combineBuilder[R]) Offset(startOffset int) {
	q.debug.Offset(startOffset)
	q.base.Offset(startOffset)
}

func (q *combineBuilder[R]) Not() {
	q.debug.Not()
	q.base.Not()
}

func (q *combineBuilder[R]) Or() {
	q.debug.Or()
	q.base.Or()
}

func (q *combineBuilder[R]) OpenBracket() {
	q.debug.OpenBracket()
	q.base.OpenBracket()
}

func (q *combineBuilder[R]) CloseBracket() {
	q.debug.CloseBracket()
	q.base.CloseBracket()
}

func (q *combineBuilder[R]) Error(err error) {
	q.debug.Error(err)
	q.base.Error(err)
}

func (q *combineBuilder[R]) AddWhere(cmp where.FieldComparator[R]) {
	q.debug.AddWhere(cmp)
	q.base.AddWhere(cmp)
}

func (q *combineBuilder[R]) Sort(sortBy sort.ByWithOrder[R]) {
	q.debug.Sort(sortBy)
	q.base.Sort(sortBy)
}

func (q *combineBuilder[R]) OnIteration(cb func(item R)) {
	q.debug.OnIteration(cb)
	q.base.OnIteration(cb)
}

func (q *combineBuilder[R]) Append(options ...query.BuilderOption) {
	q.debug.Append(options...)
	q.base.Append(options...)
}

func (q *combineBuilder[R]) MakeCopy() query.BuilderGeneric[R] {
	return &combineBuilder[R]{
		debug: q.debug.MakeCopy().(QueryBuilderWithDumper[R]),
		base:  q.base.MakeCopy(),
	}
}

func (q *combineBuilder[R]) Query() query.Query[R] {
	return NewQueryWithDumper[R](
		q.base.Query(),
		q.debug.Dump(),
	)
}

func WrapQueryBuilder[R record.Record](
	qb query.BuilderGeneric[R],
	options ...query.BuilderOption,
) query.BuilderGeneric[R] {
	return &combineBuilder[R]{
		debug: CreateDebugQueryBuilder[R](options...),
		base:  qb,
	}
}

type QueryBuilderConstructor[R record.Record] func(options ...query.BuilderOption) query.BuilderGeneric[R]

func WrapCreateQueryBuilder[R record.Record](constructor QueryBuilderConstructor[R]) QueryBuilderConstructor[R] {
	return func(options ...query.BuilderOption) query.BuilderGeneric[R] {
		return WrapQueryBuilder(
			constructor(options...),
			options...,
		)
	}
}

func WrapCreateQueryBuilderWithDumper[R record.Record](
	constructor QueryBuilderConstructor[R],
	dumper FieldComparatorDumper[R],
) QueryBuilderConstructor[R] {
	return func(options ...query.BuilderOption) query.BuilderGeneric[R] {
		debug := CreateDebugQueryBuilder[R]()
		debug.SetFieldComparatorDumper(dumper) // setup dumper before apply options
		debug.Append(options...)
		return &combineBuilder[R]{
			debug: debug,
			base:  constructor(options...),
		}
	}
}
