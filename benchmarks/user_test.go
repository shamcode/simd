package benchmarks

import (
	"github.com/shamcode/simd/record"
)

type StatusEnum uint8

const (
	StatusActive StatusEnum = iota + 1
	StatusDisabled
)

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

func (u *User) GetID() int64 { return u.ID }

var userFields = record.NewFields()

var userID = record.NewIDGetter[*User]()

var userName = record.ComparableGetter[*User, string]{
	Field: userFields.New("name"),
	Get:   func(item *User) string { return item.Name },
}

var userStatus = record.ComparableGetter[*User, StatusEnum]{
	Field: userFields.New("status"),
	Get:   func(item *User) StatusEnum { return item.Status },
}

var userAge = record.ComparableGetter[*User, int64]{
	Field: userFields.New("age"),
	Get:   func(item *User) int64 { return item.Age },
}

var userIsOnline = record.BoolGetter[*User]{
	Field: userFields.New("is_online"),
	Get:   func(item *User) bool { return item.IsOnline },
}
