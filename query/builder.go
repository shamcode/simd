package query

import (
	"github.com/hashicorp/go-multierror"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
)

type Builder interface { //nolint:interfacebloat
	Limit(limitItems int)
	Offset(startOffset int)
	Not()
	Or()
	OpenBracket()
	CloseBracket()
	AddWhere(cmp where.FieldComparator)
	Sort(by sort.ByWithOrder)

	// OnIteration registers a callback to be called for each record before sorting and applying offset/limits
	// but after applying WHERE conditions
	OnIteration(cb func(item record.Record))

	// Append apply new options to builder
	Append(options ...BuilderOption)

	MakeCopy() Builder

	// Query return build Query
	Query() Query
}

type BuilderOption interface {
	Apply(b Builder)
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
	sortBy       []sort.ByWithOrder
	onIteration  *func(item record.Record)
	// TODO: wait go1.20 https://go-review.googlesource.com/c/go/+/432898/11/src/errors/join.go
	error *multierror.Error
}

func (qb *queryBuilder) Limit(limitItems int) {
	qb.limitItems = limitItems
	qb.withLimit = true
}

func (qb *queryBuilder) Offset(startOffset int) {
	qb.startOffset = startOffset
}

func (qb *queryBuilder) Not() {
	qb.withNot = !qb.withNot
}

func (qb *queryBuilder) Or() {
	qb.isOr = true
	if !qb.conditionSet {
		qb.error = multierror.Append(qb.error, ErrOrBeforeAnyConditions)
	}
}

func (qb *queryBuilder) OpenBracket() {
	if qb.withNot {
		qb.error = multierror.Append(qb.error, ErrNotOpenBracket)
	}
	qb.conditionSet = false
	qb.bracketLevel += 1
}

func (qb *queryBuilder) CloseBracket() {
	qb.bracketLevel -= 1
	if qb.bracketLevel == -1 {
		qb.error = multierror.Append(qb.error, ErrCloseBracketWithoutOpen)
	}
	qb.conditionSet = true
}

func (qb *queryBuilder) Sort(sortBy sort.ByWithOrder) {
	qb.sortBy = append(qb.sortBy, sortBy)
}

func (qb *queryBuilder) AddWhere(cmp where.FieldComparator) {
	qb.where = append(qb.where, where.Condition{
		WithNot:      qb.withNot,
		IsOr:         qb.isOr,
		BracketLevel: 1 + qb.bracketLevel,
		Cmp:          cmp,
	})
	qb.withNot = false
	qb.isOr = false
	qb.conditionSet = true
}

func (qb *queryBuilder) OnIteration(cb func(item record.Record)) {
	qb.onIteration = &cb
}

func (qb *queryBuilder) Append(options ...BuilderOption) {
	for _, opt := range options {
		opt.Apply(qb)
	}
}

func (qb *queryBuilder) MakeCopy() Builder {
	cpy := &queryBuilder{
		limitItems:   qb.limitItems,
		startOffset:  qb.startOffset,
		withLimit:    qb.withLimit,
		withNot:      qb.withNot,
		isOr:         qb.isOr,
		conditionSet: qb.conditionSet,
		bracketLevel: qb.bracketLevel,
		where:        make(where.Conditions, len(qb.where)),
		sortBy:       make([]sort.ByWithOrder, len(qb.sortBy)),
		onIteration:  qb.onIteration,
		error:        qb.error,
	}
	copy(cpy.where, qb.where)
	copy(cpy.sortBy, qb.sortBy)
	return cpy
}

func (qb *queryBuilder) Query() Query {
	if qb.bracketLevel > 0 {
		qb.error = multierror.Append(qb.error, ErrInvalidBracketBalance)
	}
	return query{
		offset:              qb.startOffset,
		limit:               qb.limitItems,
		withLimit:           qb.withLimit,
		conditions:          qb.where,
		sorting:             qb.sortBy,
		onIterationCallback: qb.onIteration,
		error:               qb.error.ErrorOrNil(),
	}
}

func NewBuilder(options ...BuilderOption) Builder {
	b := &queryBuilder{} //nolint:exhaustruct
	for _, opt := range options {
		opt.Apply(b)
	}
	return b
}
