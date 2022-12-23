package benchmarks

import (
	"context"
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/indexes/fields"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"strconv"
	"testing"
)

func Benchmark_FetchAllAndTotal(b *testing.B) {
	store := indexes.CreateNamespace()
	store.AddIndex(fields.NewInt64Index(userID))
	store.AddIndex(fields.NewStringIndex(userName))
	store.AddIndex(fields.NewEnum8Index(userStatus))
	store.AddIndex(fields.NewBoolIndex(userIsOnline))

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
		Query namespace.Query
	}{
		{
			Name:  "is_online = true",
			Query: store.Query().WhereBool(userIsOnline, where.EQ, true),
		},
		{
			Name: "is_online = true offset 1000 limit 100",
			Query: store.Query().
				WhereBool(userIsOnline, where.EQ, true).
				Offset(1000).
				Limit(100),
		},
		{
			Name:  "id > 1000",
			Query: store.Query().WhereInt64(userID, where.GT, 1000),
		},
		{
			Name: "id > 1000 limit 100 asc",
			Query: store.Query().
				WhereInt64(userID, where.GT, 1000).
				Limit(100).
				Sort(sort.ByInt64Index(&byIDAsc{})),
		},
		{
			Name: "id > 1000 limit 100 desc",
			Query: store.Query().
				WhereInt64(userID, where.GT, 1000).
				Limit(100).
				Sort(sort.ByInt64Index(&byIDDesc{})),
		},
		{
			Name: "id > 1000 and is_online = true and status = ACTIVE",
			Query: store.Query().
				WhereInt64(userID, where.GT, 1000).
				WhereBool(userIsOnline, where.EQ, true).
				WhereEnum8(userStatus, where.EQ, StatusActive),
		},
		{
			Name: "id > 1000 and is_online = true and status = ACTIVE limit 100 desc",
			Query: store.Query().
				WhereInt64(userID, where.GT, 1000).
				WhereBool(userIsOnline, where.EQ, true).
				WhereEnum8(userStatus, where.EQ, StatusActive).
				Sort(sort.ByInt64Index(&byIDDesc{})).
				Limit(100),
		},
	}

	for _, bench := range benchmarks {
		b.Run(bench.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, err := bench.Query.FetchAllAndTotal(context.Background())
				if nil != err {
					b.Fatal(err)
				}
			}
		})
	}
}
