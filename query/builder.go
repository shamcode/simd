package query

import (
	"errors"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
)

type Builder interface {
	Limit(limitItems int)
	Offset(startOffset int)
	Not()
	Or()
	OpenBracket()
	CloseBracket()

	// Error save error to builder
	Error(err error)
}

type BuilderGeneric[R record.Record] interface {
	Builder

	AddWhere(cmp where.FieldComparator[R])
	Sort(by sort.ByWithOrder[R])

	// OnIteration registers a callback to be called for each record before sorting and applying offset/limits
	// but after applying WHERE conditions
	OnIteration(cb func(item R))

	MakeCopy() BuilderGeneric[R]

	// Query return build Query
	Query() Query[R]
}

type queryBuilder[R record.Record] struct {
	limitItems   int
	startOffset  int
	withLimit    bool
	withNot      bool
	isOr         bool
	conditionSet bool
	bracketLevel int
	where        where.Conditions[R]
	sortBy       []sort.ByWithOrder[R]
	onIteration  *func(item R)
	errors       []error
}

func (qb *queryBuilder[R]) Limit(limitItems int) {
	qb.limitItems = limitItems
	qb.withLimit = true
}

func (qb *queryBuilder[R]) Offset(startOffset int) {
	qb.startOffset = startOffset
}

func (qb *queryBuilder[R]) Not() {
	qb.withNot = !qb.withNot
}

func (qb *queryBuilder[R]) Or() {
	qb.isOr = true
	if !qb.conditionSet {
		qb.errors = append(qb.errors, ErrOrBeforeAnyConditions)
	}
}

func (qb *queryBuilder[R]) OpenBracket() {
	if qb.withNot {
		qb.errors = append(qb.errors, ErrNotOpenBracket)
	}

	qb.conditionSet = false
	qb.bracketLevel += 1
}

func (qb *queryBuilder[R]) CloseBracket() {
	qb.bracketLevel -= 1
	if qb.bracketLevel == -1 {
		qb.errors = append(qb.errors, ErrCloseBracketWithoutOpen)
	}

	qb.conditionSet = true
}

func (qb *queryBuilder[R]) Error(err error) {
	if err != nil {
		qb.errors = append(qb.errors, err)
	}
}

func (qb *queryBuilder[R]) Sort(sortBy sort.ByWithOrder[R]) {
	qb.sortBy = append(qb.sortBy, sortBy)
}

func (qb *queryBuilder[R]) AddWhere(cmp where.FieldComparator[R]) {
	qb.where = append(qb.where, where.Condition[R]{
		WithNot:      qb.withNot,
		IsOr:         qb.isOr,
		BracketLevel: 1 + qb.bracketLevel,
		Cmp:          cmp,
	})
	qb.withNot = false
	qb.isOr = false
	qb.conditionSet = true
}

func (qb *queryBuilder[R]) OnIteration(cb func(item R)) {
	qb.onIteration = &cb
}

func (qb *queryBuilder[R]) MakeCopy() BuilderGeneric[R] {
	cpy := &queryBuilder[R]{
		limitItems:   qb.limitItems,
		startOffset:  qb.startOffset,
		withLimit:    qb.withLimit,
		withNot:      qb.withNot,
		isOr:         qb.isOr,
		conditionSet: qb.conditionSet,
		bracketLevel: qb.bracketLevel,
		where:        make(where.Conditions[R], len(qb.where)),
		sortBy:       make([]sort.ByWithOrder[R], len(qb.sortBy)),
		onIteration:  qb.onIteration,
		errors:       make([]error, len(qb.errors)),
	}
	copy(cpy.where, qb.where)
	copy(cpy.sortBy, qb.sortBy)
	copy(cpy.errors, qb.errors)

	return cpy
}

func (qb *queryBuilder[R]) Query() Query[R] {
	if qb.bracketLevel > 0 {
		qb.errors = append(qb.errors, ErrInvalidBracketBalance)
	}

	return query[R]{
		offset:              qb.startOffset,
		limit:               qb.limitItems,
		withLimit:           qb.withLimit,
		conditions:          qb.where,
		sorting:             qb.sortBy,
		onIterationCallback: qb.onIteration,
		error:               errors.Join(qb.errors...),
	}
}

func NewBuilder[R record.Record]() BuilderGeneric[R] {
	return &queryBuilder[R]{} //nolint:exhaustruct
}
