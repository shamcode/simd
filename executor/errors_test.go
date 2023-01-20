package executor

import (
	"errors"
	"github.com/shamcode/simd/asserts"
	"io/fs"
	"testing"
)

func TestErrors(t *testing.T) {
	testCases := []struct {
		Error          error
		IsError        error
		ExpectedString string
	}{
		{
			Error:          NewErrValidateQuery(fs.ErrPermission),
			IsError:        ErrValidateQuery{},
			ExpectedString: "validate query: permission denied",
		},
		{
			Error:          NewErrExecuteQuery(fs.ErrPermission),
			IsError:        ErrExecuteQuery{},
			ExpectedString: "execute query: permission denied",
		},
	}

	for _, err := range testCases {
		t.Run(err.Error.Error(), func(t *testing.T) {
			asserts.Equals(t, true, errors.Is(err.Error, err.IsError), "is error")
			asserts.Equals(t, true, errors.Is(err.Error, fs.ErrPermission), "is wrapped error")
			asserts.Equals(t, err.ExpectedString, err.Error.Error(), "string")
		})
	}
}
