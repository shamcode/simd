package sort

import (
	"fmt"
	"github.com/shamcode/simd/record"
)

var (
	_ By = (*byInt64Index)(nil)
)

// Int64IndexCalculation is a special case for sorting by comparing int64 values
type Int64IndexCalculation interface {
	CalcIndex(item record.Record) int64
}

type byInt64Index struct {
	Int64IndexCalculation
}

func (bi *byInt64Index) Equal(a, b record.Record) bool {
	return bi.CalcIndex(a) == bi.CalcIndex(b)
}

func (bi *byInt64Index) Less(a, b record.Record) bool {
	return bi.CalcIndex(a) < bi.CalcIndex(b)
}

func (bi *byInt64Index) String() string {
	return fmt.Sprintf("%#v", bi.Int64IndexCalculation)
}

func ByInt64Index(index Int64IndexCalculation) By {
	return &byInt64Index{
		Int64IndexCalculation: index,
	}
}

func Int64IndexDesc(value int64) int64 {
	return -value
}
