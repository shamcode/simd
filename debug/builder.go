package debug

import (
	"fmt"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"regexp"
	"strconv"
	"strings"
)

type QueryBuilderWithDump interface {
	query.Builder
	SaveWhereForDump(field string, condition where.ComparatorType, value ...interface{})
}

const (
	chunkLimit uint8 = iota + 1
	chunkOffset
	chunkWhere
	chunkSort
)

var _ QueryBuilderWithDump = (*debugQueryBuilder)(nil)

type debugQueryBuilder struct {
	builder   query.Builder
	chunks    map[uint8]*strings.Builder
	requireOp bool
	withNot   bool
	isOr      bool
}

func (q *debugQueryBuilder) OnIteration(cb func(item record.Record)) query.Builder {
	q.builder.OnIteration(cb)
	return q
}

func (q *debugQueryBuilder) Limit(limitItems int) query.Builder {
	w := q.chunks[chunkLimit]
	w.WriteString("LIMIT ")
	w.WriteString(strconv.Itoa(limitItems))
	q.builder.Limit(limitItems)
	return q
}

func (q *debugQueryBuilder) Offset(startOffset int) query.Builder {
	w := q.chunks[chunkOffset]
	w.WriteString("OFFSET ")
	w.WriteString(strconv.Itoa(startOffset))
	q.builder.Offset(startOffset)
	return q
}

func (q *debugQueryBuilder) Or() query.Builder {
	q.builder.Or()
	q.isOr = true
	return q
}

func (q *debugQueryBuilder) Not() query.Builder {
	q.builder.Not()
	q.withNot = !q.withNot
	return q
}

func (q *debugQueryBuilder) OpenBracket() query.Builder {
	w := q.chunks[chunkWhere]
	if q.requireOp {
		if q.isOr {
			w.WriteString(" OR ")
		} else {
			w.WriteString(" AND ")
		}
	}
	w.WriteString("(")
	q.requireOp = false
	q.builder.OpenBracket()
	return q
}

func (q *debugQueryBuilder) CloseBracket() query.Builder {
	q.chunks[chunkWhere].WriteString(")")
	q.requireOp = true
	q.builder.CloseBracket()
	return q
}

func (q *debugQueryBuilder) SaveWhereForDump(field string, condition where.ComparatorType, value ...interface{}) {
	w := q.chunks[chunkWhere]
	if q.requireOp {
		if q.isOr {
			w.WriteString(" OR ")
		} else {
			w.WriteString(" AND ")
		}
	} else {
		q.requireOp = true
	}
	if q.withNot {
		w.WriteString("NOT ")
	}
	w.WriteString(field)
	switch condition {
	case where.EQ:
		w.WriteString(" = ")
		w.WriteString(fmt.Sprintf("%v", value[0]))
	case where.GT:
		w.WriteString(" > ")
		w.WriteString(fmt.Sprintf("%v", value[0]))
	case where.GE:
		w.WriteString(" >= ")
		w.WriteString(fmt.Sprintf("%v", value[0]))
	case where.LT:
		w.WriteString(" < ")
		w.WriteString(fmt.Sprintf("%v", value[0]))
	case where.LE:
		w.WriteString(" <= ")
		w.WriteString(fmt.Sprintf("%v", value[0]))
	case where.InArray:
		w.WriteString(" IN (")
		first := true
		for _, item := range value {
			if !first {
				w.WriteString(", ")
			} else {
				first = false
			}
			w.WriteString(fmt.Sprintf("%v", item))
		}
		w.WriteString(")")
	case where.Like:
		w.WriteString(" LIKE ")
		w.WriteString(fmt.Sprintf("\"%v\"", value[0]))
	case where.Regexp:
		w.WriteString(" REGEXP ")
		w.WriteString(fmt.Sprintf("%v", value[0]))
	case where.SetHas:
		w.WriteString(" SET_HAS ")
		w.WriteString(fmt.Sprintf("%v", value[0]))
	case where.MapHasValue:
		w.WriteString(" MAP_HAS_VALUE FIELD ")
		w.WriteString(value[0].(where.FieldComparator).GetField())
		w.WriteString(fmt.Sprintf(" COMPARE %v", value[0].(where.FieldComparator)))
	case where.MapHasKey:
		w.WriteString(" MAP_HAS_KEY ")
		w.WriteString(fmt.Sprintf("\"%v\"", value[0]))
	}
	q.withNot = false
	q.isOr = false
}

func (q *debugQueryBuilder) AddWhere(cmp where.FieldComparator) query.Builder {
	q.builder.AddWhere(cmp)
	q.withNot = false
	q.isOr = false
	return q
}

func (q *debugQueryBuilder) Where(getter *record.InterfaceGetter, condition where.ComparatorType, value ...interface{}) query.Builder {
	q.SaveWhereForDump(getter.Field, condition, value)
	q.builder.Where(getter, condition, value...)
	return q
}

func (q *debugQueryBuilder) WhereInt(getter *record.IntGetter, condition where.ComparatorType, value ...int) query.Builder {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.SaveWhereForDump(getter.Field, condition, items...)
	q.builder.WhereInt(getter, condition, value...)
	return q
}

func (q *debugQueryBuilder) WhereInt32(getter *record.Int32Getter, condition where.ComparatorType, value ...int32) query.Builder {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.SaveWhereForDump(getter.Field, condition, items...)
	q.builder.WhereInt32(getter, condition, value...)
	return q
}

func (q *debugQueryBuilder) WhereInt64(getter *record.Int64Getter, condition where.ComparatorType, value ...int64) query.Builder {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.SaveWhereForDump(getter.Field, condition, items...)
	q.builder.WhereInt64(getter, condition, value...)
	return q
}

func (q *debugQueryBuilder) WhereString(getter *record.StringGetter, condition where.ComparatorType, value ...string) query.Builder {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.SaveWhereForDump(getter.Field, condition, items...)
	q.builder.WhereString(getter, condition, value...)
	return q
}

func (q *debugQueryBuilder) WhereStringRegexp(getter *record.StringGetter, value *regexp.Regexp) query.Builder {
	q.SaveWhereForDump(getter.Field, where.Regexp, []interface{}{value})
	q.builder.WhereStringRegexp(getter, value)
	return q
}

func (q *debugQueryBuilder) WhereBool(getter *record.BoolGetter, condition where.ComparatorType, value ...bool) query.Builder {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.SaveWhereForDump(getter.Field, condition, items...)
	q.builder.WhereBool(getter, condition, value...)
	return q
}

func (q *debugQueryBuilder) WhereEnum8(getter *record.Enum8Getter, condition where.ComparatorType, value ...record.Enum8) query.Builder {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.SaveWhereForDump(getter.Field, condition, items...)
	q.builder.WhereEnum8(getter, condition, value...)
	return q
}

func (q *debugQueryBuilder) WhereEnum16(getter *record.Enum16Getter, condition where.ComparatorType, value ...record.Enum16) query.Builder {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.SaveWhereForDump(getter.Field, condition, items...)
	q.builder.WhereEnum16(getter, condition, value...)
	return q
}

func (q *debugQueryBuilder) WhereMap(getter *record.MapGetter, condition where.ComparatorType, value ...interface{}) query.Builder {
	items := make([]interface{}, len(value))
	copy(items, value)
	q.SaveWhereForDump(getter.Field, condition, items...)
	q.builder.WhereMap(getter, condition, value...)
	return q
}

func (q *debugQueryBuilder) WhereSet(getter *record.SetGetter, condition where.ComparatorType, value ...interface{}) query.Builder {
	items := make([]interface{}, len(value))
	copy(items, value)
	q.SaveWhereForDump(getter.Field, condition, items...)
	q.builder.WhereSet(getter, condition, value...)
	return q
}

func (q *debugQueryBuilder) MakeCopy() query.Builder {
	chunks := make(map[uint8]*strings.Builder, len(q.chunks))
	for key := range q.chunks {
		chunks[key] = &strings.Builder{}
		chunks[key].WriteString(q.chunks[key].String())
	}
	return &debugQueryBuilder{
		builder:   q.builder.MakeCopy(),
		chunks:    chunks,
		requireOp: q.requireOp,
		withNot:   q.withNot,
		isOr:      q.isOr,
	}
}

func (q *debugQueryBuilder) Sort(sortBy sort.By) query.Builder {
	w := q.chunks[chunkSort]
	if w.Len() > 0 {
		w.WriteString(", ")
	}
	w.WriteString(sortBy.String())
	q.builder.Sort(sortBy)
	return q
}

func (q *debugQueryBuilder) Query() query.Query {
	var result strings.Builder
	if q.chunks[chunkWhere].Len() > 0 {
		result.WriteString(" WHERE ")
		result.WriteString(q.chunks[chunkWhere].String())
	}
	if q.chunks[chunkSort].Len() > 0 {
		result.WriteString(" ORDER BY ")
		result.WriteString(q.chunks[chunkSort].String())
	}
	if q.chunks[chunkOffset].Len() > 0 {
		result.WriteString(" ")
		result.WriteString(q.chunks[chunkOffset].String())
	}
	if q.chunks[chunkLimit].Len() > 0 {
		result.WriteString(" ")
		result.WriteString(q.chunks[chunkLimit].String())
	}
	return &debugQuery{
		Query:     q.builder.Query(),
		queryDump: result.String(),
	}
}

func WrapQueryBuilderWithDebug(qb query.Builder) query.Builder {
	return &debugQueryBuilder{
		builder: qb,
		chunks: map[uint8]*strings.Builder{
			chunkLimit:  {},
			chunkOffset: {},
			chunkWhere:  {},
			chunkSort:   {},
		},
	}
}

type QueryBuilderConstructor func() query.Builder

func WrapCreateQueryBuilderWithDebug(constructor QueryBuilderConstructor) QueryBuilderConstructor {
	return func() query.Builder {
		return WrapQueryBuilderWithDebug(constructor())
	}
}
