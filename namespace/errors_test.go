package namespace

import (
	"errors"
	"testing"

	asserts "github.com/shamcode/assert"
)

func TestErrors(t *testing.T) {
	testCases := []struct {
		Error          error
		IsError        error
		ExpectedString string
	}{
		{
			Error:          NewRecordAlreadyExists(10),
			IsError:        RecordAlreadyExistsError{},
			ExpectedString: "simd: record with passed id already exists: ID == 10",
		},
	}

	for _, err := range testCases {
		t.Run(err.Error.Error(), func(t *testing.T) {
			asserts.Equals(t, true, errors.Is(err.Error, err.IsError), "is error")
			asserts.Equals(t, err.ExpectedString, err.Error.Error(), "string")
		})
	}
}
