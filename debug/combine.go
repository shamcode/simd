package debug

import (
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"regexp"
)

var _ query.BaseQueryBuilder = (*combineBuilder)(nil)

type combineBuilder struct {
	debug BaseQueryBuilderWithDump
	base  query.BaseQueryBuilder
}

func (q *combineBuilder) OnIteration(cb func(item record.Record)) query.BaseQueryBuilder {
	q.debug.OnIteration(cb)
	q.base.OnIteration(cb)
	return q
}

func (q *combineBuilder) Limit(limitItems int) query.BaseQueryBuilder {
	q.debug.Limit(limitItems)
	q.base.Limit(limitItems)
	return q
}

func (q *combineBuilder) Offset(startOffset int) query.BaseQueryBuilder {
	q.debug.Offset(startOffset)
	q.base.Offset(startOffset)
	return q
}

func (q *combineBuilder) Or() query.BaseQueryBuilder {
	q.debug.Or()
	q.base.Or()
	return q
}

func (q *combineBuilder) Not() query.BaseQueryBuilder {
	q.debug.Not()
	q.base.Not()
	return q
}

func (q *combineBuilder) OpenBracket() query.BaseQueryBuilder {
	q.debug.OpenBracket()
	q.base.OpenBracket()
	return q
}

func (q *combineBuilder) CloseBracket() query.BaseQueryBuilder {
	q.debug.CloseBracket()
	q.base.CloseBracket()
	return q
}

func (q *combineBuilder) SaveWhereForDump(field string, condition where.ComparatorType, value ...interface{}) {
	q.debug.SaveWhereForDump(field, condition, value...)
}

func (q *combineBuilder) AddWhere(cmp where.FieldComparator) query.BaseQueryBuilder {
	q.debug.AddWhere(cmp)
	q.base.AddWhere(cmp)
	return q
}

func (q *combineBuilder) Where(getter *record.InterfaceGetter, condition where.ComparatorType, value ...interface{}) query.BaseQueryBuilder {
	q.debug.Where(getter, condition, value...)
	q.base.Where(getter, condition, value...)
	return q
}

func (q *combineBuilder) WhereInt(getter *record.IntGetter, condition where.ComparatorType, value ...int) query.BaseQueryBuilder {
	q.debug.WhereInt(getter, condition, value...)
	q.base.WhereInt(getter, condition, value...)
	return q
}

func (q *combineBuilder) WhereInt32(getter *record.Int32Getter, condition where.ComparatorType, value ...int32) query.BaseQueryBuilder {
	q.debug.WhereInt32(getter, condition, value...)
	q.base.WhereInt32(getter, condition, value...)
	return q
}

func (q *combineBuilder) WhereInt64(getter *record.Int64Getter, condition where.ComparatorType, value ...int64) query.BaseQueryBuilder {
	q.debug.WhereInt64(getter, condition, value...)
	q.base.WhereInt64(getter, condition, value...)
	return q
}

func (q *combineBuilder) WhereString(getter *record.StringGetter, condition where.ComparatorType, value ...string) query.BaseQueryBuilder {
	q.debug.WhereString(getter, condition, value...)
	q.base.WhereString(getter, condition, value...)
	return q
}

func (q *combineBuilder) WhereStringRegexp(getter *record.StringGetter, value *regexp.Regexp) query.BaseQueryBuilder {
	q.debug.WhereStringRegexp(getter, value)
	q.base.WhereStringRegexp(getter, value)
	return q
}

func (q *combineBuilder) WhereBool(getter *record.BoolGetter, condition where.ComparatorType, value ...bool) query.BaseQueryBuilder {
	q.debug.WhereBool(getter, condition, value...)
	q.base.WhereBool(getter, condition, value...)
	return q
}

func (q *combineBuilder) WhereEnum8(getter *record.Enum8Getter, condition where.ComparatorType, value ...record.Enum8) query.BaseQueryBuilder {
	q.debug.WhereEnum8(getter, condition, value...)
	q.base.WhereEnum8(getter, condition, value...)
	return q
}

func (q *combineBuilder) WhereEnum16(getter *record.Enum16Getter, condition where.ComparatorType, value ...record.Enum16) query.BaseQueryBuilder {
	q.debug.WhereEnum16(getter, condition, value...)
	q.base.WhereEnum16(getter, condition, value...)
	return q
}

func (q *combineBuilder) WhereMap(getter *record.MapGetter, condition where.ComparatorType, value ...interface{}) query.BaseQueryBuilder {
	q.debug.WhereMap(getter, condition, value...)
	q.base.WhereMap(getter, condition, value...)
	return q
}

func (q *combineBuilder) WhereSet(getter *record.SetGetter, condition where.ComparatorType, value ...interface{}) query.BaseQueryBuilder {
	q.debug.WhereSet(getter, condition, value...)
	q.base.WhereSet(getter, condition, value...)
	return q
}

func (q *combineBuilder) MakeCopy() query.BaseQueryBuilder {
	return &combineBuilder{
		debug: q.debug.MakeCopy().(BaseQueryBuilderWithDump),
		base:  q.base.MakeCopy(),
	}
}

func (q *combineBuilder) Sort(sortBy sort.By) query.BaseQueryBuilder {
	q.debug.Sort(sortBy)
	q.base.Sort(sortBy)
	return q
}

func (q *combineBuilder) Query() query.Query {
	return NewQueryWithDumper(
		q.base.Query(),
		q.debug.Dump(),
	)
}

func WrapQueryBuilder(qb query.BaseQueryBuilder) query.BaseQueryBuilder {
	return &combineBuilder{
		debug: CreateDebugQueryBuilder(),
		base:  qb,
	}
}

type QueryBuilderConstructor func() query.BaseQueryBuilder

func WrapCreateQueryBuilder(constructor QueryBuilderConstructor) QueryBuilderConstructor {
	return func() query.BaseQueryBuilder {
		return WrapQueryBuilder(constructor())
	}
}
