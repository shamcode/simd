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
	Age      int64
	Score    int
	IsOnline bool
}

func (u *User) GetID() int64   { return u.ID }
func (u *User) ComputeFields() {}

var userFields = record.NewFields()

var userName = record.StringGetter{
	Field: userFields.New("name"),
	Get:   func(item record.Record) string { return item.(*User).Name },
}

var userStatus = record.Enum8Getter{
	Field: userFields.New("status"),
	Get:   func(item record.Record) record.Enum8 { return item.(*User).Status },
}

var userAge = record.Int64Getter{
	Field: userFields.New("age"),
	Get:   func(item record.Record) int64 { return item.(*User).Age },
}

var userIsOnline = record.BoolGetter{
	Field: userFields.New("is_online"),
	Get:   func(item record.Record) bool { return item.(*User).IsOnline },
}
