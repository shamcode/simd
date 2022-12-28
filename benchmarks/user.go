package benchmarks

import (
	"github.com/shamcode/simd/record"
)

type StatusEnum uint8

const (
	StatusActive StatusEnum = iota + 1
	StatusDisabled
)

func (s StatusEnum) Value() uint8 { return uint8(s) }
func (s StatusEnum) String() string {
	switch s {
	case StatusActive:
		return "ACTIVE"
	case StatusDisabled:
		return "DISABLED"
	default:
		return ""
	}
}

type User struct {
	ID       int64
	Name     string
	Status   StatusEnum
	Score    int
	IsOnline bool
}

func (u *User) GetID() int64   { return u.ID }
func (u *User) ComputeFields() {}

var userID = &record.Int64Getter{
	Field: "id",
	Get:   func(item interface{}) int64 { return item.(*User).ID },
}

var userName = &record.StringGetter{
	Field: "name",
	Get:   func(item interface{}) string { return item.(*User).Name },
}

var userStatus = &record.Enum8Getter{
	Field: "status",
	Get:   func(item interface{}) record.Enum8 { return item.(*User).Status },
}

var userIsOnline = &record.BoolGetter{
	Field: "is_online",
	Get:   func(item interface{}) bool { return item.(*User).IsOnline },
}

type byID struct{}

func (sorting *byID) CalcIndex(item record.Record) int64 { return item.(*User).ID }
