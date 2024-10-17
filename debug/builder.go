package debug

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
)

type FieldComparatorDumper func(builder *strings.Builder, cmp where.FieldComparator)

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
	chunk := q.chunks[chunkLimit]
	chunk.WriteString("LIMIT ")
	chunk.WriteString(strconv.Itoa(limitItems))
}

func (q *debugQueryBuilder) Offset(startOffset int) {
	chunk := q.chunks[chunkOffset]
	chunk.WriteString("OFFSET ")
	chunk.WriteString(strconv.Itoa(startOffset))
}

func (q *debugQueryBuilder) Not() {
	q.withNot = !q.withNot
}

func (q *debugQueryBuilder) Or() {
	q.isOr = true
}

func (q *debugQueryBuilder) OpenBracket() {
	chunk := q.chunks[chunkWhere]
	if q.requireOp {
		if q.isOr {
			chunk.WriteString(" OR ")
		} else {
			chunk.WriteString(" AND ")
		}
	}
	chunk.WriteString("(")
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
	chunk := q.chunks[chunkSort]
	if chunk.Len() > 0 {
		chunk.WriteString(", ")
	}
	chunk.WriteString(sortBy.String())
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
		chunks:                chunks,
		requireOp:             q.requireOp,
		withNot:               q.withNot,
		isOr:                  q.isOr,
		fieldComparatorDumper: q.fieldComparatorDumper,
	}
}

func (q *debugQueryBuilder) Query() query.Query {
	return nil
}

func (q *debugQueryBuilder) saveFieldComparatorForDump(cmp where.FieldComparator) { //nolint:funlen,cyclop
	chunk := q.chunks[chunkWhere]
	if q.requireOp {
		if q.isOr {
			chunk.WriteString(" OR ")
		} else {
			chunk.WriteString(" AND ")
		}
	}
	if q.withNot {
		chunk.WriteString("NOT ")
	}
	chunk.WriteString(cmp.GetField().String())
	switch cmp.GetType() {
	case where.EQ:
		chunk.WriteString(" = ")
		chunk.WriteString(fmt.Sprintf("%v", cmp.ValueAt(0)))
	case where.GT:
		chunk.WriteString(" > ")
		chunk.WriteString(fmt.Sprintf("%v", cmp.ValueAt(0)))
	case where.GE:
		chunk.WriteString(" >= ")
		chunk.WriteString(fmt.Sprintf("%v", cmp.ValueAt(0)))
	case where.LT:
		chunk.WriteString(" < ")
		chunk.WriteString(fmt.Sprintf("%v", cmp.ValueAt(0)))
	case where.LE:
		chunk.WriteString(" <= ")
		chunk.WriteString(fmt.Sprintf("%v", cmp.ValueAt(0)))
	case where.InArray:
		chunk.WriteString(" IN (")
		for i := 0; i < cmp.ValuesCount(); i++ {
			if i != 0 {
				chunk.WriteString(", ")
			}
			chunk.WriteString(fmt.Sprintf("%v", cmp.ValueAt(i)))
		}
		chunk.WriteString(")")
	case where.Like:
		chunk.WriteString(" LIKE ")
		chunk.WriteString(fmt.Sprintf("\"%v\"", cmp.ValueAt(0)))
	case where.Regexp:
		chunk.WriteString(" REGEXP ")
		chunk.WriteString(fmt.Sprintf("%v", cmp.ValueAt(0)))
	case where.SetHas:
		chunk.WriteString(" SET_HAS ")
		chunk.WriteString(fmt.Sprintf("%v", cmp.ValueAt(0)))
	case where.MapHasValue:
		chunk.WriteString(" MAP_HAS_VALUE FIELD ")
		mapCmp := cmp.ValueAt(0).(where.FieldComparator)
		chunk.WriteString(mapCmp.GetField().String())
		chunk.WriteString(fmt.Sprintf(" COMPARE %v", mapCmp))
	case where.MapHasKey:
		chunk.WriteString(" MAP_HAS_KEY ")
		chunk.WriteString(fmt.Sprintf("\"%v\"", cmp.ValueAt(0)))
	default:
		if nil == q.fieldComparatorDumper {
			chunk.WriteString(fmt.Sprintf(" (ComparatorType(%d) ", cmp.GetType()))
			for i := 0; i < cmp.ValuesCount(); i++ {
				if i != 0 {
					chunk.WriteString(" ")
				}
				chunk.WriteString(fmt.Sprintf("%v", cmp.ValueAt(i)))
			}
			chunk.WriteString(")")
		} else {
			(*q.fieldComparatorDumper)(chunk, cmp)
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
	debugQB := &debugQueryBuilder{
		chunks: map[uint8]*strings.Builder{
			chunkLimit:  {},
			chunkOffset: {},
			chunkWhere:  {},
			chunkSort:   {},
		},
		requireOp:             false,
		withNot:               false,
		isOr:                  false,
		fieldComparatorDumper: nil,
	}
	for _, opt := range options {
		opt.Apply(debugQB)
	}
	return debugQB
}
