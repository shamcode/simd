package namespace

import (
	"errors"
	"fmt"
)

var ErrExecuteQuery = errors.New("execute query")

func wrapErrors(simdErr, err error) error {
	return fmt.Errorf("%w: %s", simdErr, err)
}
