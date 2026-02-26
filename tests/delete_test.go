package tests

import (
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

	cur, err := executor.CreateQueryExecutor[*User](store).FetchAll(
		t.Context(),
		query.NewBuilder[*User]().
			Sort(sort.Asc(userID)).
			Query(),
	)
	asserts.Success(t, err)

	var ids []int64 //nolint:prealloc

	for item := range cur.Seq(t.Context()) {
		ids = append(ids, item.GetID())
	}

	asserts.Success(t, cur.Err())
	asserts.Equals(t, []int64{1, 3}, ids, "ids from cursor")
}
