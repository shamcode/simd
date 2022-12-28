package sort

import (
	"fmt"
	"github.com/shamcode/simd/record"
)

var _ By = (*byString)(nil)

type byString struct {
	getter *record.StringGetter
	asc    bool
}

func (bs *byString) Equal(a, b record.Record) bool {
	return bs.getter.Get(a) == bs.getter.Get(b)
}

func (bs *byString) Less(a, b record.Record) bool {
	if bs.asc {
		return bs.getter.Get(a) < bs.getter.Get(b)
	} else {
		return bs.getter.Get(a) > bs.getter.Get(b)
	}
}

func (bs *byString) String() string {
	return fmt.Sprintf("%#v", bs)
}

func ByStringAsc(getter *record.StringGetter) By {
	return &byString{
		getter: getter,
		asc:    true,
	}
}

func ByStringDesc(getter *record.StringGetter) By {
	return &byString{
		getter: getter,
		asc:    false,
	}
}
