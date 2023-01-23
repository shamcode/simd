package main

import "github.com/shamcode/simd/record"

type Status uint8

func (s Status) Value() uint8 { return uint8(s) }

const (
	StatusActive Status = iota + 1
	StatusDisabled
)

type User struct {
	ID     int64
	Name   string
	Status Status
	Score  int64
}

func (u *User) GetID() int64   { return u.ID }
func (u *User) ComputeFields() {}

var id = &record.Int64Getter{
	Field: "id",
	Get:   func(item record.Record) int64 { return item.(*User).ID },
}

var name = &record.StringGetter{
	Field: "name",
	Get:   func(item record.Record) string { return item.(*User).Name },
}

var status = &record.Enum8Getter{
	Field: "status",
	Get:   func(item record.Record) record.Enum8 { return item.(*User).Status },
}
