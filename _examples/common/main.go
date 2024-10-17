package main

import (
	"context"
	"flag"
	"log"

	"github.com/shamcode/simd/debug"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/indexes/hash"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
)

type User struct {
	ID   int64
	Name string
}

func (u *User) GetID() int64 { return u.ID }

var userFields = record.NewFields()

var name = record.StringGetter{
	Field: userFields.New("name"),
	Get:   func(item record.Record) string { return item.(*User).Name },
}

func main() {
	debugEnabled := flag.Bool("debug", false, "enabled debug")
	flag.Parse()

	store := namespace.CreateNamespace()
	queryBuilder := query.NewBuilder
	queryExecutor := executor.CreateQueryExecutor(store)

	// if debug enabled, add logging for query
	if *debugEnabled {
		queryBuilder = debug.WrapCreateQueryBuilder(queryBuilder)
		queryExecutor = debug.WrapQueryExecutor(queryExecutor, func(s string) {
			log.Printf("SIMD QUERY: %s", s)
		})
	}

	store.AddIndex(hash.NewInt64HashIndex(record.ID, true))
	store.AddIndex(hash.NewStringHashIndex(name, false))

	for _, user := range []*User{
		{
			ID:   1,
			Name: "Foo",
		},
		{
			ID:   2,
			Name: "Bar",
		},
		{
			ID:   3,
			Name: "Faz",
		},
	} {
		err := store.Insert(user)
		if nil != err {
			log.Fatal(err)
		}
	}

	query := queryBuilder(
		query.WhereInt64(record.ID, where.GT, 1),
		query.Sort(sort.Asc(name)),
	).Query()

	ctx := context.Background()
	cur, total, err := queryExecutor.FetchAllAndTotal(ctx, query)
	if nil != err {
		log.Fatal(err)
	}
	for cur.Next(ctx) {
		log.Printf("%#v", cur.Item().(*User))
	}
	if err := cur.Err(); nil != err {
		log.Fatal(err)
	}
	log.Printf("total: %d", total)
}
