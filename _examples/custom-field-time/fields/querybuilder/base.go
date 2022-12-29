package querybuilder

import (
	"github.com/shamcode/simd/_examples/custom-field-time/fields"
	"github.com/shamcode/simd/_examples/custom-field-time/fields/comparators"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"regexp"
	"time"
)

type QueryBuilder interface {
	MakeCopy() QueryBuilder
	Limit(limitItems int) QueryBuilder
	Offset(startOffset int) QueryBuilder
	Not() QueryBuilder
	Or() QueryBuilder
	OpenBracket() QueryBuilder
	CloseBracket() QueryBuilder
	AddWhere(cmp where.FieldComparator) QueryBuilder
	Where(getter *record.InterfaceGetter, condition where.ComparatorType, values ...interface{}) QueryBuilder
	WhereInt(getter *record.IntGetter, condition where.ComparatorType, values ...int) QueryBuilder
	WhereInt32(getter *record.Int32Getter, condition where.ComparatorType, values ...int32) QueryBuilder
	WhereInt64(getter *record.Int64Getter, condition where.ComparatorType, values ...int64) QueryBuilder
	WhereString(getter *record.StringGetter, condition where.ComparatorType, values ...string) QueryBuilder
	WhereStringRegexp(getter *record.StringGetter, value *regexp.Regexp) QueryBuilder
	WhereBool(getter *record.BoolGetter, condition where.ComparatorType, values ...bool) QueryBuilder
	WhereEnum8(getter *record.Enum8Getter, condition where.ComparatorType, values ...record.Enum8) QueryBuilder
	WhereEnum16(getter *record.Enum16Getter, condition where.ComparatorType, values ...record.Enum16) QueryBuilder
	WhereMap(getter *record.MapGetter, condition where.ComparatorType, values ...interface{}) QueryBuilder
	WhereSet(getter *record.SetGetter, condition where.ComparatorType, values ...interface{}) QueryBuilder
	Sort(by sort.By) QueryBuilder
	OnIteration(cb func(item record.Record)) QueryBuilder
	Query() query.Query

	// WhereTime add condition for check field with time.Time type
	WhereTime(getter *fields.TimeGetter, condition where.ComparatorType, value ...time.Time) QueryBuilder
}

var _ QueryBuilder = (*queryBuilder)(nil)

type queryBuilder struct {
	builder query.BaseQueryBuilder
}

func (q *queryBuilder) Query() query.Query {
	return q.builder.Query()
}

func (q *queryBuilder) MakeCopy() QueryBuilder {
	return &queryBuilder{
		builder: q.builder.MakeCopy(),
	}
}

func (q *queryBuilder) Limit(limitItems int) QueryBuilder {
	q.builder.Limit(limitItems)
	return q
}

func (q *queryBuilder) Offset(startOffset int) QueryBuilder {
	q.builder.Offset(startOffset)
	return q
}

func (q *queryBuilder) Not() QueryBuilder {
	q.builder.Not()
	return q
}

func (q *queryBuilder) Or() QueryBuilder {
	q.builder.Not()
	return q
}

func (q *queryBuilder) OpenBracket() QueryBuilder {
	q.builder.OpenBracket()
	return q
}

func (q *queryBuilder) CloseBracket() QueryBuilder {
	q.builder.CloseBracket()
	return q
}

func (q *queryBuilder) AddWhere(cmp where.FieldComparator) QueryBuilder {
	q.builder.AddWhere(cmp)
	return q
}

func (q *queryBuilder) Where(getter *record.InterfaceGetter, condition where.ComparatorType, values ...interface{}) QueryBuilder {
	q.builder.Where(getter, condition, values...)
	return q
}

func (q *queryBuilder) WhereInt(getter *record.IntGetter, condition where.ComparatorType, values ...int) QueryBuilder {
	q.builder.WhereInt(getter, condition, values...)
	return q
}

func (q *queryBuilder) WhereInt32(getter *record.Int32Getter, condition where.ComparatorType, values ...int32) QueryBuilder {
	q.builder.WhereInt32(getter, condition, values...)
	return q
}

func (q *queryBuilder) WhereInt64(getter *record.Int64Getter, condition where.ComparatorType, values ...int64) QueryBuilder {
	q.builder.WhereInt64(getter, condition, values...)
	return q
}

func (q *queryBuilder) WhereString(getter *record.StringGetter, condition where.ComparatorType, values ...string) QueryBuilder {
	q.builder.WhereString(getter, condition, values...)
	return q
}

func (q *queryBuilder) WhereStringRegexp(getter *record.StringGetter, value *regexp.Regexp) QueryBuilder {
	q.builder.WhereStringRegexp(getter, value)
	return q
}

func (q *queryBuilder) WhereBool(getter *record.BoolGetter, condition where.ComparatorType, values ...bool) QueryBuilder {
	q.builder.WhereBool(getter, condition, values...)
	return q
}

func (q *queryBuilder) WhereEnum8(getter *record.Enum8Getter, condition where.ComparatorType, values ...record.Enum8) QueryBuilder {
	q.builder.WhereEnum8(getter, condition, values...)
	return q
}

func (q *queryBuilder) WhereEnum16(getter *record.Enum16Getter, condition where.ComparatorType, values ...record.Enum16) QueryBuilder {
	q.builder.WhereEnum16(getter, condition, values...)
	return q
}

func (q *queryBuilder) WhereMap(getter *record.MapGetter, condition where.ComparatorType, values ...interface{}) QueryBuilder {
	q.builder.WhereMap(getter, condition, values...)
	return q
}

func (q *queryBuilder) WhereSet(getter *record.SetGetter, condition where.ComparatorType, values ...interface{}) QueryBuilder {
	q.builder.WhereSet(getter, condition, values...)
	return q
}

func (q *queryBuilder) Sort(by sort.By) QueryBuilder {
	q.builder.Sort(by)
	return q
}

func (q *queryBuilder) OnIteration(cb func(item record.Record)) QueryBuilder {
	q.builder.OnIteration(cb)
	return q
}

func (q *queryBuilder) WhereTime(getter *fields.TimeGetter, condition where.ComparatorType, value ...time.Time) QueryBuilder {
	q.builder.AddWhere(comparators.TimeFieldComparator{
		Cmp:    condition,
		Getter: getter,
		Value:  value,
	})
	return q
}

func Create() QueryBuilder {
	return &queryBuilder{
		builder: query.NewBuilder(),
	}
}
