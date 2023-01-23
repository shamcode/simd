package benchmarks

import (
	"context"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/indexes/hash"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"strconv"
	"testing"
)

func Benchmark_FetchAllAndTotal(b *testing.B) {
	store := namespace.CreateNamespace()
	store.AddIndex(hash.NewInt64HashIndex(userID))
	store.AddIndex(hash.NewStringHashIndex(userName))
	store.AddIndex(hash.NewEnum8HashIndex(userStatus))
	store.AddIndex(hash.NewBoolHashIndex(userIsOnline))

	for i := 1; i < 10_000; i++ {
		err := store.Upsert(&User{
			ID:       int64(i),
			Name:     "user_" + strconv.Itoa(i),
			Status:   StatusEnum(1 + i%2),
			Score:    i % 150,
			IsOnline: 0 == i%2,
		})
		if nil != err {
			b.Fatal(err)
		}
	}

	benchmarks := []struct {
		Name  string
		Query query.Query
	}{
		{
			Name:  "is_online = true",
			Query: query.NewBuilder(query.WhereBool(userIsOnline, where.EQ, true)).Query(),
		},
		{
			Name: "is_online = true offset 1000 limit 100",
			Query: query.NewBuilder(
				query.WhereBool(userIsOnline, where.EQ, true),
				query.Offset(1000),
				query.Limit(100),
			).Query(),
		},
		{
			Name:  "id > 1000",
			Query: query.NewBuilder(query.WhereInt64(userID, where.GT, 1000)).Query(),
		},
		{
			Name: "id > 1000 limit 100 asc",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.GT, 1000),
				query.Limit(100),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
		},
		{
			Name: "id > 1000 limit 100 desc",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.GT, 1000),
				query.Limit(100),
				query.Sort(sort.ByInt64IndexDesc(&byID{})),
			).Query(),
		},
		{
			Name: "id > 1000 and is_online = true and status = ACTIVE",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.GT, 1000),
				query.WhereBool(userIsOnline, where.EQ, true),
				query.WhereEnum8(userStatus, where.EQ, StatusActive),
			).Query(),
		},
		{
			Name: "id > 1000 and is_online = true and status = ACTIVE limit 100 desc",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.GT, 1000),
				query.WhereBool(userIsOnline, where.EQ, true),
				query.WhereEnum8(userStatus, where.EQ, StatusActive),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
				query.Limit(100),
			).Query(),
		},
	}

	qe := executor.CreateQueryExecutor(store)
	for _, bench := range benchmarks {
		b.Run(bench.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, err := qe.FetchAllAndTotal(context.Background(), bench.Query)
				if nil != err {
					b.Fatal(err)
				}
			}
		})
	}
}
