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
	"github.com/shamcode/simd/where"
)

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

	store.AddIndex(hash.NewEnum8HashIndex(status, false))

	for _, user := range []*User{
		{
			ID:     1,
			Name:   "Foo",
			Status: StatusActive,
			Score:  10,
		},
		{
			ID:     2,
			Name:   "Bar",
			Status: StatusDisabled,
			Score:  5,
		},
		{
			ID:     3,
			Name:   "Faz",
			Status: StatusActive,
			Score:  30,
		},
	} {
		err := store.Insert(user)
		if nil != err {
			log.Fatal(err)
		}
	}

	query := NewUserQueryBuilder(queryBuilder()).
		WhereStatus(where.EQ, StatusActive).
		Not().
		WhereName(where.EQ, "Foo").
		Query()

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
