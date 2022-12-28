package comparators

import (
	"errors"
)

var (
	ErrNotImplementComparator = errors.New("not implemented ComparatorType")
	ErrFailCastType           = errors.New("can't cast type")
)
