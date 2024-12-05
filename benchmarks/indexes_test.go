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

type discardLogger struct{}

func (discardLogger) Println(...any) {}

func Benchmark_Indexes(b *testing.B) {
	storeWithoutIndexes := namespace.CreateNamespace[*User]()
	storeWithoutIndexes.SetLogger(discardLogger{})

	storeWithHash := namespace.CreateNamespace[*User]()
	storeWithHash.AddIndex(hash.NewComparableHashIndex(userID, false))
	storeWithHash.AddIndex(hash.NewComparableHashIndex(userAge, false))

	storeWithHashUnique := namespace.CreateNamespace[*User]()
	storeWithHashUnique.AddIndex(hash.NewComparableHashIndex(userID, true))
	storeWithHashUnique.AddIndex(hash.NewComparableHashIndex(userAge, false))

	storeWithBtree := namespace.CreateNamespace[*User]()
	storeWithBtree.AddIndex(btree.NewComparableBTreeIndex(userID, 64, false))
	storeWithBtree.AddIndex(btree.NewComparableBTreeIndex(userAge, 8, false))

	storeWithBtreeUnique := namespace.CreateNamespace[*User]()
	storeWithBtreeUnique.AddIndex(btree.NewComparableBTreeIndex(userID, 64, true))
	storeWithBtreeUnique.AddIndex(btree.NewComparableBTreeIndex(userAge, 8, false))

	for i := 1; i < 10_000; i++ {
		for _, store := range []*namespace.WithIndexes[*User]{
			storeWithoutIndexes,
			storeWithHash,
			storeWithHashUnique,
			storeWithBtree,
			storeWithBtreeUnique,
		} {
			err := store.Upsert(&User{
				ID:       int64(i),
				Name:     "user_" + strconv.Itoa(i),
				Status:   StatusEnum(1 + i%2), //nolint:gosec
				Age:      int64(1 + i%100),
				Score:    i % 150,
				IsOnline: i%2 == 0,
			})
			if nil != err {
				b.Fatal(err)
			}
		}
	}

	benchmarks := []struct {
		Name  string
		Query query.Query[*User]
	}{
		{
			Name:  "id = 500",
			Query: query.NewBuilder[*User](query.Where(userID, where.EQ, 500)).Query(),
		},
		{
			Name:  "id IN (500, 1000, 1500)",
			Query: query.NewBuilder[*User](query.Where(userID, where.InArray, 500, 1000, 1500)).Query(),
		},
		{
			Name:  "id <= 1000",
			Query: query.NewBuilder[*User](query.Where(userID, where.LE, 1000)).Query(),
		},
		{
			Name:  "id > 1000",
			Query: query.NewBuilder[*User](query.Where(userID, where.GT, 1000)).Query(),
		},
		{
			Name:  "id <= 5000",
			Query: query.NewBuilder[*User](query.Where(userID, where.LE, 5000)).Query(),
		},
		{
			Name:  "id > 5000",
			Query: query.NewBuilder[*User](query.Where(userID, where.GT, 5000)).Query(),
		},
		{
			Name: "id > 2000 and id < 3000",
			Query: query.NewBuilder[*User](
				query.Where(userID, where.GT, 2000),
				query.Where(userID, where.LT, 3000),
			).Query(),
		},
		{
			Name: "id < 2000 or id > 8000",
			Query: query.NewBuilder[*User](
				query.Where(userID, where.LT, 2000),
				query.Or(),
				query.Where(userID, where.GT, 8000),
			).Query(),
		},
		{
			Name: "id < 1000 limit 100 asc",
			Query: query.NewBuilder[*User](
				query.Where(userID, where.LT, 1000),
				query.Limit(100),
				query.Sort(sort.Asc[*User](userID)),
			).Query(),
		},
		{
			Name: "id < 1000 limit 100 desc",
			Query: query.NewBuilder[*User](
				query.Where(userID, where.LT, 1000),
				query.Limit(100),
				query.Sort(sort.Desc[*User](userID)),
			).Query(),
		},
		{
			Name: "id < 1000 limit 100 offset 50 asc",
			Query: query.NewBuilder[*User](
				query.Where(userID, where.LT, 1000),
				query.Limit(100),
				query.Offset(50),
				query.Sort(sort.Asc[*User](userID)),
			).Query(),
		},
		{
			Name: "id < 1000 limit 100 offset 50 desc",
			Query: query.NewBuilder[*User](
				query.Where(userID, where.LT, 1000),
				query.Limit(100),
				query.Offset(50),
				query.Sort(sort.Desc[*User](userID)),
			).Query(),
		},
		{
			Name:  "age < 18",
			Query: query.NewBuilder[*User](query.Where(userAge, where.LT, 18)).Query(),
		},
		{
			Name:  "age <= 17",
			Query: query.NewBuilder[*User](query.Where(userAge, where.LE, 17)).Query(),
		},
		{
			Name: "age > 18 and age < 45",
			Query: query.NewBuilder[*User](
				query.Where(userAge, where.GT, 18),
				query.Where(userAge, where.LT, 45),
			).Query(),
		},
		{
			Name: "age > 18 and age < 45 and id > 2000",
			Query: query.NewBuilder[*User](
				query.Where(userAge, where.GT, 18),
				query.Where(userAge, where.LT, 45),
				query.Where(userID, where.GT, 2000),
			).Query(),
		},
	}

	executors := []struct {
		name string
		qe   executor.QueryExecutor[*User]
	}{
		{
			name: "without",
			qe:   executor.CreateQueryExecutor[*User](storeWithoutIndexes),
		},
		{
			name: "hash",
			qe:   executor.CreateQueryExecutor[*User](storeWithHash),
		},
		{
			name: "hash unique",
			qe:   executor.CreateQueryExecutor[*User](storeWithHashUnique),
		},
		{
			name: "btree",
			qe:   executor.CreateQueryExecutor[*User](storeWithBtree),
		},
		{
			name: "btree unique",
			qe:   executor.CreateQueryExecutor[*User](storeWithBtreeUnique),
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
	stores := make([]*namespace.WithIndexes[*User], 0, maxChildren)
	for i := 1; i <= maxChildren; i++ {
		store := namespace.CreateNamespace[*User]()
		store.AddIndex(btree.NewComparableBTreeIndex(userAge, i, false))
		stores = append(stores, store)
	}

	for i := 1; i < 10_000; i++ {
		for _, store := range stores {
			err := store.Upsert(&User{
				ID:       int64(i),
				Name:     "user_" + strconv.Itoa(i),
				Status:   StatusEnum(1 + i%2), //nolint:gosec
				Age:      int64(i%100 + 1),
				Score:    i % 150,
				IsOnline: i%2 == 0,
			})
			if nil != err {
				b.Fatal(err)
			}
		}
	}

	q := query.NewBuilder[*User](query.Where(userAge, where.LT, 5)).Query()

	for i, store := range stores {
		exec := executor.CreateQueryExecutor[*User](store)
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
