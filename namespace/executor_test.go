package namespace

import (
	"context"
	"fmt"
	"github.com/shamcode/simd/asserts"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
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

var id = &record.Int64Getter{
	Field: "id",
	Get: func(item interface{}) int64 {
		return item.(*user).ID
	},
}

var name = &record.StringGetter{
	Field: "name",
	Get: func(item interface{}) string {
		return item.(*user).Name
	},
}

var age = &record.IntGetter{
	Field: "age",
	Get: func(item interface{}) int {
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

func (s *storage) SelectForExecutor(conditions where.Conditions) ([]record.Record, error) {
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

func TestQueryExecutor(t *testing.T) {
	ns := &storage{
		data: make(map[int64]record.Record),
	}
	ns.insert(&user{ID: 1, Name: "first", Age: 18})
	ns.insert(&user{ID: 2, Name: "second", Age: 19})
	ns.insert(&user{ID: 3, Name: "third", Age: 20})
	ns.insert(&user{ID: 4, Name: "fourth", Age: 21})
	ns.insert(&user{ID: 5, Name: "fifth", Age: 22})

	tests := []struct {
		name     string
		query    query.Query
		expected []int64
	}{
		{
			name:     "order by id asc",
			query:    query.NewBuilder().Sort(sort.ByInt64IndexAsc(&byID{})).Query(),
			expected: []int64{1, 2, 3, 4, 5},
		},
		{
			name:     "order by id desc",
			query:    query.NewBuilder().Sort(sort.ByInt64IndexDesc(&byID{})).Query(),
			expected: []int64{5, 4, 3, 2, 1},
		},
		{
			name:     "where id = int64(3)",
			query:    query.NewBuilder().WhereInt64(id, where.EQ, 3).Query(),
			expected: []int64{3},
		},
		{
			name:     "where id = int64(3) and age == int(20)",
			query:    query.NewBuilder().WhereInt64(id, where.EQ, 3).WhereInt(age, where.EQ, 20).Query(),
			expected: []int64{3},
		},
		{
			name:     "where id = int64(3) and age == int(18)",
			query:    query.NewBuilder().WhereInt64(id, where.EQ, 3).WhereInt(age, where.EQ, 18).Query(),
			expected: []int64{},
		},
		{
			name:     "where age > 18 and age < 22",
			query:    query.NewBuilder().WhereInt(age, where.GT, 18).WhereInt(age, where.LT, 22).Sort(sort.ByInt64IndexAsc(&byID{})).Query(),
			expected: []int64{2, 3, 4},
		},
		{
			name:     "where age >= 18 and age <= 22",
			query:    query.NewBuilder().WhereInt(age, where.GE, 18).WhereInt(age, where.LE, 22).Sort(sort.ByInt64IndexAsc(&byID{})).Query(),
			expected: []int64{1, 2, 3, 4, 5},
		},
		{
			name:     "where id = 2 or id = 5",
			query:    query.NewBuilder().WhereInt64(id, where.EQ, 2).Or().WhereInt64(id, where.EQ, 5).Sort(sort.ByInt64IndexAsc(&byID{})).Query(),
			expected: []int64{2, 5},
		},
		{
			name:     "where id = 2 or age > 20",
			query:    query.NewBuilder().WhereInt64(id, where.EQ, 2).Or().WhereInt(age, where.GT, 20).Sort(sort.ByInt64IndexAsc(&byID{})).Query(),
			expected: []int64{2, 4, 5},
		},
		{
			name: "where id = 1 or ( age > 20 and age < 22)",
			query: query.NewBuilder().
				WhereInt64(id, where.EQ, 1).
				Or().
				OpenBracket().
				WhereInt(age, where.GT, 20).
				WhereInt(age, where.LT, 22).
				CloseBracket().
				Sort(sort.ByInt64IndexAsc(&byID{})).
				Query(),
			expected: []int64{1, 4},
		},
		{
			name: "where ( age > 20 and age < 22) or id = 1",
			query: query.NewBuilder().
				OpenBracket().
				WhereInt(age, where.GT, 20).
				WhereInt(age, where.LT, 22).
				CloseBracket().
				Or().
				WhereInt64(id, where.EQ, 1).
				Sort(sort.ByInt64IndexAsc(&byID{})).
				Query(),
			expected: []int64{1, 4},
		},
		{
			name: "where age > 20 and age < 22 or id = 1",
			query: query.NewBuilder().
				WhereInt(age, where.GT, 20).
				WhereInt(age, where.LT, 22).
				Or().
				WhereInt64(id, where.EQ, 1).
				Sort(sort.ByInt64IndexAsc(&byID{})).
				Query(),
			expected: []int64{1, 4},
		},
		{
			name: "where (age > 20 and age < 22) or id = 1",
			query: query.NewBuilder().
				OpenBracket().
				WhereInt(age, where.GT, 20).
				WhereInt(age, where.LT, 22).
				CloseBracket().
				Or().
				WhereInt64(id, where.EQ, 1).
				Sort(sort.ByInt64IndexAsc(&byID{})).
				Query(),
			expected: []int64{1, 4},
		},
		{
			name: "where age > 20 and age < 22 and (id = 1 or id = 2)",
			query: query.NewBuilder().
				WhereInt(age, where.GT, 20).
				WhereInt(age, where.LT, 22).
				OpenBracket().
				WhereInt64(id, where.EQ, 1).
				Or().
				WhereInt64(id, where.EQ, 2).
				CloseBracket().
				Query(),
			expected: []int64{},
		},
		{
			name: "where age > 20 and age < 22 and (id = 1 or id = 2 or id = 4)",
			query: query.NewBuilder().
				WhereInt(age, where.GT, 20).
				WhereInt(age, where.LT, 22).
				OpenBracket().
				WhereInt64(id, where.EQ, 1).
				Or().
				WhereInt64(id, where.EQ, 2).
				Or().
				WhereInt64(id, where.EQ, 4).
				CloseBracket().
				Query(),
			expected: []int64{4},
		},
		{
			name: "where (age > 20 and age < 22) and id = 4",
			query: query.NewBuilder().
				OpenBracket().
				WhereInt(age, where.GT, 20).
				WhereInt(age, where.LT, 22).
				CloseBracket().
				WhereInt64(id, where.EQ, 4).
				Query(),
			expected: []int64{4},
		},
		{
			name: "where age in {20, 21, 22} and id > 3",
			query: query.NewBuilder().
				WhereInt(age, where.InArray, 20, 21, 22).
				WhereInt64(id, where.GT, 3).
				Sort(sort.ByInt64IndexAsc(&byID{})).
				Query(),
			expected: []int64{4, 5},
		},
		{
			name: "where name like \"th\"",
			query: query.NewBuilder().WhereString(name, where.Like, "th").Sort(sort.ByInt64IndexAsc(&byID{})).
				Query(),
			expected: []int64{3, 4, 5},
		},
		{
			name: "where name like \"th\" or name like \"first\"",
			query: query.NewBuilder().
				WhereString(name, where.Like, "th").
				Or().
				WhereString(name, where.Like, "first").
				Sort(sort.ByInt64IndexAsc(&byID{})).
				Query(),
			expected: []int64{1, 3, 4, 5},
		},
		{
			name: "where ((id = 1) or (id = 2))",
			query: query.NewBuilder().
				OpenBracket().
				OpenBracket().
				WhereInt64(id, where.EQ, 1).
				CloseBracket().
				Or().
				OpenBracket().
				WhereInt64(id, where.EQ, 2).
				CloseBracket().
				CloseBracket().
				Sort(sort.ByInt64IndexAsc(&byID{})).
				Query(),
			expected: []int64{1, 2},
		},
		{
			name: "where (((id = 1) or (id = 2)) or id = 3) or id = 4",
			query: query.NewBuilder().
				OpenBracket().
				OpenBracket().
				OpenBracket().
				WhereInt64(id, where.EQ, 1).
				CloseBracket().
				Or().
				OpenBracket().
				WhereInt64(id, where.EQ, 2).
				CloseBracket().
				CloseBracket().
				Or().
				WhereInt64(id, where.EQ, 3).
				CloseBracket().
				Or().
				WhereInt64(id, where.EQ, 4).
				Sort(sort.ByInt64IndexAsc(&byID{})).
				Query(),
			expected: []int64{1, 2, 3, 4},
		},
	}

	for _, test := range tests {
		ctx := context.Background()
		iter, err := CreateQueryExecutor(ns).FetchAll(ctx, test.query)
		res := make([]int64, 0, iter.Size())
		for iter.Next(ctx) {
			res = append(res, iter.Item().(*user).ID)
		}
		asserts.Equals(t, nil, err, fmt.Sprintf("err is nil for test \"%s\"", test.name))
		asserts.Equals(t, test.expected, res, fmt.Sprintf("res for test \"%s\"", test.name))
	}
}
