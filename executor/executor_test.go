package executor

import (
	"context"
	"testing"

	asserts "github.com/shamcode/assert"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
)

type user struct {
	ID   int64
	Name string
	Age  int
}

func (u *user) GetID() int64 {
	return u.ID
}

var userFields = record.NewFields()

var id = record.NewIDGetter[*user]()

var name = record.StringGetter[*user]{
	Field: userFields.New("name"),
	Get: func(item *user) string {
		return item.Name
	},
}

var age = record.ComparableGetter[*user, int]{
	Field: userFields.New("age"),
	Get: func(item *user) int {
		return item.Age
	},
}

type storage struct {
	data map[int64]*user
}

func (s *storage) Get(id int64) *user {
	return s.data[id]
}

func (s *storage) Insert(item *user) error {
	s.data[item.GetID()] = item
	return nil
}

func (s *storage) Delete(id int64) error {
	delete(s.data, id)
	return nil
}

func (s *storage) Upsert(item *user) error {
	s.data[item.GetID()] = item
	return nil
}

func (s *storage) PreselectForExecutor(_ where.Conditions[*user]) ([]*user, error) {
	items := make([]*user, 0, len(s.data))
	for _, item := range s.data {
		items = append(items, item)
	}
	return items, nil
}

func TestQueryExecutor(t *testing.T) { //nolint:maintidx
	ns := &storage{
		data: make(map[int64]*user),
	}
	asserts.Success(t, ns.Insert(&user{ID: 1, Name: "first", Age: 18}))
	asserts.Success(t, ns.Insert(&user{ID: 2, Name: "second", Age: 19}))
	asserts.Success(t, ns.Insert(&user{ID: 3, Name: "third", Age: 20}))
	asserts.Success(t, ns.Insert(&user{ID: 4, Name: "fourth", Age: 21}))
	asserts.Success(t, ns.Insert(&user{ID: 5, Name: "fifth", Age: 22}))

	tests := []struct {
		name     string
		query    query.Query[*user]
		expected []int64
	}{
		{
			name:     "order by id asc",
			query:    query.NewBuilder[*user](query.Sort(sort.Asc[*user](id))).Query(),
			expected: []int64{1, 2, 3, 4, 5},
		},
		{
			name:     "order by id desc",
			query:    query.NewBuilder[*user](query.Sort(sort.Desc[*user](id))).Query(),
			expected: []int64{5, 4, 3, 2, 1},
		},
		{
			name:     "where id = int64(3)",
			query:    query.NewBuilder[*user](query.Where(id, where.EQ, 3)).Query(),
			expected: []int64{3},
		},
		{
			name: "where id = int64(3) and age == int(20)",
			query: query.NewBuilder[*user](
				query.Where(id, where.EQ, 3),
				query.Where(age, where.EQ, 20),
			).Query(),
			expected: []int64{3},
		},
		{
			name: "where id = int64(3) and age == int(18)",
			query: query.NewBuilder[*user](
				query.Where(id, where.EQ, 3),
				query.Where(age, where.EQ, 18),
			).Query(),
			expected: []int64{},
		},
		{
			name: "where age > 18 and age < 22",
			query: query.NewBuilder[*user](
				query.Where(age, where.GT, 18),
				query.Where(age, where.LT, 22),
				query.Sort(sort.Asc[*user](id)),
			).Query(),
			expected: []int64{2, 3, 4},
		},
		{
			name: "where age >= 18 and age <= 22",
			query: query.NewBuilder[*user](
				query.Where(age, where.GE, 18),
				query.Where(age, where.LE, 22),
				query.Sort(sort.Asc[*user](id)),
			).Query(),
			expected: []int64{1, 2, 3, 4, 5},
		},
		{
			name: "where id = 2 or id = 5",
			query: query.NewBuilder[*user](
				query.Where(id, where.EQ, 2),
				query.Or(),
				query.Where(id, where.EQ, 5),
				query.Sort(sort.Asc[*user](id)),
			).Query(),
			expected: []int64{2, 5},
		},
		{
			name: "where id = 2 or age > 20",
			query: query.NewBuilder[*user](
				query.Where(id, where.EQ, 2),
				query.Or(),
				query.Where(age, where.GT, 20),
				query.Sort(sort.Asc[*user](id)),
			).Query(),
			expected: []int64{2, 4, 5},
		},
		{
			name: "where id = 1 or ( age > 20 and age < 22)",
			query: query.NewBuilder[*user](
				query.Where(id, where.EQ, 1),
				query.Or(),
				query.OpenBracket(),
				query.Where(age, where.GT, 20),
				query.Where(age, where.LT, 22),
				query.CloseBracket(),
				query.Sort(sort.Asc[*user](id)),
			).Query(),
			expected: []int64{1, 4},
		},
		{
			name: "where ( age > 20 and age < 22) or id = 1",
			query: query.NewBuilder[*user](
				query.OpenBracket(),
				query.Where(age, where.GT, 20),
				query.Where(age, where.LT, 22),
				query.CloseBracket(),
				query.Or(),
				query.Where(id, where.EQ, 1),
				query.Sort(sort.Asc[*user](id)),
			).Query(),
			expected: []int64{1, 4},
		},
		{
			name: "where age > 20 and age < 22 or id = 1",
			query: query.NewBuilder[*user](
				query.Where(age, where.GT, 20),
				query.Where(age, where.LT, 22),
				query.Or(),
				query.Where(id, where.EQ, 1),
				query.Sort(sort.Asc[*user](id)),
			).Query(),
			expected: []int64{1, 4},
		},
		{
			name: "where (age > 20 and age < 22) or id = 1",
			query: query.NewBuilder[*user](
				query.OpenBracket(),
				query.Where(age, where.GT, 20),
				query.Where(age, where.LT, 22),
				query.CloseBracket(),
				query.Or(),
				query.Where(id, where.EQ, 1),
				query.Sort(sort.Asc[*user](id)),
			).Query(),
			expected: []int64{1, 4},
		},
		{
			name: "where age > 20 and age < 22 and (id = 1 or id = 2)",
			query: query.NewBuilder[*user](
				query.Where(age, where.GT, 20),
				query.Where(age, where.LT, 22),
				query.OpenBracket(),
				query.Where(id, where.EQ, 1),
				query.Or(),
				query.Where(id, where.EQ, 2),
				query.CloseBracket(),
			).Query(),
			expected: []int64{},
		},
		{
			name: "where age > 20 and age < 22 and (id = 1 or id = 2 or id = 4)",
			query: query.NewBuilder[*user](
				query.Where(age, where.GT, 20),
				query.Where(age, where.LT, 22),
				query.OpenBracket(),
				query.Where(id, where.EQ, 1),
				query.Or(),
				query.Where(id, where.EQ, 2),
				query.Or(),
				query.Where(id, where.EQ, 4),
				query.CloseBracket(),
			).Query(),
			expected: []int64{4},
		},
		{
			name: "where (age > 20 and age < 22) and id = 4",
			query: query.NewBuilder[*user](
				query.OpenBracket(),
				query.Where(age, where.GT, 20),
				query.Where(age, where.LT, 22),
				query.CloseBracket(),
				query.Where(id, where.EQ, 4),
			).Query(),
			expected: []int64{4},
		},
		{
			name: "where age in {20, 21, 22} and id > 3",
			query: query.NewBuilder[*user](
				query.Where(age, where.InArray, 20, 21, 22),
				query.Where(id, where.GT, 3),
				query.Sort(sort.Asc[*user](id)),
			).Query(),
			expected: []int64{4, 5},
		},
		{
			name: "where name like \"th\"",
			query: query.NewBuilder[*user](
				query.WhereString(name, where.Like, "th"),
				query.Sort(sort.Asc[*user](id)),
			).Query(),
			expected: []int64{3, 4, 5},
		},
		{
			name: "where name like \"th\" or name like \"first\"",
			query: query.NewBuilder[*user](
				query.WhereString(name, where.Like, "th"),
				query.Or(),
				query.WhereString(name, where.Like, "first"),
				query.Sort(sort.Asc[*user](id)),
			).Query(),
			expected: []int64{1, 3, 4, 5},
		},
		{
			name: "where ((id = 1) or (id = 2))",
			query: query.NewBuilder[*user](
				query.OpenBracket(),
				query.OpenBracket(),
				query.Where(id, where.EQ, 1),
				query.CloseBracket(),
				query.Or(),
				query.OpenBracket(),
				query.Where(id, where.EQ, 2),
				query.CloseBracket(),
				query.CloseBracket(),
				query.Sort(sort.Asc[*user](id)),
			).Query(),
			expected: []int64{1, 2},
		},
		{
			name: "where (((id = 1) or (id = 2)) or id = 3) or id = 4",
			query: query.NewBuilder[*user](
				query.OpenBracket(),
				query.OpenBracket(),
				query.OpenBracket(),
				query.Where(id, where.EQ, 1),
				query.CloseBracket(),
				query.Or(),
				query.OpenBracket(),
				query.Where(id, where.EQ, 2),
				query.CloseBracket(),
				query.CloseBracket(),
				query.Or(),
				query.Where(id, where.EQ, 3),
				query.CloseBracket(),
				query.Or(),
				query.Where(id, where.EQ, 4),
				query.Sort(sort.Asc[*user](id)),
			).Query(),
			expected: []int64{1, 2, 3, 4},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			iter, err := CreateQueryExecutor[*user](ns).FetchAll(ctx, test.query)
			asserts.Success(t, err)
			res := make([]int64, 0, iter.Size())
			for iter.Next(ctx) {
				res = append(res, iter.Item().ID)
			}
			asserts.Equals(t, test.expected, res, "ids")
		})
	}
}
