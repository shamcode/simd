//nolint:exhaustruct
package tests

import (
	"context"
	"errors"
	"regexp"
	_sort "sort"
	"testing"

	asserts "github.com/shamcode/assert"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/indexes/btree"
	"github.com/shamcode/simd/indexes/hash"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
)

func Test_FetchAllAndTotal(t *testing.T) { //nolint:maintidx
	store := namespace.CreateNamespace[*User]()
	store.AddIndex(hash.NewComparableHashIndex(userID, true))
	store.AddIndex(hash.NewComparableHashIndex(userName, false))
	store.AddIndex(hash.NewComparableHashIndex(userStatus, false))
	store.AddIndex(hash.NewBoolHashIndex(userIsOnline, false))
	store.AddIndex(btree.NewComparableBTreeIndex(userScore, 16, false))
	asserts.Success(t, store.Insert(&User{
		ID:     1,
		Name:   "First",
		Status: StatusActive,
		Score:  10,
		Tags: map[Tag]struct{}{
			TagTester:    {},
			TagConfirmed: {},
		},
		Counters: Counters{
			CounterKeyUnreadMessages: 10,
			CounterKeyPendingTasks:   1,
		},
	}))
	asserts.Success(t, store.Insert(&User{
		ID:     2,
		Name:   "Second",
		Status: StatusDisabled,
		Score:  15,
		Tags: map[Tag]struct{}{
			TagConfirmed: {},
		},
		Counters: Counters{
			CounterKeyUnreadMessages: 2,
		},
	}))
	asserts.Success(t, store.Insert(&User{
		ID:     3,
		Name:   "Third",
		Status: StatusDisabled,
		Score:  20,
		Tags: map[Tag]struct{}{
			TagConfirmed: {},
			TagFree:      {},
		},
	}))
	asserts.Success(t, store.Insert(&User{
		ID:       4,
		Name:     "Fourth",
		Status:   StatusActive,
		Score:    25,
		IsOnline: true,
		Counters: Counters{
			CounterKeyPendingTasks: 1,
		},
	}))

	testCases := []struct {
		Name          string
		Query         query.Query[*User]
		ExpectedCount int
		ExpectedIDs   []int64
	}{
		{
			Name: "SELECT *, COUNT(*) WHERE status = ACTIVE ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.Where(userStatus, where.EQ, StatusActive),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE status = ACTIVE ORDER BY id DESC",
			Query: query.NewBuilder[*User](
				query.Where(userStatus, where.EQ, StatusActive),
				query.Sort(sort.Desc(userID)),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{4, 1},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE status != DISABLED ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.Not(),
				query.Where(userStatus, where.EQ, StatusDisabled),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.Where(userScore, where.GE, 10),
				query.Where(userScore, where.LT, 20),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 2},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC LIMIT 1",
			Query: query.NewBuilder[*User](
				query.Where(userScore, where.GE, 10),
				query.Where(userScore, where.LT, 20),
				query.Sort(sort.Asc(userID)),
				query.Limit(1),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC OFFSET 1 LIMIT 3",
			Query: query.NewBuilder[*User](
				query.Where(userScore, where.GE, 10),
				query.Where(userScore, where.LT, 20),
				query.Sort(sort.Asc(userID)),
				query.Offset(1),
				query.Limit(3),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{2},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC OFFSET 2 LIMIT 3",
			Query: query.NewBuilder[*User](
				query.Where(userScore, where.GE, 10),
				query.Where(userScore, where.LT, 20),
				query.Sort(sort.Asc(userID)),
				query.Offset(2),
				query.Limit(3),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name = 'Fourth' AND status == ACTIVE",
			Query: query.NewBuilder[*User](
				query.Where(userName, where.EQ, "Fourth"),
				query.Where(userStatus, where.EQ, StatusActive),
			).Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name = 'Fourth' AND status == DISABLED",
			Query: query.NewBuilder[*User](
				query.Where(userName, where.EQ, "Fourth"),
				query.Where(userStatus, where.EQ, StatusDisabled),
			).Query(),
			ExpectedCount: 0,
			ExpectedIDs:   []int64{},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name LIKE 'th' ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.Where(userName, where.Like, "t"),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR status == DISABLED ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.Where(userID, where.EQ, 1),
				query.Or(),
				query.Where(userStatus, where.EQ, StatusDisabled),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR (NOT status == DISABLED) ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.Where(userID, where.EQ, 1),
				query.Or(),
				query.OpenBracket(),
				query.Not(),
				query.Where(userStatus, where.EQ, StatusDisabled),
				query.CloseBracket(),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR (NOT status == ACTIVE OR NOT is_online = true) ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.Where(userID, where.EQ, 1),
				query.Or(),
				query.OpenBracket(),
				query.Not(),
				query.Where(userStatus, where.EQ, StatusActive),
				query.Or(),
				query.Not(),
				query.WhereBool(userIsOnline, where.EQ, true),
				query.CloseBracket(),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR (status == DISABLED OR is_online = false) ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.Where(userID, where.EQ, 1),
				query.Or(),
				query.OpenBracket(),
				query.Where(userStatus, where.EQ, StatusDisabled),
				query.Or(),
				query.WhereBool(userIsOnline, where.EQ, false),
				query.CloseBracket(),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE (status == DISABLED OR is_online = false) OR id = 1 ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.OpenBracket(),
				query.Where(userStatus, where.EQ, StatusDisabled),
				query.Or(),
				query.WhereBool(userIsOnline, where.EQ, false),
				query.CloseBracket(),
				query.Or(),
				query.Where(userID, where.EQ, 1),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 4 OR (status == DISABLED OR is_online = false) ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.Where(userID, where.EQ, 4),
				query.Or(),
				query.OpenBracket(),
				query.Where(userStatus, where.EQ, StatusDisabled),
				query.Or(),
				query.WhereBool(userIsOnline, where.EQ, false),
				query.CloseBracket(),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{1, 2, 3, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 4 AND (status == DISABLED OR is_online = true) ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.Where(userID, where.EQ, 4),
				query.OpenBracket(),
				query.Where(userStatus, where.EQ, StatusDisabled),
				query.Or(),
				query.WhereBool(userIsOnline, where.EQ, true),
				query.CloseBracket(),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE is_online = true ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.WhereBool(userIsOnline, where.EQ, true),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id IN (4, 2) ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.Where(userID, where.InArray, 4, 2),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{2, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name REGEXP [tT]) ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.WhereStringRegexp(userName, regexp.MustCompile("[tT]")),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 3, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name IN (Second, Third) ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.Where(userName, where.InArray, "Second", "Third"),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE ( id = 1 ) AND id IN (1, 2, 3) ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.OpenBracket(),
				query.Where(userID, where.EQ, 1),
				query.CloseBracket(),
				query.Where(userID, where.InArray, 1, 2, 3),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{1},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE True ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{1, 2, 3, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE tags HAS confirmed ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.WhereSet(userTags, where.SetHas, TagConfirmed),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE tags HAS free ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.WhereSet(userTags, where.SetHas, TagFree),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE counter MAP_HAS_KEY UnreadMessages ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.WhereMap(userCounters, where.MapHasKey, CounterKeyUnreadMessages),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 2},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE counter MAP_HAS_VALUE HasCounterValueEqual(2) ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.WhereMap(userCounters, where.MapHasValue, HasCounterValueEqual(2)),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{2},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE counter MAP_HAS_VALUE HasCounterValue(1) ORDER BY id ASC",
			Query: query.NewBuilder[*User](
				query.WhereMap(userCounters, where.MapHasValue, HasCounterValueEqual(1)),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE True ORDER BY byOnline ASC id ASC",
			Query: query.NewBuilder[*User](
				query.Sort(sort.Asc(sort.ByScalar(byOnline{onlineToUp: true}))),
				query.Sort(sort.Asc(userID)),
			).Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{4, 1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE True ORDER BY name ASC",
			Query: query.NewBuilder[*User](
				query.Sort(sort.Asc(userName)),
			).Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{1, 4, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE True ORDER BY name DESC",
			Query: query.NewBuilder[*User](
				query.Sort(sort.Desc(userName)),
			).Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{3, 2, 4, 1},
		},
	}

	qe := executor.CreateQueryExecutor[*User](store)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			cursor, count, err := qe.FetchAllAndTotal(ctx, testCase.Query)
			asserts.Success(t, err)

			ids := make([]int64, 0, cursor.Size())
			for cursor.Next(ctx) {
				ids = append(ids, cursor.Item().ID)
			}

			asserts.Success(t, cursor.Err())
			asserts.Equals(t, testCase.ExpectedIDs, ids, "ids")
			asserts.Equals(t, testCase.ExpectedCount, count, "total count")
		})
	}
}

func Test_Context(t *testing.T) {
	store := namespace.CreateNamespace[*User]()
	store.AddIndex(hash.NewComparableHashIndex(userID, true))
	asserts.Success(t, store.Insert(&User{
		ID:     1,
		Name:   "First",
		Status: StatusActive,
		Score:  10,
	}))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := executor.CreateQueryExecutor[*User](store).FetchTotal(ctx, query.NewBuilder[*User]().Query())

	asserts.Equals(t, "context canceled", err.Error(), "check error")
	asserts.Equals(t, true, errors.Is(err, context.Canceled), "error is context.Canceled")
}

func Test_CallbackOnIteration(t *testing.T) {
	store := namespace.CreateNamespace[*User]()
	store.AddIndex(hash.NewComparableHashIndex(userID, true))
	asserts.Success(t, store.Insert(&User{
		ID:     1,
		Status: StatusActive,
	}))
	asserts.Success(t, store.Insert(&User{
		ID:     2,
		Status: StatusDisabled,
	}))
	asserts.Success(t, store.Insert(&User{
		ID:     3,
		Status: StatusActive,
	}))

	var (
		idsFromCallback []int
		idsFromCursor   []int64
	)

	cur, err := executor.CreateQueryExecutor[*User](store).FetchAll(
		context.Background(),
		query.NewBuilder[*User](
			query.Where(userStatus, where.EQ, StatusActive),
			query.Limit(1),
			query.Sort(sort.Asc(userID)),
			query.OnIteration(func(item *User) {
				idsFromCallback = append(idsFromCallback, int(item.GetID()))
			}),
		).Query(),
	)
	asserts.Success(t, err)

	for cur.Next(context.Background()) {
		idsFromCursor = append(idsFromCursor, cur.Item().GetID())
	}

	_sort.Ints(idsFromCallback)
	asserts.Success(t, cur.Err())
	asserts.Equals(t, []int{1, 3}, idsFromCallback, "ids from callback")
	asserts.Equals(t, []int64{1}, idsFromCursor, "ids from cursor")
}
