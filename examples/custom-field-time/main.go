package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/shamcode/simd/debug"
	"github.com/shamcode/simd/examples/custom-field-time/types"
	indexesByType "github.com/shamcode/simd/examples/custom-field-time/types/indexes"
	"github.com/shamcode/simd/examples/custom-field-time/types/querybuilder"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
)

type Item struct {
	ID       int64
	CreateAt time.Time
}

func (u *Item) GetID() int64 { return u.ID }

var itemFields = record.NewFields()

var id = record.NewIDGetter[*Item]()

var createdAt = types.TimeGetter[*Item]{
	Field: itemFields.New("created_at"),
	Get: func(item *Item) time.Time {
		return item.CreateAt
	},
}

func main() {
	debugEnabled := flag.Bool("debug", false, "enabled debug")
	flag.Parse()

	store := namespace.CreateNamespace[*Item]()

	queryBuilder := query.NewBuilder[*Item]
	queryExecutor := executor.CreateQueryExecutor(store)

	if *debugEnabled {
		queryBuilder = debug.WrapCreateQueryBuilder(queryBuilder)
		queryExecutor = debug.WrapQueryExecutor(queryExecutor, func(s string) {
			log.Printf("SIMD QUERY: %s", s)
		})
	}

	store.AddIndex(indexesByType.NewTimeBTreeIndex(createdAt, 8, false))

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

	query := queryBuilder(
		querybuilder.WhereTime(createdAt, where.LT, time.Date(2022, time.January, 1, 0, 0, 0, 0, time.Local)),
		query.Sort(sort.Asc(id)),
	).Query()

	ctx := context.Background()
	cur, total, err := queryExecutor.FetchAllAndTotal(ctx, query)
	if nil != err {
		log.Fatal(err)
	}
	for cur.Next(ctx) {
		log.Printf("%#v", cur.Item())
	}
	if err := cur.Err(); nil != err {
		log.Fatal(err)
	}
	log.Printf("total: %d", total)
}
