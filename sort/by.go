package sort

import "github.com/shamcode/simd/record"

type By interface {
	Less(a, b record.Record) bool
	String() string
}
