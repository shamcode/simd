package query

import (
	"github.com/hashicorp/go-multierror"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
	"regexp"
)

// Builder is a helper for build Query
type Builder interface {
	MakeCopy() Builder

	Limit(limitItems int) Builder
	Offset(startOffset int) Builder

	Not() Builder
	Or() Builder

	OpenBracket() Builder
	CloseBracket() Builder

	AddWhere(cmp where.FieldComparator) Builder

	Where(getter *record.InterfaceGetter, condition where.ComparatorType, values ...interface{}) Builder
	WhereInt(getter *record.IntGetter, condition where.ComparatorType, values ...int) Builder
	WhereInt32(getter *record.Int32Getter, condition where.ComparatorType, values ...int32) Builder
	WhereInt64(getter *record.Int64Getter, condition where.ComparatorType, values ...int64) Builder
	WhereString(getter *record.StringGetter, condition where.ComparatorType, values ...string) Builder
	WhereStringRegexp(getter *record.StringGetter, value *regexp.Regexp) Builder
	WhereBool(getter *record.BoolGetter, condition where.ComparatorType, values ...bool) Builder
	WhereEnum8(getter *record.Enum8Getter, condition where.ComparatorType, values ...record.Enum8) Builder
	WhereEnum16(getter *record.Enum16Getter, condition where.ComparatorType, values ...record.Enum16) Builder
	WhereMap(getter *record.MapGetter, condition where.ComparatorType, values ...interface{}) Builder
	WhereSet(getter *record.SetGetter, condition where.ComparatorType, values ...interface{}) Builder

	Sort(by sort.By) Builder

	// OnIteration registers a callback to be called for each record before sorting and applying offset/limits
	// but after applying WHERE conditions
	OnIteration(cb func(item record.Record)) Builder

	Query() Query
}

var _ Builder = (*queryBuilder)(nil)

type queryBuilder struct {
	limitItems   int
	startOffset  int
	withLimit    bool
	withNot      bool
	isOr         bool
	conditionSet bool
	bracketLevel int
	where        where.Conditions
	sortBy       []sort.By
	onIteration  *func(item record.Record)
	error        *multierror.Error // TODO: wait go1.20 https://go-review.googlesource.com/c/go/+/432898/11/src/errors/join.go
}

func (q *queryBuilder) Sorting() []sort.By {
	return q.sortBy
}

func (q *queryBuilder) Conditions() where.Conditions {
	return q.where
}

func (q *queryBuilder) OnIterationCallback() *func(item record.Record) {
	return q.onIteration
}

func (q *queryBuilder) Error() error {
	return q.error.ErrorOrNil()
}

func (q *queryBuilder) AddWhere(cmp where.FieldComparator) Builder {
	q.where = append(q.where, where.Condition{
		WithNot:      q.withNot,
		IsOr:         q.isOr,
		BracketLevel: 1 + q.bracketLevel,
		Cmp:          cmp,
	})
	q.withNot = false
	q.isOr = false
	q.conditionSet = true
	return q
}

func (q *queryBuilder) MakeCopy() Builder {
	cpy := &queryBuilder{
		limitItems:   q.limitItems,
		startOffset:  q.startOffset,
		withLimit:    q.withLimit,
		withNot:      q.withNot,
		isOr:         q.isOr,
		bracketLevel: q.bracketLevel,
		where:        make(where.Conditions, len(q.where), cap(q.where)),
		sortBy:       make([]sort.By, len(q.sortBy), cap(q.sortBy)),
		onIteration:  q.onIteration,
	}
	for i, item := range q.where {
		cpy.where[i] = item
	}
	for i, item := range q.sortBy {
		cpy.sortBy[i] = item
	}
	return cpy
}

func (q *queryBuilder) OnIteration(cb func(item record.Record)) Builder {
	q.onIteration = &cb
	return q
}

func (q *queryBuilder) Limit(limitItems int) Builder {
	q.limitItems = limitItems
	q.withLimit = true
	return q
}

func (q *queryBuilder) Offset(startOffset int) Builder {
	q.startOffset = startOffset
	return q
}

func (q *queryBuilder) Or() Builder {
	q.isOr = true
	if !q.conditionSet {
		q.error = multierror.Append(q.error, ErrOrBeforeAnyConditions)
	}
	return q
}

func (q *queryBuilder) Not() Builder {
	q.withNot = !q.withNot
	return q
}

func (q *queryBuilder) OpenBracket() Builder {
	if q.withNot {
		q.error = multierror.Append(q.error, ErrNotOpenBracket)
	}
	q.conditionSet = false
	q.bracketLevel += 1
	return q
}

func (q *queryBuilder) CloseBracket() Builder {
	q.bracketLevel -= 1
	if -1 == q.bracketLevel {
		q.error = multierror.Append(q.error, ErrCloseBracketWithoutOpen)
	}
	q.conditionSet = true
	return q
}

func (q *queryBuilder) Where(getter *record.InterfaceGetter, condition where.ComparatorType, value ...interface{}) Builder {
	return q.AddWhere(comparators.InterfaceFieldComparator{
		Cmp:    condition,
		Getter: getter,
		Value:  value,
	})
}

func (q *queryBuilder) WhereInt(getter *record.IntGetter, condition where.ComparatorType, value ...int) Builder {
	return q.AddWhere(comparators.IntFieldComparator{
		Cmp:    condition,
		Getter: getter,
		Value:  value,
	})
}

func (q *queryBuilder) WhereInt32(getter *record.Int32Getter, condition where.ComparatorType, value ...int32) Builder {
	return q.AddWhere(comparators.Int32FieldComparator{
		Cmp:    condition,
		Getter: getter,
		Value:  value,
	})
}

func (q *queryBuilder) WhereInt64(getter *record.Int64Getter, condition where.ComparatorType, value ...int64) Builder {
	return q.AddWhere(comparators.Int64FieldComparator{
		Cmp:    condition,
		Getter: getter,
		Value:  value,
	})
}

func (q *queryBuilder) WhereString(getter *record.StringGetter, condition where.ComparatorType, value ...string) Builder {
	return q.AddWhere(comparators.StringFieldComparator{
		Cmp:    condition,
		Getter: getter,
		Value:  value,
	})
}

func (q *queryBuilder) WhereStringRegexp(getter *record.StringGetter, value *regexp.Regexp) Builder {
	return q.AddWhere(comparators.StringFieldRegexpComparator{
		Cmp:    where.Regexp,
		Getter: getter,
		Value:  value,
	})
}

func (q *queryBuilder) WhereBool(getter *record.BoolGetter, condition where.ComparatorType, value ...bool) Builder {
	return q.AddWhere(comparators.BoolFieldComparator{
		Cmp:    condition,
		Getter: getter,
		Value:  value,
	})
}

func (q *queryBuilder) WhereEnum8(getter *record.Enum8Getter, condition where.ComparatorType, value ...record.Enum8) Builder {
	return q.AddWhere(comparators.Enum8FieldComparator{
		Cmp:    condition,
		Getter: getter,
		Value:  value,
	})
}

func (q *queryBuilder) WhereEnum16(getter *record.Enum16Getter, condition where.ComparatorType, value ...record.Enum16) Builder {
	return q.AddWhere(comparators.Enum16FieldComparator{
		Cmp:    condition,
		Getter: getter,
		Value:  value,
	})
}

func (q *queryBuilder) WhereMap(getter *record.MapGetter, condition where.ComparatorType, value ...interface{}) Builder {
	return q.AddWhere(comparators.MapFieldComparator{
		Cmp:    condition,
		Getter: getter,
		Value:  value,
	})
}

func (q *queryBuilder) WhereSet(getter *record.SetGetter, condition where.ComparatorType, value ...interface{}) Builder {
	return q.AddWhere(comparators.SetFieldComparator{
		Cmp:    condition,
		Getter: getter,
		Value:  value,
	})
}

func (q *queryBuilder) Sort(sortBy sort.By) Builder {
	q.sortBy = append(q.sortBy, sortBy)
	return q
}

func (q *queryBuilder) Query() Query {
	if q.bracketLevel > 0 {
		q.error = multierror.Append(q.error, ErrInvalidBracketBalance)
	}
	return &query{
		offset:              q.startOffset,
		limit:               q.limitItems,
		withLimit:           q.withLimit,
		conditions:          q.where,
		sorting:             q.sortBy,
		onIterationCallback: q.onIteration,
		error:               q.error.ErrorOrNil(),
	}
}

func NewBuilder() Builder {
	return &queryBuilder{}
}
