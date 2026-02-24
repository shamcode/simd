//nolint:exhaustruct
package tests

import (
	"context"
	"regexp"
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
	// Arrange
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
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userStatus, where.EQ, StatusActive)).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE status = ACTIVE ORDER BY id DESC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userStatus, where.EQ, StatusActive)).
				Sort(sort.Desc(userID)).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{4, 1},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE status != DISABLED ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				Not().
				AddWhere(query.Where(userStatus, where.EQ, StatusDisabled)).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userScore, where.GE, 10)).
				AddWhere(query.Where(userScore, where.LT, 20)).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 2},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC LIMIT 1",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userScore, where.GE, 10)).
				AddWhere(query.Where(userScore, where.LT, 20)).
				Sort(sort.Asc(userID)).
				Limit(1).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC OFFSET 1 LIMIT 3",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userScore, where.GE, 10)).
				AddWhere(query.Where(userScore, where.LT, 20)).
				Sort(sort.Asc(userID)).
				Offset(1).
				Limit(3).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{2},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC OFFSET 2 LIMIT 3",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userScore, where.GE, 10)).
				AddWhere(query.Where(userScore, where.LT, 20)).
				Sort(sort.Asc(userID)).
				Offset(2).
				Limit(3).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name = 'Fourth' AND status == ACTIVE",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userName, where.EQ, "Fourth")).
				AddWhere(query.Where(userStatus, where.EQ, StatusActive)).
				Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name = 'Fourth' AND status == DISABLED",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userName, where.EQ, "Fourth")).
				AddWhere(query.Where(userStatus, where.EQ, StatusDisabled)).
				Query(),
			ExpectedCount: 0,
			ExpectedIDs:   []int64{},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name LIKE 'th' ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userName, where.Like, "t")).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR status == DISABLED ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userID, where.EQ, 1)).
				Or().
				AddWhere(query.Where(userStatus, where.EQ, StatusDisabled)).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR (NOT status == DISABLED) ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userID, where.EQ, 1)).
				Or().
				OpenBracket().
				Not().
				AddWhere(query.Where(userStatus, where.EQ, StatusDisabled)).
				CloseBracket().
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR (NOT status == ACTIVE OR NOT is_online = true) ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userID, where.EQ, 1)).
				Or().
				OpenBracket().
				Not().
				AddWhere(query.Where(userStatus, where.EQ, StatusActive)).
				Or().
				Not().
				AddWhere(query.WhereBool(userIsOnline, where.EQ, true)).
				CloseBracket().
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR (status == DISABLED OR is_online = false) ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userID, where.EQ, 1)).
				Or().
				OpenBracket().
				AddWhere(query.Where(userStatus, where.EQ, StatusDisabled)).
				Or().
				AddWhere(query.WhereBool(userIsOnline, where.EQ, false)).
				CloseBracket().
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE (status == DISABLED OR is_online = false) OR id = 1 ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				OpenBracket().
				AddWhere(query.Where(userStatus, where.EQ, StatusDisabled)).
				Or().
				AddWhere(query.WhereBool(userIsOnline, where.EQ, false)).
				CloseBracket().
				Or().
				AddWhere(query.Where(userID, where.EQ, 1)).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 4 OR (status == DISABLED OR is_online = false) ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userID, where.EQ, 4)).
				Or().
				OpenBracket().
				AddWhere(query.Where(userStatus, where.EQ, StatusDisabled)).
				Or().
				AddWhere(query.WhereBool(userIsOnline, where.EQ, false)).
				CloseBracket().
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{1, 2, 3, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 4 AND (status == DISABLED OR is_online = true) ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userID, where.EQ, 4)).
				OpenBracket().
				AddWhere(query.Where(userStatus, where.EQ, StatusDisabled)).
				Or().
				AddWhere(query.WhereBool(userIsOnline, where.EQ, true)).
				CloseBracket().
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE is_online = true ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.WhereBool(userIsOnline, where.EQ, true)).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id IN (4, 2) ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userID, where.InArray, 4, 2)).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{2, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name REGEXP [tT]) ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.WhereStringRegexp(userName, regexp.MustCompile("[tT]"))).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 3, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name IN (Second, Third) ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.Where(userName, where.InArray, "Second", "Third")).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE ( id = 1 ) AND id IN (1, 2, 3) ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				OpenBracket().
				AddWhere(query.Where(userID, where.EQ, 1)).
				CloseBracket().
				AddWhere(query.Where(userID, where.InArray, 1, 2, 3)).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{1},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE True ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{1, 2, 3, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE tags HAS confirmed ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.WhereSet(userTags, where.SetHas, TagConfirmed)).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE tags HAS free ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.WhereSet(userTags, where.SetHas, TagFree)).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE counter MAP_HAS_KEY UnreadMessages ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.WhereMap(userCounters, where.MapHasKey, CounterKeyUnreadMessages)).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 2},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE counter MAP_HAS_VALUE HasCounterValueEqual(2) ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.WhereMap(userCounters, where.MapHasValue, HasCounterValueEqual(2))).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{2},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE counter MAP_HAS_VALUE HasCounterValue(1) ORDER BY id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				AddWhere(query.WhereMap(userCounters, where.MapHasValue, HasCounterValueEqual(1))).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE True ORDER BY byOnline ASC id ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				Sort(sort.Asc(sort.ByScalar(byOnline{onlineToUp: true}))).
				Sort(sort.Asc(userID)).
				Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{4, 1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE True ORDER BY name ASC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				Sort(sort.Asc(userName)).
				Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{1, 4, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE True ORDER BY name DESC",
			Query: query.NewChainBuilder(query.NewBuilder[*User]()).
				Sort(sort.Desc(userName)).
				Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{3, 2, 4, 1},
		},
	}

	qe := executor.CreateQueryExecutor[*User](store)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			// Act
			cursor, count, err := qe.FetchAllAndTotal(ctx, testCase.Query)

			// Assert
			asserts.Success(t, err)

			ids := make([]int64, 0, cursor.Size())
			for item := range cursor.Seq(ctx) {
				ids = append(ids, item.ID)
			}

			asserts.Success(t, cursor.Err())
			asserts.Equals(t, testCase.ExpectedIDs, ids, "ids")
			asserts.Equals(t, testCase.ExpectedCount, count, "total count")
		})
	}
}
