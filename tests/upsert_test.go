package tests

import (
	"testing"

	asserts "github.com/shamcode/assert"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/indexes/hash"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/where"
)

func Test_Upsert(t *testing.T) {
	// Arrange
	store := namespace.CreateNamespace[*User]()
	store.AddIndex(hash.NewComparableHashIndex(userID, true))
	store.AddIndex(hash.NewComparableHashIndex(userStatus, true))
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
	err := store.Upsert(&User{ //nolint:exhaustruct
		ID:     2,
		Status: StatusActive,
	})

	// Assert
	asserts.Success(t, err)

	q := query.NewChainBuilder(query.NewBuilder[*User]()).
		AddWhere(query.Where(userID, where.EQ, 2)).
		Query()

	cur, err := executor.CreateQueryExecutor[*User](store).FetchAll(t.Context(), q)

	asserts.Success(t, err)
	asserts.Success(t, cur.Err())
	asserts.Equals(t, StatusActive, cur.Item().Status, "status")
}
