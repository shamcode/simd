package namespace

import "errors"

var (
	ErrRecordExists = errors.New("simd: record with passed id already exists")
)
