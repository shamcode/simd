package main

import (
	"context"
	"flag"
	"github.com/shamcode/simd/debug"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/indexes/bytype"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"log"
)

type User struct {
	ID   int64
	Name string
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

	store.AddIndex(bytype.NewInt64Index(id))
	store.AddIndex(bytype.NewStringIndex(name))

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

	q := queryBuilder(
		query.WhereInt64(id, where.GT, 1),
		query.Sort(sort.ByStringAsc(name)),
	).Query()

	ctx := context.Background()
	cur, total, err := queryExecutor.FetchAllAndTotal(ctx, q)
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
