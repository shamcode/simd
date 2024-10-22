package sort

import (
	"fmt"

	"github.com/shamcode/simd/record"
)

// Scalar is a special case for sorting by comparing int64 values.
type Scalar interface {
	Calc(item record.Record) int64
}

type byScalar struct {
	Scalar
}

func (bi byScalar) Less(a, b record.Record) bool {
	return bi.Calc(a) < bi.Calc(b)
}

func (bi byScalar) String() string {
	return fmt.Sprintf("%#v", bi.Scalar)
}

// ByScalar create sorting by int64 values.
func ByScalar(s Scalar) By {
	return byScalar{s}
}
