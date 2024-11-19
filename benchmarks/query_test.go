//nolint:gosec
package benchmarks

import (
	"context"
	"strconv"
	"testing"

	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/indexes/btree"
	"github.com/shamcode/simd/indexes/hash"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
)

func Benchmark_Query(b *testing.B) {
	store := namespace.CreateNamespace[*User]()
	store.AddIndex(hash.NewComparableHashIndex(userID, true))
	store.AddIndex(btree.NewComparableBTreeIndex(userID, 8, true))
	store.AddIndex(hash.NewComparableHashIndex(userName, false))
	store.AddIndex(hash.NewComparableHashIndex(userStatus, false))
	store.AddIndex(hash.NewBoolHashIndex(userIsOnline, false))

	for i := 1; i < 10_000; i++ {
		err := store.Upsert(&User{ //nolint:exhaustruct
			ID:       int64(i),
			Name:     "user_" + strconv.Itoa(i),
			Status:   StatusEnum(1 + i%2),
			Score:    i % 150,
			IsOnline: i%2 == 0,
		})
		if nil != err {
			b.Fatal(err)
		}
	}

	b.Run("FetchAllAndTotal", func(b *testing.B) {
		benchmarks := []struct {
			Name  string
			Query query.Query[*User]
		}{
			{
				Name:  "is_online = true",
				Query: query.NewBuilder[*User](query.WhereBool(userIsOnline, where.EQ, true)).Query(),
			},
			{
				Name: "is_online = true offset 1000 limit 100",
				Query: query.NewBuilder[*User](
					query.WhereBool(userIsOnline, where.EQ, true),
					query.Offset(1000),
					query.Limit(100),
				).Query(),
			},
			{
				Name:  "id <= 5000",
				Query: query.NewBuilder[*User](query.Where(userID, where.LE, 5000)).Query(),
			},
			{
				Name:  "id > 1000",
				Query: query.NewBuilder[*User](query.Where(userID, where.GT, 1000)).Query(),
			},
			{
				Name: "id > 1000 limit 100 asc",
				Query: query.NewBuilder[*User](
					query.Where(userID, where.GT, 1000),
					query.Limit(100),
					query.Sort(sort.Asc[*User](userID)),
				).Query(),
			},
			{
				Name: "id > 1000 limit 100 desc",
				Query: query.NewBuilder[*User](
					query.Where(userID, where.GT, 1000),
					query.Limit(100),
					query.Sort(sort.Desc[*User](userID)),
				).Query(),
			},
			{
				Name: "id > 1000 and is_online = true and status = ACTIVE",
				Query: query.NewBuilder[*User](
					query.Where(userID, where.GT, 1000),
					query.WhereBool(userIsOnline, where.EQ, true),
					query.Where(userStatus, where.EQ, StatusActive),
				).Query(),
			},
			{
				Name: "id > 1000 and is_online = true and status = ACTIVE limit 100 desc",
				Query: query.NewBuilder[*User](
					query.Where(userID, where.GT, 1000),
					query.WhereBool(userIsOnline, where.EQ, true),
					query.Where(userStatus, where.EQ, StatusActive),
					query.Sort(sort.Asc[*User](userID)),
					query.Limit(100),
				).Query(),
			},
		}

		qe := executor.CreateQueryExecutor[*User](store)
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
	})
}
