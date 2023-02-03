package debug

import (
	"fmt"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"strconv"
	"strings"
)

type FieldComparatorDumper func(w *strings.Builder, cmp where.FieldComparator)

type QueryBuilderWithDumper interface {
	query.Builder
	SetFieldComparatorDumper(FieldComparatorDumper)
	Dump() string
}

const (
	chunkLimit uint8 = iota + 1
	chunkOffset
	chunkWhere
	chunkSort
)

var _ QueryBuilderWithDumper = (*debugQueryBuilder)(nil)

type debugQueryBuilder struct {
	chunks                map[uint8]*strings.Builder
	requireOp             bool
	withNot               bool
	isOr                  bool
	fieldComparatorDumper *FieldComparatorDumper
}

func (q *debugQueryBuilder) SetFieldComparatorDumper(dumper FieldComparatorDumper) {
	q.fieldComparatorDumper = &dumper
}

func (q *debugQueryBuilder) Limit(limitItems int) {
	w := q.chunks[chunkLimit]
	w.WriteString("LIMIT ")
	w.WriteString(strconv.Itoa(limitItems))
}

func (q *debugQueryBuilder) Offset(startOffset int) {
	w := q.chunks[chunkOffset]
	w.WriteString("OFFSET ")
	w.WriteString(strconv.Itoa(startOffset))
}

func (q *debugQueryBuilder) Not() {
	q.withNot = !q.withNot
}

func (q *debugQueryBuilder) Or() {
	q.isOr = true
}

func (q *debugQueryBuilder) OpenBracket() {
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
}

func (q *debugQueryBuilder) CloseBracket() {
	q.chunks[chunkWhere].WriteString(")")
	q.requireOp = true
}
func (q *debugQueryBuilder) AddWhere(cmp where.FieldComparator) {
	q.saveFieldComparatorForDump(cmp)
	q.withNot = false
	q.isOr = false
	q.requireOp = true
}

func (q *debugQueryBuilder) Sort(sortBy sort.ByWithOrder) {
	w := q.chunks[chunkSort]
	if w.Len() > 0 {
		w.WriteString(", ")
	}
	w.WriteString(sortBy.String())
}

func (q *debugQueryBuilder) OnIteration(_ func(item record.Record)) {
}

func (q *debugQueryBuilder) Append(options ...query.BuilderOption) {
	for _, opt := range options {
		opt.Apply(q)
	}
}

func (q *debugQueryBuilder) MakeCopy() query.Builder {
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

func (q *debugQueryBuilder) Query() query.Query {
	return nil
}

func (q *debugQueryBuilder) saveFieldComparatorForDump(cmp where.FieldComparator) {
	w := q.chunks[chunkWhere]
	if q.requireOp {
		if q.isOr {
			w.WriteString(" OR ")
		} else {
			w.WriteString(" AND ")
		}
	}
	if q.withNot {
		w.WriteString("NOT ")
	}
	w.WriteString(cmp.GetField().String())
	switch cmp.GetType() {
	case where.EQ:
		w.WriteString(" = ")
		w.WriteString(fmt.Sprintf("%v", cmp.ValueAt(0)))
	case where.GT:
		w.WriteString(" > ")
		w.WriteString(fmt.Sprintf("%v", cmp.ValueAt(0)))
	case where.GE:
		w.WriteString(" >= ")
		w.WriteString(fmt.Sprintf("%v", cmp.ValueAt(0)))
	case where.LT:
		w.WriteString(" < ")
		w.WriteString(fmt.Sprintf("%v", cmp.ValueAt(0)))
	case where.LE:
		w.WriteString(" <= ")
		w.WriteString(fmt.Sprintf("%v", cmp.ValueAt(0)))
	case where.InArray:
		w.WriteString(" IN (")
		for i := 0; i < cmp.ValuesCount(); i++ {
			if 0 != i {
				w.WriteString(", ")
			}
			w.WriteString(fmt.Sprintf("%v", cmp.ValueAt(i)))
		}
		w.WriteString(")")
	case where.Like:
		w.WriteString(" LIKE ")
		w.WriteString(fmt.Sprintf("\"%v\"", cmp.ValueAt(0)))
	case where.Regexp:
		w.WriteString(" REGEXP ")
		w.WriteString(fmt.Sprintf("%v", cmp.ValueAt(0)))
	case where.SetHas:
		w.WriteString(" SET_HAS ")
		w.WriteString(fmt.Sprintf("%v", cmp.ValueAt(0)))
	case where.MapHasValue:
		w.WriteString(" MAP_HAS_VALUE FIELD ")
		mapCmp := cmp.ValueAt(0).(where.FieldComparator)
		w.WriteString(mapCmp.GetField().String())
		w.WriteString(fmt.Sprintf(" COMPARE %v", mapCmp))
	case where.MapHasKey:
		w.WriteString(" MAP_HAS_KEY ")
		w.WriteString(fmt.Sprintf("\"%v\"", cmp.ValueAt(0)))
	default:
		if nil == q.fieldComparatorDumper {
			w.WriteString(fmt.Sprintf(" (ComparatorType(%d) ", cmp.GetType()))
			for i := 0; i < cmp.ValuesCount(); i++ {
				if 0 != i {
					w.WriteString(" ")
				}
				w.WriteString(fmt.Sprintf("%v", cmp.ValueAt(i)))
			}
			w.WriteString(")")
		} else {
			(*q.fieldComparatorDumper)(w, cmp)
		}
	}
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

func CreateDebugQueryBuilder(options ...query.BuilderOption) QueryBuilderWithDumper {
	qb := &debugQueryBuilder{
		chunks: map[uint8]*strings.Builder{
			chunkLimit:  {},
			chunkOffset: {},
			chunkWhere:  {},
			chunkSort:   {},
		},
	}
	for _, opt := range options {
		opt.Apply(qb)
	}
	return qb
}
