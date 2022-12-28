package namespace

import (
	"errors"
	"fmt"
)

var (
	ErrValidateQuery = errors.New("validate query")
	ErrExecuteQuery  = errors.New("execute query")
)

func wrapErrors(simdErr, err error) error {
	return fmt.Errorf("%w: %s", simdErr, err)
}
