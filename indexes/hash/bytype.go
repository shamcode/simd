package hash

import (
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/indexes/compute"
	"github.com/shamcode/simd/record"
)

func NewBoolHashIndex(getter *record.BoolGetter, unique bool) indexes.Index {
	return NewIndex(
		getter.Field,
		compute.CreateBoolIndexComputation(getter),
		CreateHashTable(),
		unique,
	)
}

func NewEnum8HashIndex(getter *record.Enum8Getter, unique bool) indexes.Index {
	return NewIndex(
		getter.Field,
		compute.CreateEnum8IndexComputation(getter),
		CreateHashTable(),
		unique,
	)
}

func NewEnum16HashIndex(getter *record.Enum16Getter, unique bool) indexes.Index {
	return NewIndex(
		getter.Field,
		compute.CreateEnum16IndexComputation(getter),
		CreateHashTable(),
		unique,
	)
}

func NewInt32HashIndex(getter *record.Int32Getter, unique bool) indexes.Index {
	return NewIndex(
		getter.Field,
		compute.CreateInt32IndexComputation(getter),
		CreateHashTable(),
		unique,
	)
}

func NewInt64HashIndex(getter *record.Int64Getter, unique bool) indexes.Index {
	return NewIndex(
		getter.Field,
		compute.CreateInt64IndexComputation(getter),
		CreateHashTable(),
		unique,
	)
}

func NewStringHashIndex(getter *record.StringGetter, unique bool) indexes.Index {
	return NewIndex(
		getter.Field,
		compute.CreateStringIndexComputation(getter),
		CreateHashTable(),
		unique,
	)
}
