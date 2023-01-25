package benchmarks

import (
	"context"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/indexes/btree"
	"github.com/shamcode/simd/indexes/hash"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"strconv"
	"testing"
)

type discardLogger struct{}

func (_ discardLogger) Println(...interface{}) {}

func Benchmark_Indexes(b *testing.B) {
	storeWithoutIndexes := namespace.CreateNamespace()
	storeWithoutIndexes.SetLogger(discardLogger{})

	storeWithHash := namespace.CreateNamespace()
	storeWithHash.AddIndex(hash.NewInt64HashIndex(userID, false))
	storeWithHash.AddIndex(hash.NewInt64HashIndex(userAge, false))

	storeWithHashUnique := namespace.CreateNamespace()
	storeWithHashUnique.AddIndex(hash.NewInt64HashIndex(userID, true))
	storeWithHashUnique.AddIndex(hash.NewInt64HashIndex(userAge, false))

	storeWithBtree := namespace.CreateNamespace()
	storeWithBtree.AddIndex(btree.NewInt64BTreeIndex(userID, 64, false))
	storeWithBtree.AddIndex(btree.NewInt64BTreeIndex(userAge, 8, false))

	storeWithBtreeUnique := namespace.CreateNamespace()
	storeWithBtreeUnique.AddIndex(btree.NewInt64BTreeIndex(userID, 64, true))
	storeWithBtreeUnique.AddIndex(btree.NewInt64BTreeIndex(userAge, 8, false))

	for i := 1; i < 10_000; i++ {
		for _, store := range []*namespace.WithIndexes{
			storeWithoutIndexes,
			storeWithHash,
			storeWithHashUnique,
			storeWithBtree,
			storeWithBtreeUnique,
		} {
			err := store.Upsert(&User{
				ID:       int64(i),
				Name:     "user_" + strconv.Itoa(i),
				Status:   StatusEnum(1 + i%2),
				Age:      int64(1 + i%100),
				Score:    i % 150,
				IsOnline: 0 == i%2,
			})
			if nil != err {
				b.Fatal(err)
			}
		}
	}

	benchmarks := []struct {
		Name  string
		Query query.Query
	}{
		{
			Name:  "id = 500",
			Query: query.NewBuilder(query.WhereInt64(userID, where.EQ, 500)).Query(),
		},
		{
			Name:  "id IN (500, 1000, 1500)",
			Query: query.NewBuilder(query.WhereInt64(userID, where.InArray, 500, 1000, 1500)).Query(),
		},
		{
			Name:  "id <= 1000",
			Query: query.NewBuilder(query.WhereInt64(userID, where.LE, 1000)).Query(),
		},
		{
			Name:  "id > 1000",
			Query: query.NewBuilder(query.WhereInt64(userID, where.GT, 1000)).Query(),
		},
		{
			Name:  "id <= 5000",
			Query: query.NewBuilder(query.WhereInt64(userID, where.LE, 5000)).Query(),
		},
		{
			Name:  "id > 5000",
			Query: query.NewBuilder(query.WhereInt64(userID, where.GT, 5000)).Query(),
		},
		{
			Name: "id > 2000 and id < 3000",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.GT, 2000),
				query.WhereInt64(userID, where.LT, 3000),
			).Query(),
		},
		{
			Name: "id < 2000 or id > 8000",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.LT, 2000),
				query.Or(),
				query.WhereInt64(userID, where.GT, 8000),
			).Query(),
		},
		{
			Name: "id < 1000 limit 100 asc",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.LT, 1000),
				query.Limit(100),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
		},
		{
			Name: "id < 1000 limit 100 desc",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.LT, 1000),
				query.Limit(100),
				query.Sort(sort.ByInt64IndexDesc(&byID{})),
			).Query(),
		},
		{
			Name: "id < 1000 limit 100 offset 50 asc",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.LT, 1000),
				query.Limit(100),
				query.Offset(50),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
		},
		{
			Name: "id < 1000 limit 100 offset 50 desc",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.LT, 1000),
				query.Limit(100),
				query.Offset(50),
				query.Sort(sort.ByInt64IndexDesc(&byID{})),
			).Query(),
		},
		{
			Name:  "age < 18",
			Query: query.NewBuilder(query.WhereInt64(userAge, where.LT, 18)).Query(),
		},
		{
			Name:  "age <= 17",
			Query: query.NewBuilder(query.WhereInt64(userAge, where.LE, 17)).Query(),
		},
		{
			Name: "age > 18 and age < 45",
			Query: query.NewBuilder(
				query.WhereInt64(userAge, where.GT, 18),
				query.WhereInt64(userAge, where.LT, 45),
			).Query(),
		},
		{
			Name: "age > 18 and age < 45 and id > 2000",
			Query: query.NewBuilder(
				query.WhereInt64(userAge, where.GT, 18),
				query.WhereInt64(userAge, where.LT, 45),
				query.WhereInt64(userID, where.GT, 2000),
			).Query(),
		},
	}

	executors := []struct {
		name string
		qe   executor.QueryExecutor
	}{
		{
			name: "without",
			qe:   executor.CreateQueryExecutor(storeWithoutIndexes),
		},
		{
			name: "hash",
			qe:   executor.CreateQueryExecutor(storeWithHash),
		},
		{
			name: "hash unique",
			qe:   executor.CreateQueryExecutor(storeWithHashUnique),
		},
		{
			name: "btree",
			qe:   executor.CreateQueryExecutor(storeWithBtree),
		},
		{
			name: "btree unique",
			qe:   executor.CreateQueryExecutor(storeWithBtreeUnique),
		},
	}

	for _, bench := range benchmarks {
		b.Run(bench.Name, func(b *testing.B) {
			for _, exec := range executors {
				b.Run(exec.name, func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						_, _, err := exec.qe.FetchAllAndTotal(context.Background(), bench.Query)
						if nil != err {
							b.Fatal(err)
						}
					}
				})
			}
		})
	}
}

func Benchmark_BTreeIndexesMaxChildren(b *testing.B) {
	maxChildren := 64
	stores := make([]*namespace.WithIndexes, 0, maxChildren)
	for i := 1; i <= maxChildren; i++ {
		store := namespace.CreateNamespace()
		store.AddIndex(btree.NewInt64BTreeIndex(userAge, i, false))
		stores = append(stores, store)
	}

	for i := 1; i < 10_000; i++ {
		for _, store := range stores {
			err := store.Upsert(&User{
				ID:       int64(i),
				Name:     "user_" + strconv.Itoa(i),
				Status:   StatusEnum(1 + i%2),
				Age:      int64(i%100 + 1),
				Score:    i % 150,
				IsOnline: 0 == i%2,
			})
			if nil != err {
				b.Fatal(err)
			}
		}
	}

	q := query.NewBuilder(query.WhereInt64(userAge, where.LT, 5)).Query()

	for i, store := range stores {
		exec := executor.CreateQueryExecutor(store)
		b.Run(strconv.Itoa(i+1), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, err := exec.FetchAllAndTotal(context.Background(), q)
				if nil != err {
					b.Fatal(err)
				}
			}
		})
	}
}
