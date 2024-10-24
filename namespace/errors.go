package namespace

import (
	"fmt"
)

type RecordAlreadyExistsError struct {
	ID int64
}

func (e RecordAlreadyExistsError) Error() string {
	return fmt.Sprintf("simd: record with passed id already exists: ID == %d", e.ID)
}

func (e RecordAlreadyExistsError) Is(err error) bool {
	_, ok := err.(RecordAlreadyExistsError)
	return ok
}

func NewRecordAlreadyExists(id int64) error {
	return RecordAlreadyExistsError{ID: id}
}
