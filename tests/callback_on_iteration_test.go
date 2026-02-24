package tests

import (
	_sort "sort"
	"testing"

	asserts "github.com/shamcode/assert"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/indexes/hash"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
)

func Test_CallbackOnIteration(t *testing.T) {
	// Arrange
	store := namespace.CreateNamespace[*User]()
	store.AddIndex(hash.NewComparableHashIndex(userID, true))
	asserts.Success(t, store.Insert(&User{ //nolint:exhaustruct
		ID:     1,
		Status: StatusActive,
	}))
	asserts.Success(t, store.Insert(&User{ //nolint:exhaustruct
		ID:     2,
		Status: StatusDisabled,
	}))
	asserts.Success(t, store.Insert(&User{ //nolint:exhaustruct
		ID:     3,
		Status: StatusActive,
	}))

	var (
		idsFromCallback []int
		idsFromCursor   []int64
	)

	// Act
	cur, err := executor.CreateQueryExecutor[*User](store).FetchAll(
		t.Context(),
		query.NewChainBuilder(query.NewBuilder[*User]()).
			AddWhere(query.Where(userStatus, where.EQ, StatusActive)).
			Limit(1).
			Sort(sort.Asc(userID)).
			OnIteration(func(item *User) {
				idsFromCallback = append(idsFromCallback, int(item.GetID()))
			}).
			Query(),
	)
	asserts.Success(t, err)

	for cur.Next(t.Context()) {
		idsFromCursor = append(idsFromCursor, cur.Item().GetID())
	}

	// Assert
	_sort.Ints(idsFromCallback)
	asserts.Success(t, cur.Err())
	asserts.Equals(t, []int{1, 3}, idsFromCallback, "ids from callback")
	asserts.Equals(t, []int64{1}, idsFromCursor, "ids from cursor")
}
