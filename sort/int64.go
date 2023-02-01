package sort

import (
	"fmt"
	"github.com/shamcode/simd/record"
)

var _ By = (*byInt64Index)(nil)

// Int64IndexCalculation is a special case for sorting by comparing int64 values
type Int64IndexCalculation interface {
	CalcIndex(item record.Record) int64
}

type byInt64Index struct {
	Int64IndexCalculation
	asc bool
}

func (bi *byInt64Index) Less(a, b record.Record) bool {
	if bi.asc {
		return bi.CalcIndex(a) < bi.CalcIndex(b)
	} else {
		return bi.CalcIndex(a) > bi.CalcIndex(b)
	}
}

func (bi *byInt64Index) String() string {
	var direction string
	if bi.asc {
		direction = "ASC"
	} else {
		direction = "DESC"
	}
	return fmt.Sprintf("%#v %s", bi.Int64IndexCalculation, direction)
}

// ByInt64IndexAsc create sorting by int64 index in ASC direction
func ByInt64IndexAsc(index Int64IndexCalculation) By {
	return &byInt64Index{
		Int64IndexCalculation: index,
		asc:                   true,
	}
}

// ByInt64IndexAsc create sorting by int64 index in DESC direction
func ByInt64IndexDesc(index Int64IndexCalculation) By {
	return &byInt64Index{
		Int64IndexCalculation: index,
		asc:                   false,
	}
}
