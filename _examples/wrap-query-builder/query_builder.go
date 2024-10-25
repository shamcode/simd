package main

import (
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

type UserQueryBuilder interface { //nolint:interfacebloat
	Limit(limitItems int) UserQueryBuilder
	Offset(startOffset int) UserQueryBuilder
	Not() UserQueryBuilder
	Or() UserQueryBuilder
	OpenBracket() UserQueryBuilder
	CloseBracket() UserQueryBuilder

	WhereID(condition where.ComparatorType, id ...int64) UserQueryBuilder
	WhereName(condition where.ComparatorType, name ...string) UserQueryBuilder
	WhereStatus(condition where.ComparatorType, status ...Status) UserQueryBuilder

	Sort(by sort.ByWithOrder[*User]) UserQueryBuilder
	MakeCopy() UserQueryBuilder
	Query() query.Query[*User]
}

type userQueryBuilder struct {
	builder query.BuilderGeneric[*User]
}

func (uq userQueryBuilder) Limit(limitItems int) UserQueryBuilder {
	uq.builder.Limit(limitItems)
	return uq
}

func (uq userQueryBuilder) Offset(startOffset int) UserQueryBuilder {
	uq.builder.Offset(startOffset)
	return uq
}

func (uq userQueryBuilder) Not() UserQueryBuilder {
	uq.builder.Not()
	return uq
}

func (uq userQueryBuilder) Or() UserQueryBuilder {
	uq.builder.Or()
	return uq
}

func (uq userQueryBuilder) OpenBracket() UserQueryBuilder {
	uq.builder.OpenBracket()
	return uq
}

func (uq userQueryBuilder) CloseBracket() UserQueryBuilder {
	uq.builder.CloseBracket()
	return uq
}

func (uq userQueryBuilder) WhereID(condition where.ComparatorType, value ...int64) UserQueryBuilder {
	uq.builder.AddWhere(comparators.ComparableFieldComparator[*User, int64]{
		Cmp:    condition,
		Getter: id,
		Value:  value,
	})
	return uq
}

func (uq userQueryBuilder) WhereName(condition where.ComparatorType, value ...string) UserQueryBuilder {
	uq.builder.AddWhere(comparators.StringFieldComparator[*User]{
		Cmp:    condition,
		Getter: name,
		Value:  value,
	})
	return uq
}

func (uq userQueryBuilder) WhereStatus(condition where.ComparatorType, value ...Status) UserQueryBuilder {
	enums := make([]record.Enum[uint8], len(value))
	for i, x := range value {
		enums[i] = x
	}
	uq.builder.AddWhere(comparators.EnumFieldComparator[*User, uint8]{
		Cmp:    condition,
		Getter: status,
		Value:  enums,
	})
	return uq
}

func (uq userQueryBuilder) Sort(by sort.ByWithOrder[*User]) UserQueryBuilder {
	uq.builder.Sort(by)
	return uq
}

func (uq userQueryBuilder) MakeCopy() UserQueryBuilder {
	return &userQueryBuilder{
		builder: uq.builder.MakeCopy(),
	}
}

func (uq userQueryBuilder) Query() query.Query[*User] {
	return uq.builder.Query()
}

func NewUserQueryBuilder(b query.BuilderGeneric[*User]) UserQueryBuilder {
	return userQueryBuilder{
		builder: b,
	}
}
