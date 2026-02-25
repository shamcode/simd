package main

import (
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/where"
)

type UserQueryBuilder interface {
	query.ChainBuilder[*User, UserQueryBuilder]

	WhereID(condition where.ComparatorType, id ...int64) UserQueryBuilder
	WhereName(condition where.ComparatorType, name ...string) UserQueryBuilder
	WhereStatus(condition where.ComparatorType, status ...Status) UserQueryBuilder
}

type userQueryBuilder struct {
	query.ChainBuilder[*User, UserQueryBuilder]
}

func (uq userQueryBuilder) WhereID(condition where.ComparatorType, value ...int64) UserQueryBuilder {
	uq.AddWhere(query.Where(id, condition, value...))
	return uq
}

func (uq userQueryBuilder) WhereName(condition where.ComparatorType, value ...string) UserQueryBuilder {
	uq.AddWhere(query.Where(name, condition, value...))
	return uq
}

func (uq userQueryBuilder) WhereStatus(condition where.ComparatorType, value ...Status) UserQueryBuilder {
	uq.AddWhere(query.Where(status, condition, value...))
	return uq
}

func NewUserQueryBuilder(
	builder query.BuilderGeneric[*User],
	wrapChain func(qb query.ChainBuilder[*User, UserQueryBuilder]) query.ChainBuilder[*User, UserQueryBuilder],
) UserQueryBuilder {
	var queryBuilder userQueryBuilder

	chain := wrapChain(query.NewCustomChainBuilder[*User, UserQueryBuilder](builder))
	chain.SetOnChain(func() UserQueryBuilder {
		return queryBuilder
	})
	chain.SetOnCopy(func(bcb query.ChainBuilder[*User, UserQueryBuilder]) UserQueryBuilder {
		return userQueryBuilder{
			ChainBuilder: bcb,
		}
	})

	queryBuilder = userQueryBuilder{
		ChainBuilder: chain,
	}

	return queryBuilder
}
