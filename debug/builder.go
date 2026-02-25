package debug

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
)

type FieldComparatorDumper[R record.Record] func(builder *strings.Builder, cmp where.FieldComparator[R])

const (
	chunkLimit uint8 = iota + 1
	chunkOffset
	chunkWhere
	chunkSort
)

type debugQueryBuilder[R record.Record] struct {
	chunks                map[uint8]*strings.Builder
	requireOp             bool
	withNot               bool
	isOr                  bool
	fieldComparatorDumper *FieldComparatorDumper[R]
}

func (q *debugQueryBuilder[R]) setFieldComparatorDumper(dumper FieldComparatorDumper[R]) {
	q.fieldComparatorDumper = &dumper
}

func (q *debugQueryBuilder[R]) limit(limitItems int) {
	chunk := q.chunks[chunkLimit]
	chunk.WriteString("LIMIT ")
	chunk.WriteString(strconv.Itoa(limitItems))
}

func (q *debugQueryBuilder[R]) offset(startOffset int) {
	chunk := q.chunks[chunkOffset]
	chunk.WriteString("OFFSET ")
	chunk.WriteString(strconv.Itoa(startOffset))
}

func (q *debugQueryBuilder[R]) not() {
	q.withNot = !q.withNot
}

func (q *debugQueryBuilder[R]) or() {
	q.isOr = true
}

func (q *debugQueryBuilder[R]) openBracket() {
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

func (q *debugQueryBuilder[R]) closeBracket() {
	q.chunks[chunkWhere].WriteString(")")
	q.requireOp = true
}

func (q *debugQueryBuilder[R]) addWhere(cmp where.FieldComparator[R]) {
	q.saveFieldComparatorForDump(cmp)
	q.withNot = false
	q.isOr = false
	q.requireOp = true
}

func (q *debugQueryBuilder[R]) sort(sortBy sort.ByWithOrder[R]) {
	chunk := q.chunks[chunkSort]
	if chunk.Len() > 0 {
		chunk.WriteString(", ")
	}

	chunk.WriteString(sortBy.String())
}
func (q *debugQueryBuilder[R]) makeCopy() *debugQueryBuilder[R] {
	chunks := make(map[uint8]*strings.Builder, len(q.chunks))
	for key := range q.chunks {
		chunks[key] = &strings.Builder{}
		chunks[key].WriteString(q.chunks[key].String())
	}

	return &debugQueryBuilder[R]{
		chunks:                chunks,
		requireOp:             q.requireOp,
		withNot:               q.withNot,
		isOr:                  q.isOr,
		fieldComparatorDumper: q.fieldComparatorDumper,
	}
}

func (q *debugQueryBuilder[R]) saveFieldComparatorForDump(cmp where.FieldComparator[R]) { //nolint:funlen,cyclop
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
		fmt.Fprintf(chunk, "%v", cmp.ValueAt(0))
	case where.GT:
		chunk.WriteString(" > ")
		fmt.Fprintf(chunk, "%v", cmp.ValueAt(0))
	case where.GE:
		chunk.WriteString(" >= ")
		fmt.Fprintf(chunk, "%v", cmp.ValueAt(0))
	case where.LT:
		chunk.WriteString(" < ")
		fmt.Fprintf(chunk, "%v", cmp.ValueAt(0))
	case where.LE:
		chunk.WriteString(" <= ")
		fmt.Fprintf(chunk, "%v", cmp.ValueAt(0))
	case where.InArray:
		chunk.WriteString(" IN (")

		for i := range cmp.ValuesCount() {
			if i != 0 {
				chunk.WriteString(", ")
			}

			fmt.Fprintf(chunk, "%v", cmp.ValueAt(i))
		}

		chunk.WriteString(")")
	case where.Like:
		chunk.WriteString(" LIKE ")
		fmt.Fprintf(chunk, "\"%v\"", cmp.ValueAt(0))
	case where.Regexp:
		chunk.WriteString(" REGEXP ")
		fmt.Fprintf(chunk, "%v", cmp.ValueAt(0))
	case where.SetHas:
		chunk.WriteString(" SET_HAS ")
		fmt.Fprintf(chunk, "%v", cmp.ValueAt(0))
	case where.MapHasValue:
		chunk.WriteString(" MAP_HAS_VALUE FIELD ")

		mapCmp := cmp.ValueAt(0).(where.FieldComparator[R])
		chunk.WriteString(mapCmp.GetField().String())
		fmt.Fprintf(chunk, " COMPARE %v", mapCmp)
	case where.MapHasKey:
		chunk.WriteString(" MAP_HAS_KEY ")
		fmt.Fprintf(chunk, "\"%v\"", cmp.ValueAt(0))
	default:
		if nil == q.fieldComparatorDumper {
			fmt.Fprintf(chunk, " (ComparatorType(%d) ", cmp.GetType())

			for i := range cmp.ValuesCount() {
				if i != 0 {
					chunk.WriteString(" ")
				}

				fmt.Fprintf(chunk, "%v", cmp.ValueAt(i))
			}

			chunk.WriteString(")")
		} else {
			(*q.fieldComparatorDumper)(chunk, cmp)
		}
	}
}

func (q *debugQueryBuilder[R]) Dump() string {
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

func newDebugQueryBuilder[R record.Record]() *debugQueryBuilder[R] {
	return &debugQueryBuilder[R]{
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
}
