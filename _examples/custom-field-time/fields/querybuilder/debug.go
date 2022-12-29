package querybuilder

import (
	"github.com/shamcode/simd/_examples/custom-field-time/fields"
	"github.com/shamcode/simd/debug"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"regexp"
	"time"
)

var _ QueryBuilder = (*debugQueryBuilder)(nil)

type debugQueryBuilder struct {
	builder QueryBuilder
	debug   debug.BaseQueryBuilderWithDump
}

func (d *debugQueryBuilder) Query() query.Query {
	return debug.NewQueryWithDumper(
		d.builder.Query(),
		d.debug.Dump(),
	)
}

func (d *debugQueryBuilder) MakeCopy() QueryBuilder {
	return &debugQueryBuilder{
		builder: d.MakeCopy(),
		debug:   d.debug.MakeCopy().(debug.BaseQueryBuilderWithDump),
	}
}

func (d *debugQueryBuilder) Limit(limitItems int) QueryBuilder {
	d.builder.Limit(limitItems)
	d.debug.Limit(limitItems)
	return d
}

func (d *debugQueryBuilder) Offset(startOffset int) QueryBuilder {
	d.builder.Offset(startOffset)
	d.debug.Offset(startOffset)
	return d
}

func (d *debugQueryBuilder) Not() QueryBuilder {
	d.builder.Not()
	d.debug.Not()
	return d
}

func (d *debugQueryBuilder) Or() QueryBuilder {
	d.builder.Or()
	d.debug.Or()
	return d
}

func (d *debugQueryBuilder) OpenBracket() QueryBuilder {
	d.builder.OpenBracket()
	d.debug.OpenBracket()
	return d
}

func (d *debugQueryBuilder) CloseBracket() QueryBuilder {
	d.builder.CloseBracket()
	d.debug.CloseBracket()
	return d
}

func (d *debugQueryBuilder) AddWhere(cmp where.FieldComparator) QueryBuilder {
	d.builder.AddWhere(cmp)
	d.debug.AddWhere(cmp)
	return d
}

func (d *debugQueryBuilder) Where(getter *record.InterfaceGetter, condition where.ComparatorType, values ...interface{}) QueryBuilder {
	d.builder.Where(getter, condition, values...)
	d.debug.Where(getter, condition, values...)
	return d
}

func (d *debugQueryBuilder) WhereInt(getter *record.IntGetter, condition where.ComparatorType, values ...int) QueryBuilder {
	d.builder.WhereInt(getter, condition, values...)
	d.debug.WhereInt(getter, condition, values...)
	return d
}

func (d *debugQueryBuilder) WhereInt32(getter *record.Int32Getter, condition where.ComparatorType, values ...int32) QueryBuilder {
	d.builder.WhereInt32(getter, condition, values...)
	d.debug.WhereInt32(getter, condition, values...)
	return d
}

func (d *debugQueryBuilder) WhereInt64(getter *record.Int64Getter, condition where.ComparatorType, values ...int64) QueryBuilder {
	d.builder.WhereInt64(getter, condition, values...)
	d.debug.WhereInt64(getter, condition, values...)
	return d
}

func (d *debugQueryBuilder) WhereString(getter *record.StringGetter, condition where.ComparatorType, values ...string) QueryBuilder {
	d.builder.WhereString(getter, condition, values...)
	d.debug.WhereString(getter, condition, values...)
	return d
}

func (d *debugQueryBuilder) WhereStringRegexp(getter *record.StringGetter, value *regexp.Regexp) QueryBuilder {
	d.builder.WhereStringRegexp(getter, value)
	d.debug.WhereStringRegexp(getter, value)
	return d
}

func (d *debugQueryBuilder) WhereBool(getter *record.BoolGetter, condition where.ComparatorType, values ...bool) QueryBuilder {
	d.builder.WhereBool(getter, condition, values...)
	d.debug.WhereBool(getter, condition, values...)
	return d
}

func (d *debugQueryBuilder) WhereEnum8(getter *record.Enum8Getter, condition where.ComparatorType, values ...record.Enum8) QueryBuilder {
	d.builder.WhereEnum8(getter, condition, values...)
	d.debug.WhereEnum8(getter, condition, values...)
	return d
}

func (d *debugQueryBuilder) WhereEnum16(getter *record.Enum16Getter, condition where.ComparatorType, values ...record.Enum16) QueryBuilder {
	d.builder.WhereEnum16(getter, condition, values...)
	d.debug.WhereEnum16(getter, condition, values...)
	return d
}

func (d *debugQueryBuilder) WhereMap(getter *record.MapGetter, condition where.ComparatorType, values ...interface{}) QueryBuilder {
	d.builder.WhereMap(getter, condition, values...)
	d.debug.WhereMap(getter, condition, values...)
	return d
}

func (d *debugQueryBuilder) WhereSet(getter *record.SetGetter, condition where.ComparatorType, values ...interface{}) QueryBuilder {
	d.builder.WhereSet(getter, condition, values...)
	d.debug.WhereSet(getter, condition, values...)
	return d
}

func (d *debugQueryBuilder) Sort(by sort.By) QueryBuilder {
	d.builder.Sort(by)
	d.debug.Sort(by)
	return d
}

func (d *debugQueryBuilder) OnIteration(cb func(item record.Record)) QueryBuilder {
	d.builder.OnIteration(cb)
	d.debug.OnIteration(cb)
	return d
}

func (d *debugQueryBuilder) WhereTime(getter *fields.TimeGetter, condition where.ComparatorType, value ...time.Time) QueryBuilder {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	d.debug.SaveWhereForDump(getter.Field, condition, items...)
	d.builder.WhereTime(getter, condition, value...)
	return d
}

func WrapQueryBuilder(qb QueryBuilder) QueryBuilder {
	return &debugQueryBuilder{
		debug:   debug.CreateDebugQueryBuilder(),
		builder: qb,
	}
}

type Constructor func() QueryBuilder

func WrapWithDebug(constructor Constructor) Constructor {
	return func() QueryBuilder {
		return WrapQueryBuilder(constructor())
	}
}
