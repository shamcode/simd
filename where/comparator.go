package where

import "github.com/shamcode/simd/record"

type ComparatorType uint8

const (
	EQ ComparatorType = iota + 1
	GT
	LT
	GE
	LE
	InArray
	Like
	Regexp
	SetHas
	MapHasValue
	MapHasKey
)

type FieldComparator[R record.Record] interface {
	GetField() record.Field
	GetType() ComparatorType
	Compare(item R) (bool, error)
	ValuesCount() int
	ValueAt(index int) interface{}
}
