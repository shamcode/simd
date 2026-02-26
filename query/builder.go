package query

import (
	"errors"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
)

type Builder[R record.Record, B Builder[R, B]] interface { //nolint:interfacebloat
	Not() B
	Or() B

	OpenBracket() B
	CloseBracket() B

	AddWhere(options AddWhereOption[R]) B

	Sort(by sort.ByWithOrder[R]) B

	Limit(limitItems int) B
	Offset(startOffset int) B

	// OnIteration registers a callback to be called for each record before sorting and applying offset/limits
	// but after applying WHERE conditions
	OnIteration(cb func(item R)) B

	MakeCopy() B

	// OnCopy call registered callback for make copy of builder
	OnCopy(cb Builder[R, B]) B

	SetOnChain(ret B)
	SetOnCopy(onCopy func(cb Builder[R, B]) B)

	// Error save error to builder
	Error(err error) B

	// Query return build Query
	Query() Query[R]
}

type BaseBuilder[R record.Record, Return Builder[R, Return]] struct {
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
	onChain      Return
	onCopy       func(cb Builder[R, Return]) Return
}

func (qb *BaseBuilder[R, Return]) Not() Return {
	qb.withNot = !qb.withNot

	return qb.onChain
}

func (qb *BaseBuilder[R, Return]) Or() Return {
	qb.isOr = true
	if !qb.conditionSet {
		qb.errors = append(qb.errors, ErrOrBeforeAnyConditions)
	}

	return qb.onChain
}

func (qb *BaseBuilder[R, Return]) OpenBracket() Return {
	if qb.withNot {
		qb.errors = append(qb.errors, ErrNotOpenBracket)
	}

	qb.conditionSet = false
	qb.bracketLevel += 1

	return qb.onChain
}

func (qb *BaseBuilder[R, Return]) CloseBracket() Return {
	qb.bracketLevel -= 1
	if qb.bracketLevel == -1 {
		qb.errors = append(qb.errors, ErrCloseBracketWithoutOpen)
	}

	qb.conditionSet = true

	return qb.onChain
}

func (qb *BaseBuilder[R, Return]) AddWhere(cmp AddWhereOption[R]) Return {
	if cmp.Error != nil {
		qb.errors = append(qb.errors, cmp.Error)

		return qb.onChain
	}

	qb.where = append(qb.where, where.Condition[R]{
		WithNot:      qb.withNot,
		IsOr:         qb.isOr,
		BracketLevel: 1 + qb.bracketLevel,
		Cmp:          cmp.Cmp,
	})
	qb.withNot = false
	qb.isOr = false
	qb.conditionSet = true

	return qb.onChain
}

func (qb *BaseBuilder[R, Return]) Sort(by sort.ByWithOrder[R]) Return {
	qb.sortBy = append(qb.sortBy, by)

	return qb.onChain
}

func (qb *BaseBuilder[R, Return]) Limit(limitItems int) Return {
	qb.limitItems = limitItems
	qb.withLimit = true

	return qb.onChain
}

func (qb *BaseBuilder[R, Return]) Offset(startOffset int) Return {
	qb.startOffset = startOffset

	return qb.onChain
}

func (qb *BaseBuilder[R, Return]) OnIteration(fn func(item R)) Return {
	qb.onIteration = &fn

	return qb.onChain
}

func (qb *BaseBuilder[R, Return]) MakeCopy() Return {
	cpy := &BaseBuilder[R, Return]{
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
		onChain:      qb.onChain,
		onCopy:       qb.onCopy,
	}
	copy(cpy.where, qb.where)
	copy(cpy.sortBy, qb.sortBy)
	copy(cpy.errors, qb.errors)

	return qb.onCopy(cpy)
}

func (qb *BaseBuilder[R, Return]) OnCopy(cb2 Builder[R, Return]) Return {
	return qb.onCopy(cb2)
}

func (qb *BaseBuilder[R, Return]) Error(err error) Return {
	if err != nil {
		qb.errors = append(qb.errors, err)
	}

	return qb.onChain
}

func (qb *BaseBuilder[R, Return]) Query() Query[R] {
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

func (qb *BaseBuilder[R, Return]) SetOnChain(onChain Return) {
	qb.onChain = onChain
}

func (qb *BaseBuilder[R, Return]) SetOnCopy(onCopy func(cb Builder[R, Return]) Return) {
	qb.onCopy = onCopy
}

func NewExtendedBuilder[
	R record.Record,
	Return Builder[R, Return],
]() Builder[R, Return] {
	return &BaseBuilder[R, Return]{} //nolint:exhaustruct
}

type DefaultBuilder[R record.Record] interface {
	Builder[R, DefaultBuilder[R]]
}

func NewBuilder[R record.Record]() DefaultBuilder[R] {
	chain := NewExtendedBuilder[R, DefaultBuilder[R]]()

	chain.SetOnChain(chain)
	chain.SetOnCopy(func(cb Builder[R, DefaultBuilder[R]]) DefaultBuilder[R] {
		return cb
	})

	return chain
}
