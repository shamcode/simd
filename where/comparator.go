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

type FieldComparator interface {
	record.Comparator
	GetField() string
	GetType() ComparatorType
}
