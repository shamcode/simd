package debug

import (
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
)

type combine[R record.Record, Return query.Builder[R, Return]] struct {
	base  query.Builder[R, Return]
	debug *debugQueryBuilder[R]
}

func (cb *combine[R, Return]) Not() Return {
	cb.debug.not()
	return cb.base.Not()
}

func (cb *combine[R, Return]) Or() Return {
	cb.debug.or()
	return cb.base.Or()
}

func (cb *combine[R, Return]) OpenBracket() Return {
	cb.debug.openBracket()
	return cb.base.OpenBracket()
}

func (cb *combine[R, Return]) CloseBracket() Return {
	cb.debug.closeBracket()
	return cb.base.CloseBracket()
}

func (cb *combine[R, Return]) AddWhere(where query.AddWhereOption[R]) Return {
	cb.debug.addWhere(where.Cmp)
	return cb.base.AddWhere(where)
}

func (cb *combine[R, Return]) Sort(by sort.ByWithOrder[R]) Return {
	cb.debug.sort(by)
	return cb.base.Sort(by)
}

func (cb *combine[R, Return]) Limit(limitItems int) Return {
	cb.debug.limit(limitItems)
	return cb.base.Limit(limitItems)
}

func (cb *combine[R, Return]) Offset(startOffset int) Return {
	cb.debug.offset(startOffset)
	return cb.base.Offset(startOffset)
}

func (cb *combine[R, Return]) OnIteration(fn func(item R)) Return {
	return cb.base.OnIteration(fn)
}

func (cb *combine[R, Return]) OnCopy(cb2 query.Builder[R, Return]) Return {
	return cb.base.OnCopy(cb2)
}

func (cb *combine[R, Return]) MakeCopy() Return {
	return cb.OnCopy(&combine[R, Return]{
		base:  cb.base.MakeCopy(),
		debug: cb.debug.makeCopy(),
	})
}

func (cb *combine[R, Return]) SetOnChain(onChain Return) {
	cb.base.SetOnChain(onChain)
}

func (cb *combine[R, Return]) SetOnCopy(onCopy func(onCopy query.Builder[R, Return]) Return) {
	cb.base.SetOnCopy(onCopy)
}

func (cb *combine[R, Return]) Error(err error) Return {
	return cb.base.Error(err)
}

func (cb *combine[R, Return]) Query() query.Query[R] {
	return NewQueryWithDumper[R](
		cb.base.Query(),
		cb.debug.Dump(),
	)
}

func (cb *combine[R, Return]) setFieldComparatorDumper(dumber FieldComparatorDumper[R]) {
	cb.debug.setFieldComparatorDumper(dumber)
}

func WrapBuilder[R record.Record, Return query.Builder[R, Return]](
	queryBuilder query.Builder[R, Return],
) query.Builder[R, Return] {
	wrapped := &combine[R, Return]{
		debug: newDebugQueryBuilder[R](),
		base:  queryBuilder,
	}

	if _, ok := queryBuilder.(query.DefaultBuilder[R]); ok {
		// Replace default combine builder with wrapped.
		cb1 := query.NewExtendedBuilder[R, query.DefaultBuilder[R]]()
		cb1.SetOnChain(any(wrapped).(query.DefaultBuilder[R]))
		cb1.SetOnCopy(func(onCopy query.Builder[R, query.DefaultBuilder[R]]) query.DefaultBuilder[R] {
			return onCopy.(query.DefaultBuilder[R])
		})

		wrapped.base = any(cb1).(query.Builder[R, Return])
	}

	return wrapped
}

func WrapBuilderWithDumper[R record.Record, Return query.Builder[R, Return]](
	cb query.Builder[R, Return],
	dumper FieldComparatorDumper[R],
) query.Builder[R, Return] {
	b := WrapBuilder[R, Return](cb).(*combine[R, Return])
	b.setFieldComparatorDumper(dumper)

	return b
}
