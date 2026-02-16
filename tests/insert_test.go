package tests

import (
	"testing"

	asserts "github.com/shamcode/assert"
	"github.com/shamcode/simd/indexes/hash"
	"github.com/shamcode/simd/namespace"
)

func Test_InsertAlreadyExisted(t *testing.T) {
	// Arrange
	store := namespace.CreateNamespace[*User]()
	store.AddIndex(hash.NewComparableHashIndex(userID, true))
	asserts.Success(t, store.Insert(&User{ //nolint:exhaustruct
		ID:     1,
		Status: StatusActive,
	}))

	// Act
	err := store.Insert(&User{ //nolint:exhaustruct
		ID:     1,
		Status: StatusDisabled,
	})

	// Assert
	asserts.Equals(t, "simd: record with passed id already exists: ID == 1", err.Error(), "check error")
}
