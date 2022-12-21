package sort

import "github.com/shamcode/simd/record"

type By interface {
	Equal(a, b record.Record) bool
	Less(a, b record.Record) bool
	String() string
}
