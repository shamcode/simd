package debug

import (
	"context"
	"fmt"
	"github.com/shamcode/simd/asserts"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"strings"
	"testing"
)

type user struct {
	ID   int64
	Name string
	Age  int
}

func (u *user) GetID() int64 {
	return u.ID
}

func (u *user) ComputeFields() {}

var userFields = record.NewFields()

var id = &record.Int64Getter{
	Field: userFields.New("id"),
	Get: func(item record.Record) int64 {
		return item.(*user).ID
	},
}

var name = &record.StringGetter{
	Field: userFields.New("name"),
	Get: func(item record.Record) string {
		return item.(*user).Name
	},
}

var age = &record.IntGetter{
	Field: userFields.New("age"),
	Get: func(item record.Record) int {
		return item.(*user).Age
	},
}

type storage struct {
	data map[int64]record.Record
}

func (s *storage) Get(id int64) record.Record {
	return s.data[id]
}

func (s *storage) Insert(item record.Record) error {
	s.insert(item)
	return nil
}

func (s *storage) insert(item record.Record) {
	s.data[item.GetID()] = item
}

func (s *storage) Delete(id int64) error {
	delete(s.data, id)
	return nil
}

func (s *storage) Upsert(item record.Record) error {
	s.data[item.GetID()] = item
	return nil
}

func (s *storage) PreselectForExecutor(_ where.Conditions) ([]record.Record, error) {
	items := make([]record.Record, 0, len(s.data))
	for _, item := range s.data {
		items = append(items, item)
	}
	return items, nil
}

type byID struct{}

func (sorting *byID) CalcIndex(item record.Record) int64 {
	return item.GetID()
}

func TestQueryExecutorWithDebug(t *testing.T) {
	ns := &storage{
		data: make(map[int64]record.Record),
	}
	ns.insert(&user{ID: 1, Name: "first", Age: 18})
	ns.insert(&user{ID: 2, Name: "second", Age: 19})
	ns.insert(&user{ID: 3, Name: "third", Age: 20})
	ns.insert(&user{ID: 4, Name: "fourth", Age: 21})
	ns.insert(&user{ID: 5, Name: "fifth", Age: 22})

	newBuilder := WrapCreateQueryBuilder(query.NewBuilder)

	tests := []struct {
		name                 string
		query                query.Query
		expected             string
		expectedErrorMessage string
	}{
		{
			name:     "order by id asc",
			query:    query.NewBuilder(query.Sort(sort.ByInt64IndexAsc(&byID{}))).Query(),
			expected: "SELECT *, COUNT(*) <Query dont implement QueryWithDumper interface, check QueryBuilder>",
		},
		{
			name:     "order by id asc",
			query:    newBuilder(query.Sort(sort.ByInt64IndexAsc(&byID{}))).Query(),
			expected: "SELECT *, COUNT(*) ORDER BY &debug.byID{} ASC",
		},
		{
			name:     "order by id desc",
			query:    newBuilder(query.Sort(sort.ByInt64IndexDesc(&byID{}))).Query(),
			expected: "SELECT *, COUNT(*) ORDER BY &debug.byID{} DESC",
		},
		{
			name:     "where id = int64(3)",
			query:    newBuilder(query.WhereInt64(id, where.EQ, 3)).Query(),
			expected: "SELECT *, COUNT(*) WHERE id = 3",
		},
		{
			name: "where id = int64(3) and age == int(20)",
			query: newBuilder(
				query.WhereInt64(id, where.EQ, 3),
				query.WhereInt(age, where.EQ, 20),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE id = 3 AND age = 20",
		},
		{
			name: "where id = int64(3) and age == int(18)",
			query: newBuilder(
				query.WhereInt64(id, where.EQ, 3),
				query.WhereInt(age, where.EQ, 18),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE id = 3 AND age = 18",
		},
		{
			name: "where age > 18 and age < 22",
			query: newBuilder(
				query.WhereInt(age, where.GT, 18),
				query.WhereInt(age, where.LT, 22),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE age > 18 AND age < 22 ORDER BY &debug.byID{} ASC",
		},
		{
			name: "where age >= 18 and age <= 22",
			query: newBuilder(
				query.WhereInt(age, where.GE, 18),
				query.WhereInt(age, where.LE, 22),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE age >= 18 AND age <= 22 ORDER BY &debug.byID{} ASC",
		},
		{
			name: "where id = 2 or id = 5",
			query: newBuilder(
				query.WhereInt64(id, where.EQ, 2),
				query.Or(),
				query.WhereInt64(id, where.EQ, 5),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE id = 2 OR id = 5 ORDER BY &debug.byID{} ASC",
		},
		{
			name: "where id = 2 or age > 20",
			query: newBuilder(
				query.WhereInt64(id, where.EQ, 2),
				query.Or(),
				query.WhereInt(age, where.GT, 20),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE id = 2 OR age > 20 ORDER BY &debug.byID{} ASC",
		},
		{
			name: "where id = 1 or ( age > 20 and age < 22)",
			query: newBuilder(
				query.WhereInt64(id, where.EQ, 1),
				query.Or(),
				query.OpenBracket(),
				query.WhereInt(age, where.GT, 20),
				query.WhereInt(age, where.LT, 22),
				query.CloseBracket(),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE id = 1 OR (age > 20 AND age < 22) ORDER BY &debug.byID{} ASC",
		},
		{
			name: "where ( age > 20 and age < 22) or id = 1",
			query: newBuilder(
				query.OpenBracket(),
				query.WhereInt(age, where.GT, 20),
				query.WhereInt(age, where.LT, 22),
				query.CloseBracket(),
				query.Or(),
				query.WhereInt64(id, where.EQ, 1),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE (age > 20 AND age < 22) OR id = 1 ORDER BY &debug.byID{} ASC",
		},
		{
			name: "where age > 20 and age < 22 or id = 1",
			query: newBuilder(
				query.WhereInt(age, where.GT, 20),
				query.WhereInt(age, where.LT, 22),
				query.Or(),
				query.WhereInt64(id, where.EQ, 1),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE age > 20 AND age < 22 OR id = 1 ORDER BY &debug.byID{} ASC",
		},
		{
			name: "where (age > 20 and age < 22) or id = 1",
			query: newBuilder(
				query.OpenBracket(),
				query.WhereInt(age, where.GT, 20),
				query.WhereInt(age, where.LT, 22),
				query.CloseBracket(),
				query.Or(),
				query.WhereInt64(id, where.EQ, 1),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE (age > 20 AND age < 22) OR id = 1 ORDER BY &debug.byID{} ASC",
		},
		{
			name: "where age > 20 and age < 22 and (id = 1 or id = 2)",
			query: newBuilder(
				query.WhereInt(age, where.GT, 20),
				query.WhereInt(age, where.LT, 22),
				query.OpenBracket(),
				query.WhereInt64(id, where.EQ, 1),
				query.Or(),
				query.WhereInt64(id, where.EQ, 2),
				query.CloseBracket(),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE age > 20 AND age < 22 AND (id = 1 OR id = 2)",
		},
		{
			name: "where age > 20 and age < 22 and (id = 1 or id = 2 or id = 4)",
			query: newBuilder(
				query.WhereInt(age, where.GT, 20),
				query.WhereInt(age, where.LT, 22),
				query.OpenBracket(),
				query.WhereInt64(id, where.EQ, 1),
				query.Or(),
				query.WhereInt64(id, where.EQ, 2),
				query.Or(),
				query.WhereInt64(id, where.EQ, 4),
				query.CloseBracket(),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE age > 20 AND age < 22 AND (id = 1 OR id = 2 OR id = 4)",
		},
		{
			name: "where (age > 20 and age < 22) and id = 4",
			query: newBuilder(
				query.OpenBracket(),
				query.WhereInt(age, where.GT, 20),
				query.WhereInt(age, where.LT, 22),
				query.CloseBracket(),
				query.WhereInt64(id, where.EQ, 4),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE (age > 20 AND age < 22) AND id = 4",
		},
		{
			name: "where age in {20, 21, 22} and id > 3",
			query: newBuilder(
				query.WhereInt(age, where.InArray, 20, 21, 22),
				query.WhereInt64(id, where.GT, 3),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE age IN (20, 21, 22) AND id > 3 ORDER BY &debug.byID{} ASC",
		},
		{
			name: "where name like \"th\"",
			query: newBuilder(
				query.WhereString(name, where.Like, "th"),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE name LIKE \"th\" ORDER BY &debug.byID{} ASC",
		},
		{
			name: "where name like \"th\" or name like \"first\"",
			query: newBuilder(
				query.WhereString(name, where.Like, "th"),
				query.Or(),
				query.WhereString(name, where.Like, "first"),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE name LIKE \"th\" OR name LIKE \"first\" ORDER BY &debug.byID{} ASC",
		},
		{
			name: "where ((id = 1) or (id = 2))",
			query: newBuilder(
				query.OpenBracket(),
				query.OpenBracket(),
				query.WhereInt64(id, where.EQ, 1),
				query.CloseBracket(),
				query.Or(),
				query.OpenBracket(),
				query.WhereInt64(id, where.EQ, 2),
				query.CloseBracket(),
				query.CloseBracket(),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE ((id = 1) OR (id = 2)) ORDER BY &debug.byID{} ASC",
		},
		{
			name: "where (((id = 1) or (id = 2)) or id = 3) or id = 4",
			query: newBuilder(
				query.OpenBracket(),
				query.OpenBracket(),
				query.OpenBracket(),
				query.WhereInt64(id, where.EQ, 1),
				query.CloseBracket(),
				query.Or(),
				query.OpenBracket(),
				query.WhereInt64(id, where.EQ, 2),
				query.CloseBracket(),
				query.CloseBracket(),
				query.Or(),
				query.WhereInt64(id, where.EQ, 3),
				query.CloseBracket(),
				query.Or(),
				query.WhereInt64(id, where.EQ, 4),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE (((id = 1) OR (id = 2)) OR id = 3) OR id = 4 ORDER BY &debug.byID{} ASC",
		},
		{
			name: "where id > 1 limit 2 offset 1 order by &debug.byID{} ASC",
			query: newBuilder(
				query.WhereInt64(id, where.GT, 1),
				query.Limit(2),
				query.Offset(1),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected: "SELECT *, COUNT(*) WHERE id > 1 ORDER BY &debug.byID{} ASC OFFSET 1 LIMIT 2",
		},
		{
			name: "unknown comparator",
			query: newBuilder(
				query.WhereInt64(id, where.ComparatorType(100), 1, 2),
				query.Limit(2),
				query.Offset(1),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected:             "SELECT *, COUNT(*) WHERE id (ComparatorType(100) 1 2) ORDER BY &debug.byID{} ASC OFFSET 1 LIMIT 2",
			expectedErrorMessage: "execute query: not implemented ComparatorType: 100, field = id",
		},
	}

	qe := executor.CreateQueryExecutor(ns)

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			_, err := WrapQueryExecutor(qe, func(q string) {
				asserts.Equals(t, test.expected, q, "query")
			}).FetchAll(ctx, test.query)

			var errMsg string
			if nil != err {
				errMsg = err.Error()
			}
			asserts.Equals(t, test.expectedErrorMessage, errMsg, "error")
		})
	}
}

func TestFieldComparatorDumper(t *testing.T) {
	ns := &storage{
		data: make(map[int64]record.Record),
	}
	ns.insert(&user{ID: 1, Name: "first", Age: 18})
	ns.insert(&user{ID: 2, Name: "second", Age: 19})
	ns.insert(&user{ID: 3, Name: "third", Age: 20})
	ns.insert(&user{ID: 4, Name: "fourth", Age: 21})
	ns.insert(&user{ID: 5, Name: "fifth", Age: 22})

	InRange := where.ComparatorType(20)

	newBuilder := WrapCreateQueryBuilderWithDumper(query.NewBuilder, func(w *strings.Builder, cmp where.FieldComparator) {
		if InRange == cmp.GetType() {
			w.WriteString(" IN RANGE (")
			w.WriteString(fmt.Sprintf("%v; ", cmp.ValueAt(0)))
			w.WriteString(fmt.Sprintf("%v", cmp.ValueAt(1)))
			w.WriteString(")")
		}
	})

	tests := []struct {
		name                 string
		query                query.Query
		expected             string
		expectedErrorMessage string
	}{
		{
			name: "IN RANGE",
			query: newBuilder(
				query.WhereInt64(id, InRange, 3, 10),
				query.Not(),
				query.WhereInt(age, where.LE, 21),
				query.Limit(2),
				query.Offset(1),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			expected:             "SELECT *, COUNT(*) WHERE id IN RANGE (3; 10) AND NOT age <= 21 ORDER BY &debug.byID{} ASC OFFSET 1 LIMIT 2",
			expectedErrorMessage: "execute query: not implemented ComparatorType: 20, field = id",
		},
	}

	qe := executor.CreateQueryExecutor(ns)

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			_, err := WrapQueryExecutor(qe, func(q string) {
				asserts.Equals(t, test.expected, q, "query")
			}).FetchAll(ctx, test.query)

			var errMsg string
			if nil != err {
				errMsg = err.Error()
			}
			asserts.Equals(t, test.expectedErrorMessage, errMsg, "error")
		})
	}
}
