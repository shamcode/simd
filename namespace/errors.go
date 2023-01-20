package namespace

import (
	"fmt"
)

type ErrRecordExists struct {
	ID int64
}

func (e ErrRecordExists) Error() string {
	return fmt.Sprintf("simd: record with passed id already exists: ID == %d", e.ID)
}

func (e ErrRecordExists) Is(err error) bool {
	_, ok := err.(ErrRecordExists)
	return ok
}

func NewErrRecordExists(id int64) error {
	return ErrRecordExists{ID: id}
}
