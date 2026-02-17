package tests

import (
	"context"
	"errors"
	"testing"

	asserts "github.com/shamcode/assert"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/indexes/hash"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/query"
)

func Test_Context(t *testing.T) {
	// Arrange
	store := namespace.CreateNamespace[*User]()
	store.AddIndex(hash.NewComparableHashIndex(userID, true))
	asserts.Success(t, store.Insert(&User{ //nolint:exhaustruct
		ID:     1,
		Name:   "First",
		Status: StatusActive,
		Score:  10,
	}))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	_, err := executor.CreateQueryExecutor[*User](store).FetchTotal(ctx, query.NewBuilder[*User]().Query())

	// Assert
	asserts.Equals(t, "context canceled", err.Error(), "check error")
	asserts.Equals(t, true, errors.Is(err, context.Canceled), "error is context.Canceled")
}
