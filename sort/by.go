package sort

import (
	"github.com/shamcode/simd/record"
)

// By is base interface for sorting.
type By interface {
	Less(a, b record.Record) bool

	// String must return human readable format for debug
	String() string
}

// ByWithOrder wrapper for By interface with order (ASC/DESC)
// ByWithOrder can't be implemented in user level, use Asc() and Desc() for wrap By implementation.
type ByWithOrder interface {
	By
	order()
}

type asc struct{ By }

func (_ asc) order() {}

func (a asc) String() string {
	return a.By.String() + " ASC"
}

// Asc wrap by for sorting in ASC (Ascending) direction.
func Asc(by By) ByWithOrder {
	return asc{by}
}

type desc struct{ By }

func (_ desc) order() {}

func (d desc) Less(a, b record.Record) bool {
	return d.By.Less(b, a)
}

func (d desc) String() string {
	return d.By.String() + " DESC"
}

// Desc wrap by for sorting in DESC (Descending) direction.
func Desc(by By) ByWithOrder {
	return desc{by}
}
