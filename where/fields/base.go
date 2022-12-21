package fields

import (
	"errors"
	"github.com/shamcode/simd/where"
)

var (
	errNotImplementComparator = errors.New("not implemented comparator")
)

type BaseFieldComparator struct {
	Cmp where.ComparatorType
}

func (fc BaseFieldComparator) GetType() where.ComparatorType {
	return fc.Cmp
}
