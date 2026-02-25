package debug

import (
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
)

type chain[R record.Record, Return query.ChainBuilder[R, Return]] struct {
	chain query.ChainBuilder[R, Return]
	debug *debugQueryBuilder[R]
}

func (cb *chain[R, Return]) Not() Return {
	cb.debug.not()
	return cb.chain.Not()
}

func (cb *chain[R, Return]) Or() Return {
	cb.debug.or()
	return cb.chain.Or()
}

func (cb *chain[R, Return]) OpenBracket() Return {
	cb.debug.openBracket()
	return cb.chain.OpenBracket()
}

func (cb *chain[R, Return]) CloseBracket() Return {
	cb.debug.closeBracket()
	return cb.chain.CloseBracket()
}

func (cb *chain[R, Return]) AddWhere(where query.AddWhereOption[R]) Return {
	cb.debug.addWhere(where.Cmp)
	return cb.chain.AddWhere(where)
}

func (cb *chain[R, Return]) Sort(by sort.ByWithOrder[R]) Return {
	cb.debug.sort(by)
	return cb.chain.Sort(by)
}

func (cb *chain[R, Return]) Limit(limitItems int) Return {
	cb.debug.limit(limitItems)
	return cb.chain.Limit(limitItems)
}

func (cb *chain[R, Return]) Offset(startOffset int) Return {
	cb.debug.offset(startOffset)
	return cb.chain.Offset(startOffset)
}

func (cb *chain[R, Return]) OnIteration(fn func(item R)) Return {
	return cb.chain.OnIteration(fn)
}

func (cb *chain[R, Return]) OnCopy(cb2 query.ChainBuilder[R, Return]) Return {
	return cb.chain.OnCopy(cb2)
}

func (cb *chain[R, Return]) MakeCopy() Return {
	return cb.OnCopy(&chain[R, Return]{
		chain: cb.chain.MakeCopy(),
		debug: cb.debug.makeCopy(),
	})
}

func (cb *chain[R, Return]) SetOnChain(onChain func() Return) {
	cb.chain.SetOnChain(onChain)
}

func (cb *chain[R, Return]) SetOnCopy(onCopy func(onCopy query.ChainBuilder[R, Return]) Return) {
	cb.chain.SetOnCopy(onCopy)
}

func (cb *chain[R, Return]) Error(err error) Return {
	return cb.chain.Error(err)
}

func (cb *chain[R, Return]) Query() query.Query[R] {
	return NewQueryWithDumper[R](
		cb.chain.Query(),
		cb.debug.Dump(),
	)
}

// TODO: remove.
func (cb *chain[R, Return]) Builder() query.BuilderGeneric[R] {
	return nil
}

func (cb *chain[R, Return]) setFieldComparatorDumper(dumber FieldComparatorDumper[R]) {
	cb.debug.setFieldComparatorDumper(dumber)
}

func WrapChainBuilder[R record.Record, Return query.ChainBuilder[R, Return]](
	queryBuilder query.ChainBuilder[R, Return],
) query.ChainBuilder[R, Return] {
	wrapped := &chain[R, Return]{
		debug: newDebugQueryBuilder[R](),
		chain: queryBuilder,
	}

	def, ok := queryBuilder.(query.DefaultChainBuilder[R])
	if ok {
		// Replace default chain builder with wrapped.
		cb1 := query.NewCustomChainBuilder[R, query.DefaultChainBuilder[R]](def.Builder())
		cb1.SetOnChain(func() query.DefaultChainBuilder[R] {
			return any(wrapped).(query.DefaultChainBuilder[R])
		})
		cb1.SetOnCopy(func(onCopy query.ChainBuilder[R, query.DefaultChainBuilder[R]]) query.DefaultChainBuilder[R] {
			return onCopy.(query.DefaultChainBuilder[R])
		})

		wrapped.chain = any(cb1).(query.ChainBuilder[R, Return])
	}

	return wrapped
}

func WrapChainBuilderWithDumper[R record.Record, Return query.ChainBuilder[R, Return]](
	cb query.ChainBuilder[R, Return],
	dumper FieldComparatorDumper[R],
) query.ChainBuilder[R, Return] {
	b := WrapChainBuilder[R, Return](cb).(*chain[R, Return])
	b.setFieldComparatorDumper(dumper)

	return b
}
