package main

import (
	"github.com/shamcode/simd/record"
	"strconv"
)

type Status uint8

func (s Status) Value() uint8 { return uint8(s) }

const (
	StatusActive Status = iota + 1
	StatusDisabled
)

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

func (u *User) GetID() int64   { return u.ID }
func (u *User) ComputeFields() {}

var userFields = record.NewFields()

var id = &record.Int64Getter{
	Field: userFields.New("id"),
	Get:   func(item record.Record) int64 { return item.(*User).ID },
}

var name = &record.StringGetter{
	Field: userFields.New("name"),
	Get:   func(item record.Record) string { return item.(*User).Name },
}

var status = &record.Enum8Getter{
	Field: userFields.New("status"),
	Get:   func(item record.Record) record.Enum8 { return item.(*User).Status },
}
