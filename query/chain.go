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

	AddWhere(options AddWhereOption[R]) B

	Sort(by sort.ByWithOrder[R]) B

	Limit(limitItems int) B
	Offset(startOffset int) B

	// OnIteration registers a callback to be called for each record before sorting and applying offset/limits
	// but after applying WHERE conditions
	OnIteration(cb func(item R)) B

	MakeCopy() B

	// OnCopy call registered callback for make copy of builder
	OnCopy(cb ChainBuilder[R, B]) B

	SetOnChain(onChain func() B)
	SetOnCopy(onCopy func(cb ChainBuilder[R, B]) B)

	// Error save error to builder
	Error(err error) B

	// Query return build Query
	Query() Query[R]
}

type BaseChainBuilder[R record.Record, Return ChainBuilder[R, Return]] struct {
	builder BuilderGeneric[R]
	onChain func() Return
	onCopy  func(cb ChainBuilder[R, Return]) Return
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

func (cb *BaseChainBuilder[R, Return]) AddWhere(where AddWhereOption[R]) Return {
	where.Apply(cb.builder)
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

func (cb *BaseChainBuilder[R, Return]) OnCopy(cb2 ChainBuilder[R, Return]) Return {
	return cb.onCopy(cb2)
}

func (cb *BaseChainBuilder[R, Return]) Error(err error) Return {
	cb.builder.Error(err)
	return cb.onChain()
}

func (cb *BaseChainBuilder[R, Return]) Query() Query[R] {
	return cb.builder.Query()
}

func (cb *BaseChainBuilder[R, Return]) SetOnChain(onChain func() Return) {
	cb.onChain = onChain
}

func (cb *BaseChainBuilder[R, Return]) SetOnCopy(onCopy func(cb ChainBuilder[R, Return]) Return) {
	cb.onCopy = onCopy
}

// Builder return query.Builder.
// TODO: remove.
func (cb *BaseChainBuilder[R, Return]) Builder() BuilderGeneric[R] { return cb.builder }

func NewCustomChainBuilder[
	R record.Record,
	Return ChainBuilder[R, Return],
](
	b BuilderGeneric[R],
) ChainBuilder[R, Return] {
	return &BaseChainBuilder[R, Return]{
		builder: b,
		onChain: nil,
		onCopy:  nil,
	}
}

type DefaultChainBuilder[R record.Record] interface {
	ChainBuilder[R, DefaultChainBuilder[R]]

	// TODO: remove
	Builder() BuilderGeneric[R]
}

type defaultChainBuilder[R record.Record] struct {
	*BaseChainBuilder[R, DefaultChainBuilder[R]]
}

func NewChainBuilder[R record.Record](builder BuilderGeneric[R]) DefaultChainBuilder[R] {
	rcb := &defaultChainBuilder[R]{
		BaseChainBuilder: nil,
	}

	chain := NewCustomChainBuilder[R, DefaultChainBuilder[R]](builder)
	chain.SetOnChain(func() DefaultChainBuilder[R] {
		return rcb
	})
	chain.SetOnCopy(func(cb ChainBuilder[R, DefaultChainBuilder[R]]) DefaultChainBuilder[R] {
		return &defaultChainBuilder[R]{
			BaseChainBuilder: cb.(*BaseChainBuilder[R, DefaultChainBuilder[R]]),
		}
	})

	rcb.BaseChainBuilder = chain.(*BaseChainBuilder[R, DefaultChainBuilder[R]])

	return rcb
}
