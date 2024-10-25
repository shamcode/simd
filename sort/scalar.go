package sort

import (
	"fmt"

	"github.com/shamcode/simd/record"
)

// Scalar is a special case for sorting by comparing int64 values.
type Scalar[R record.Record] interface {
	Calc(item R) int64
}

type byScalar[R record.Record] struct {
	Scalar[R]
}

func (bi byScalar[R]) Less(a, b R) bool {
	return bi.Calc(a) < bi.Calc(b)
}

func (bi byScalar[R]) String() string {
	return fmt.Sprintf("%#v", bi.Scalar)
}

// ByScalar create sorting by int64 values.
func ByScalar[R record.Record](s Scalar[R]) By[R] {
	return byScalar[R]{s}
}
