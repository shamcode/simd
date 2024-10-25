package sort

import (
	"github.com/shamcode/simd/record"
)

// By is base interface for sorting.
type By[R record.Record] interface {
	Less(a, b R) bool

	// String must return human readable format for debug
	String() string
}

// ByWithOrder wrapper for By interface with order (ASC/DESC)
// ByWithOrder can't be implemented in user level, use Asc() and Desc() for wrap By implementation.
type ByWithOrder[R record.Record] interface {
	By[R]
	order()
}

type asc[R record.Record] struct{ By[R] }

func (asc[R]) order() {}

func (a asc[R]) String() string {
	return a.By.String() + " ASC"
}

// Asc wrap by for sorting in ASC (Ascending) direction.
func Asc[R record.Record](by By[R]) ByWithOrder[R] {
	return asc[R]{by}
}

type desc[R record.Record] struct{ By[R] }

func (desc[R]) order() {}

func (d desc[R]) Less(a, b R) bool {
	return d.By.Less(b, a)
}

func (d desc[R]) String() string {
	return d.By.String() + " DESC"
}

// Desc wrap by for sorting in DESC (Descending) direction.
func Desc[R record.Record](by By[R]) ByWithOrder[R] {
	return desc[R]{by}
}
