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
	*query.BaseChainBuilder[*User, UserQueryBuilder]
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

func NewUserQueryBuilder(b query.BuilderGeneric[*User]) UserQueryBuilder {
	var queryBuilder UserQueryBuilder

	queryBuilder = userQueryBuilder{
		BaseChainBuilder: query.NewCustomChainBuilder(
			b,
			func() UserQueryBuilder { return queryBuilder },
			func(bcb *query.BaseChainBuilder[*User, UserQueryBuilder]) UserQueryBuilder {
				return userQueryBuilder{
					BaseChainBuilder: bcb,
				}
			},
		),
	}

	return queryBuilder
}
