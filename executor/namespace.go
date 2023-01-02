package executor

import (
	"context"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type Namespace interface {
	Get(id int64) record.Record
	Insert(item record.Record) error
	Delete(id int64) error
	Upsert(item record.Record) error
	SelectForExecutor(conditions where.Conditions) ([]record.Record, error)
}

type QueryExecutor interface {
	FetchTotal(ctx context.Context, q query.Query) (int, error)
	FetchAll(ctx context.Context, q query.Query) (Iterator, error)
	FetchAllAndTotal(ctx context.Context, q query.Query) (Iterator, int, error)
}

type Iterator interface {
	Next(ctx context.Context) bool
	Item() record.Record
	Size() int
	Err() error
}
