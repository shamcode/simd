package main

import (
	"strconv"

	"github.com/shamcode/simd/record"
)

type Status uint8

const (
	StatusActive Status = iota + 1
	StatusDisabled
)

func (s Status) GoString() string { return s.String() }

func (s Status) String() string {
	switch s {
	case StatusActive:
		return "active"
	case StatusDisabled:
		return "disabled"
	}
	return strconv.Itoa(int(s))
}

type User struct {
	ID     int64
	Name   string
	Status Status
	Score  int64
}

func (u *User) GetID() int64 { return u.ID }

var userFields = record.NewFields()

var id = record.NewIDGetter[*User]()

var name = record.ComparableGetter[*User, string]{
	Field: userFields.New("name"),
	Get:   func(item *User) string { return item.Name },
}

var status = record.ComparableGetter[*User, Status]{
	Field: userFields.New("status"),
	Get:   func(item *User) Status { return item.Status },
}
