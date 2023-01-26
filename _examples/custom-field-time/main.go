package main

import (
	"context"
	"flag"
	"github.com/shamcode/simd/_examples/custom-field-time/types"
	indexesByType "github.com/shamcode/simd/_examples/custom-field-time/types/indexes"
	"github.com/shamcode/simd/_examples/custom-field-time/types/querybuilder"
	"github.com/shamcode/simd/debug"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"log"
	"time"
)

type Item struct {
	ID       int64
	CreateAt time.Time
}

func (u *Item) GetID() int64   { return u.ID }
func (u *Item) ComputeFields() {}

var itemFields = record.NewFields()

var createdAt = &types.TimeGetter{
	Field: itemFields.New("created_at"),
	Get: func(item record.Record) time.Time {
		return item.(*Item).CreateAt
	},
}

type byID struct{}

func (sorting *byID) CalcIndex(item record.Record) int64 { return item.(*Item).ID }

func main() {
	debugEnabled := flag.Bool("debug", false, "enabled debug")
	flag.Parse()

	store := namespace.CreateNamespace()

	queryBuilder := query.NewBuilder
	queryExecutor := executor.CreateQueryExecutor(store)

	if *debugEnabled {
		queryBuilder = debug.WrapCreateQueryBuilder(queryBuilder)
		queryExecutor = debug.WrapQueryExecutor(queryExecutor, func(s string) {
			log.Printf("SIMD QUERY: %s", s)
		})
	}

	store.AddIndex(indexesByType.NewTimeHashIndex(createdAt, false))

	for _, user := range []*Item{
		{
			ID:       1,
			CreateAt: time.Date(2022, time.December, 28, 22, 58, 0, 0, time.Local),
		},
		{
			ID:       2,
			CreateAt: time.Date(2021, time.December, 28, 22, 58, 0, 0, time.Local),
		},
		{
			ID:       3,
			CreateAt: time.Date(2020, time.December, 28, 22, 58, 0, 0, time.Local),
		},
	} {
		err := store.Insert(user)
		if nil != err {
			log.Fatal(err)
		}
	}

	q := queryBuilder(
		querybuilder.WhereTime(createdAt, where.LT, time.Date(2022, time.January, 1, 0, 0, 0, 0, time.Local)),
		query.Sort(sort.ByInt64IndexAsc(&byID{})),
	).Query()

	ctx := context.Background()
	cur, total, err := queryExecutor.FetchAllAndTotal(ctx, q)
	if nil != err {
		log.Fatal(err)
	}
	for cur.Next(ctx) {
		log.Printf("%#v", cur.Item().(*Item))
	}
	if err := cur.Err(); nil != err {
		log.Fatal(err)
	}
	log.Printf("total: %d", total)
}
