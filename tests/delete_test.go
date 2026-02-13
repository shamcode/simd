package tests

import (
	"context"
	"testing"

	asserts "github.com/shamcode/assert"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/indexes/hash"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/sort"
)

func Test_Delete(t *testing.T) {
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

	// Act
	err := store.Delete(2)

	// Assert
	asserts.Success(t, err)

	var ids []int64

	cur, err := executor.CreateQueryExecutor[*User](store).FetchAll(
		context.Background(),
		query.NewBuilder[*User](
			query.Sort(sort.Asc(userID)),
		).Query(),
	)
	asserts.Success(t, err)

	for cur.Next(context.Background()) {
		ids = append(ids, cur.Item().GetID())
	}

	asserts.Success(t, cur.Err())
	asserts.Equals(t, []int64{1, 3}, ids, "ids from cursor")
}
