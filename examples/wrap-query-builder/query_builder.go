package main

import (
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/where"
)

type UserQueryBuilder interface {
	query.Builder[*User, UserQueryBuilder]

	WhereID(condition where.ComparatorType, id ...int64) UserQueryBuilder
	WhereName(condition where.ComparatorType, name ...string) UserQueryBuilder
	WhereStatus(condition where.ComparatorType, status ...Status) UserQueryBuilder
}

type userQueryBuilder struct {
	query.Builder[*User, UserQueryBuilder]
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
	baseBuilder query.Builder[*User, UserQueryBuilder],
) UserQueryBuilder {
	queryBuilder := &userQueryBuilder{
		Builder: baseBuilder,
	}

	baseBuilder.SetOnChain(queryBuilder)
	baseBuilder.SetOnCopy(func(bcb query.Builder[*User, UserQueryBuilder]) UserQueryBuilder {
		return userQueryBuilder{
			Builder: bcb,
		}
	})

	return queryBuilder
}
