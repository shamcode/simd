package query

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
)

type ChainBuilder[R record.Record, B ChainBuilder[R, B]] interface { //nolint:interfacebloat
	Not() B
	Or() B

	OpenBracket() B
	CloseBracket() B

	AddWhere(options BuilderOption) B

	Sort(by sort.ByWithOrder[R]) B

	Limit(limitItems int) B
	Offset(startOffset int) B

	// OnIteration registers a callback to be called for each record before sorting and applying offset/limits
	// but after applying WHERE conditions
	OnIteration(cb func(item R)) B

	MakeCopy() B

	// Error save error to builder
	Error(err error) B

	// Query return build Query
	Query() Query[R]
}

type BaseChainBuilder[R record.Record, Return any] struct {
	builder BuilderGeneric[R]
	onChain func() Return
	onCopy  func(cb *BaseChainBuilder[R, Return]) Return
}

func (cb *BaseChainBuilder[R, Return]) Not() Return {
	cb.builder.Not()
	return cb.onChain()
}

func (cb *BaseChainBuilder[R, Return]) Or() Return {
	cb.builder.Or()
	return cb.onChain()
}

func (cb *BaseChainBuilder[R, Return]) OpenBracket() Return {
	cb.builder.OpenBracket()
	return cb.onChain()
}

func (cb *BaseChainBuilder[R, Return]) CloseBracket() Return {
	cb.builder.CloseBracket()
	return cb.onChain()
}

func (cb *BaseChainBuilder[R, Return]) AddWhere(where BuilderOption) Return {
	cb.builder.Append(where)
	return cb.onChain()
}

func (cb *BaseChainBuilder[R, Return]) Sort(by sort.ByWithOrder[R]) Return {
	cb.builder.Sort(by)
	return cb.onChain()
}

func (cb *BaseChainBuilder[R, Return]) Limit(limitItems int) Return {
	cb.builder.Limit(limitItems)
	return cb.onChain()
}

func (cb *BaseChainBuilder[R, Return]) Offset(startOffset int) Return {
	cb.builder.Offset(startOffset)
	return cb.onChain()
}

func (cb *BaseChainBuilder[R, Return]) OnIteration(fn func(item R)) Return {
	cb.builder.OnIteration(fn)
	return cb.onChain()
}

func (cb *BaseChainBuilder[R, Return]) MakeCopy() Return {
	return cb.onCopy(&BaseChainBuilder[R, Return]{
		builder: cb.builder.MakeCopy(),
		onChain: cb.onChain,
		onCopy:  cb.onCopy,
	})
}

func (cb *BaseChainBuilder[R, Return]) Error(err error) Return {
	cb.builder.Error(err)
	return cb.onChain()
}

func (cb *BaseChainBuilder[R, Return]) Query() Query[R] {
	return cb.builder.Query()
}

func NewCustomChainBuilder[
	R record.Record,
	Return any,
](
	b BuilderGeneric[R],
	onChain func() Return,
	onCopy func(cb *BaseChainBuilder[R, Return]) Return,
) *BaseChainBuilder[R, Return] {
	return &BaseChainBuilder[R, Return]{
		builder: b,
		onCopy:  onCopy,
		onChain: onChain,
	}
}

type defaultChainBuilder[R record.Record] struct {
	*BaseChainBuilder[R, *defaultChainBuilder[R]]
}

func NewChainBuilder[R record.Record](builder BuilderGeneric[R]) ChainBuilder[R, *defaultChainBuilder[R]] {
	var rcb *defaultChainBuilder[R]

	rcb = &defaultChainBuilder[R]{
		BaseChainBuilder: NewCustomChainBuilder[
			R,
			*defaultChainBuilder[R],
		](
			builder,
			func() *defaultChainBuilder[R] { return rcb },
			func(cb *BaseChainBuilder[R, *defaultChainBuilder[R]]) *defaultChainBuilder[R] {
				return &defaultChainBuilder[R]{
					BaseChainBuilder: cb,
				}
			},
		),
	}

	return rcb
}
