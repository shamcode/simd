package namespace

import (
	"context"
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"regexp"
	"strconv"
	"strings"
)

const (
	chunkLimit uint8 = iota + 1
	chunkOffset
	chunkWhere
	chunkSort
)

var (
	_ Query = (*debugExecutor)(nil)
)

type debugExecutor struct {
	executor  Query
	chunks    map[uint8]*strings.Builder
	requireOp bool
	withNot   bool
	isOr      bool
	dump      func(string)
}

func (q *debugExecutor) OnIteration(cb func(item record.Record)) Query {
	q.executor.OnIteration(cb)
	return q
}

func (q *debugExecutor) Limit(limitItems int) Query {
	w := q.chunks[chunkLimit]
	w.WriteString("LIMIT ")
	w.WriteString(strconv.Itoa(limitItems))
	q.executor.Limit(limitItems)
	return q
}

func (q *debugExecutor) Offset(startOffset int) Query {
	w := q.chunks[chunkOffset]
	w.WriteString("OFFSET ")
	w.WriteString(strconv.Itoa(startOffset))
	q.executor.Offset(startOffset)
	return q
}

func (q *debugExecutor) Or() Query {
	q.executor.Or()
	q.isOr = true
	return q
}

func (q *debugExecutor) Not() Query {
	q.executor.Not()
	q.withNot = !q.withNot
	return q
}

func (q *debugExecutor) OpenBracket() Query {
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
	q.executor.OpenBracket()
	return q
}

func (q *debugExecutor) CloseBracket() Query {
	q.chunks[chunkWhere].WriteString(")")
	q.requireOp = true
	q.executor.CloseBracket()
	return q
}

func (q *debugExecutor) addWhere(field string, condition where.ComparatorType, value ...interface{}) {
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

func (q *debugExecutor) AddWhere(cmp where.FieldComparator) Query {
	q.executor.AddWhere(cmp)
	q.withNot = false
	q.isOr = false
	return q
}

func (q *debugExecutor) Where(getter *record.InterfaceGetter, condition where.ComparatorType, value ...interface{}) Query {
	q.addWhere(getter.Field, condition, value)
	q.executor.Where(getter, condition, value...)
	return q
}

func (q *debugExecutor) WhereInt(getter *record.IntGetter, condition where.ComparatorType, value ...int) Query {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.addWhere(getter.Field, condition, items...)
	q.executor.WhereInt(getter, condition, value...)
	return q
}

func (q *debugExecutor) WhereInt32(getter *record.Int32Getter, condition where.ComparatorType, value ...int32) Query {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.addWhere(getter.Field, condition, items...)
	q.executor.WhereInt32(getter, condition, value...)
	return q
}

func (q *debugExecutor) WhereInt64(getter *record.Int64Getter, condition where.ComparatorType, value ...int64) Query {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.addWhere(getter.Field, condition, items...)
	q.executor.WhereInt64(getter, condition, value...)
	return q
}

func (q *debugExecutor) WhereString(getter *record.StringGetter, condition where.ComparatorType, value ...string) Query {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.addWhere(getter.Field, condition, items...)
	q.executor.WhereString(getter, condition, value...)
	return q
}

func (q *debugExecutor) WhereStringRegexp(getter *record.StringGetter, value *regexp.Regexp) Query {
	q.addWhere(getter.Field, where.Regexp, []interface{}{value})
	q.executor.WhereStringRegexp(getter, value)
	return q
}

func (q *debugExecutor) WhereBool(getter *record.BoolGetter, condition where.ComparatorType, value ...bool) Query {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.addWhere(getter.Field, condition, items...)
	q.executor.WhereBool(getter, condition, value...)
	return q
}

func (q *debugExecutor) WhereEnum8(getter *record.Enum8Getter, condition where.ComparatorType, value ...record.Enum8) Query {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.addWhere(getter.Field, condition, items...)
	q.executor.WhereEnum8(getter, condition, value...)
	return q
}

func (q *debugExecutor) WhereEnum16(getter *record.Enum16Getter, condition where.ComparatorType, value ...record.Enum16) Query {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.addWhere(getter.Field, condition, items...)
	q.executor.WhereEnum16(getter, condition, value...)
	return q
}

func (q *debugExecutor) WhereMap(getter *record.MapGetter, condition where.ComparatorType, value ...interface{}) Query {
	items := make([]interface{}, len(value))
	copy(items, value)
	q.addWhere(getter.Field, condition, items...)
	q.executor.WhereMap(getter, condition, value...)
	return q
}

func (q *debugExecutor) WhereSet(getter *record.SetGetter, condition where.ComparatorType, value ...interface{}) Query {
	items := make([]interface{}, len(value))
	copy(items, value)
	q.addWhere(getter.Field, condition, items...)
	q.executor.WhereSet(getter, condition, value...)
	return q
}

func (q *debugExecutor) WhereStringsSet(getter *record.StringsSetGetter, condition where.ComparatorType, value ...string) Query {
	items := make([]interface{}, len(value))
	for i := range value {
		items[i] = value[i]
	}
	q.addWhere(getter.Field, condition, items...)
	q.executor.WhereStringsSet(getter, condition, value...)
	return q
}

func (q *debugExecutor) MakeCopy() Query {
	chunks := make(map[uint8]*strings.Builder, len(q.chunks))
	for key := range q.chunks {
		chunks[key] = &strings.Builder{}
		chunks[key].WriteString(q.chunks[key].String())
	}
	return &debugExecutor{
		chunks:    chunks,
		dump:      q.dump,
		executor:  q.executor.MakeCopy(),
		requireOp: q.requireOp,
	}
}

func (q *debugExecutor) Sort(sortBy sort.By) Query {
	w := q.chunks[chunkSort]
	if w.Len() > 0 {
		w.WriteString(", ")
	}
	w.WriteString(sortBy.String())
	q.executor.Sort(sortBy)
	return q
}

func (q *debugExecutor) FetchTotal(ctx context.Context) (int, error) {
	q.dumpQuery(true)
	return q.executor.FetchTotal(ctx)

}

func (q *debugExecutor) FetchAll(ctx context.Context) (Iterator, error) {
	q.dumpQuery(false)
	return q.executor.FetchAll(ctx)
}

func (q *debugExecutor) FetchAllAndTotal(ctx context.Context) (Iterator, int, error) {
	q.dumpQuery(false)
	return q.executor.FetchAllAndTotal(ctx)
}

func (q *debugExecutor) dumpQuery(onlyTotal bool) {
	var result strings.Builder
	result.WriteString("SELECT ")
	if !onlyTotal {
		result.WriteString("*, ")
	}
	result.WriteString("COUNT(*) WHERE ")
	result.WriteString(q.chunks[chunkWhere].String())
	if !onlyTotal {
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
	}
	q.dump(result.String())
}

func WrapQueryWithDebug(executor Query, dump func(string)) Query {
	return &debugExecutor{
		chunks: map[uint8]*strings.Builder{
			chunkLimit:  {},
			chunkOffset: {},
			chunkWhere:  {},
			chunkSort:   {},
		},
		executor: executor,
		dump:     dump,
	}
}
