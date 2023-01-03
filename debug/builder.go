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

type QueryBuilderDumper interface {
	SaveWhereForDump(field string, condition where.ComparatorType, value ...interface{})
	Dump() string
}

type BaseQueryBuilderWithDumper interface {
	query.BaseQueryBuilder
	QueryBuilderDumper
}

const (
	chunkLimit uint8 = iota + 1
	chunkOffset
	chunkWhere
	chunkSort
)

var _ BaseQueryBuilderWithDumper = (*debugQueryBuilder)(nil)

type debugQueryBuilder struct {
	chunks    map[uint8]*strings.Builder
	requireOp bool
	withNot   bool
	isOr      bool
}

func (q *debugQueryBuilder) OnIteration(_ func(item record.Record)) query.BaseQueryBuilder {
	return q
}

func (q *debugQueryBuilder) Limit(limitItems int) query.BaseQueryBuilder {
	w := q.chunks[chunkLimit]
	w.WriteString("LIMIT ")
	w.WriteString(strconv.Itoa(limitItems))
	return q
}

func (q *debugQueryBuilder) Offset(startOffset int) query.BaseQueryBuilder {
	w := q.chunks[chunkOffset]
	w.WriteString("OFFSET ")
	w.WriteString(strconv.Itoa(startOffset))
	return q
}

func (q *debugQueryBuilder) Or() query.BaseQueryBuilder {
	q.isOr = true
	return q
}

func (q *debugQueryBuilder) Not() query.BaseQueryBuilder {
	q.withNot = !q.withNot
	return q
}

func (q *debugQueryBuilder) OpenBracket() query.BaseQueryBuilder {
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
	return q
}

func (q *debugQueryBuilder) CloseBracket() query.BaseQueryBuilder {
	q.chunks[chunkWhere].WriteString(")")
	q.requireOp = true
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

func (q *debugQueryBuilder) AddWhere(_ where.FieldComparator) query.BaseQueryBuilder {
	q.withNot = false
	q.isOr = false
	return q
}

func (q *debugQueryBuilder) Where(getter *record.InterfaceGetter, condition where.ComparatorType, value ...interface{}) query.BaseQueryBuilder {
	q.SaveWhereForDump(getter.Field, condition, value)
	return q
}

func (q *debugQueryBuilder) WhereInt(getter *record.IntGetter, condition where.ComparatorType, value ...int) query.BaseQueryBuilder {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.SaveWhereForDump(getter.Field, condition, items...)
	return q
}

func (q *debugQueryBuilder) WhereInt32(getter *record.Int32Getter, condition where.ComparatorType, value ...int32) query.BaseQueryBuilder {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.SaveWhereForDump(getter.Field, condition, items...)
	return q
}

func (q *debugQueryBuilder) WhereInt64(getter *record.Int64Getter, condition where.ComparatorType, value ...int64) query.BaseQueryBuilder {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.SaveWhereForDump(getter.Field, condition, items...)
	return q
}

func (q *debugQueryBuilder) WhereString(getter *record.StringGetter, condition where.ComparatorType, value ...string) query.BaseQueryBuilder {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.SaveWhereForDump(getter.Field, condition, items...)
	return q
}

func (q *debugQueryBuilder) WhereStringRegexp(getter *record.StringGetter, value *regexp.Regexp) query.BaseQueryBuilder {
	q.SaveWhereForDump(getter.Field, where.Regexp, []interface{}{value})
	return q
}

func (q *debugQueryBuilder) WhereBool(getter *record.BoolGetter, condition where.ComparatorType, value ...bool) query.BaseQueryBuilder {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.SaveWhereForDump(getter.Field, condition, items...)
	return q
}

func (q *debugQueryBuilder) WhereEnum8(getter *record.Enum8Getter, condition where.ComparatorType, value ...record.Enum8) query.BaseQueryBuilder {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.SaveWhereForDump(getter.Field, condition, items...)
	return q
}

func (q *debugQueryBuilder) WhereEnum16(getter *record.Enum16Getter, condition where.ComparatorType, value ...record.Enum16) query.BaseQueryBuilder {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.SaveWhereForDump(getter.Field, condition, items...)
	return q
}

func (q *debugQueryBuilder) WhereMap(getter *record.MapGetter, condition where.ComparatorType, value ...interface{}) query.BaseQueryBuilder {
	items := make([]interface{}, len(value))
	copy(items, value)
	q.SaveWhereForDump(getter.Field, condition, items...)
	return q
}

func (q *debugQueryBuilder) WhereSet(getter *record.SetGetter, condition where.ComparatorType, value ...interface{}) query.BaseQueryBuilder {
	items := make([]interface{}, len(value))
	copy(items, value)
	q.SaveWhereForDump(getter.Field, condition, items...)
	return q
}

func (q *debugQueryBuilder) MakeCopy() query.BaseQueryBuilder {
	chunks := make(map[uint8]*strings.Builder, len(q.chunks))
	for key := range q.chunks {
		chunks[key] = &strings.Builder{}
		chunks[key].WriteString(q.chunks[key].String())
	}
	return &debugQueryBuilder{
		chunks:    chunks,
		requireOp: q.requireOp,
		withNot:   q.withNot,
		isOr:      q.isOr,
	}
}

func (q *debugQueryBuilder) Sort(sortBy sort.By) query.BaseQueryBuilder {
	w := q.chunks[chunkSort]
	if w.Len() > 0 {
		w.WriteString(", ")
	}
	w.WriteString(sortBy.String())
	return q
}

func (q *debugQueryBuilder) Query() query.Query {
	return nil
}

func (q *debugQueryBuilder) Dump() string {
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
	return result.String()
}

func CreateDebugQueryBuilder() BaseQueryBuilderWithDumper {
	return &debugQueryBuilder{
		chunks: map[uint8]*strings.Builder{
			chunkLimit:  {},
			chunkOffset: {},
			chunkWhere:  {},
			chunkSort:   {},
		},
	}
}
