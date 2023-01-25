package simd

import (
	"context"
	"errors"
	"github.com/shamcode/simd/asserts"
	"github.com/shamcode/simd/executor"
	"github.com/shamcode/simd/indexes/hash"
	"github.com/shamcode/simd/namespace"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"regexp"
	_sort "sort"
	"testing"
)

type StatusEnum uint8

const (
	StatusActive StatusEnum = iota + 1
	StatusDisabled
)

func (s StatusEnum) Value() uint8 { return uint8(s) }
func (s StatusEnum) String() string {
	switch s {
	case StatusActive:
		return "ACTIVE"
	case StatusDisabled:
		return "DISABLED"
	default:
		return ""
	}
}

type Tag uint16

const (
	TagTester Tag = iota + 1
	TagConfirmed
	TagFree
)

func (t Tag) String() string {
	switch t {
	case TagTester:
		return "tester"
	case TagConfirmed:
		return "confirmed"
	case TagFree:
		return "free"
	default:
		return ""
	}
}

type Tags map[Tag]struct{}

func (t Tags) Has(item interface{}) bool {
	value, ok := item.(Tag)
	if !ok {
		return false
	}
	_, ok = t[value]
	return ok
}

type CounterKey uint16

const (
	CounterKeyUnreadMessages CounterKey = iota + 1
	CounterKeyPendingTasks
)

type Counters map[CounterKey]uint32

func (c Counters) HasKey(key interface{}) bool {
	counterKey, ok := key.(CounterKey)
	if !ok {
		return false
	}
	_, ok = c[counterKey]
	return ok
}
func (c Counters) HasValue(check record.MapValueComparator) (bool, error) {
	for _, item := range c {
		res, err := check.Compare(item)
		if nil != err {
			return false, err
		}
		if res {
			return true, nil
		}
	}
	return false, nil
}

type HasCounterValueEqual uint32

func (c HasCounterValueEqual) Compare(item interface{}) (bool, error) {
	return item.(uint32) == uint32(c), nil
}

type User struct {
	ID       int64
	Name     string
	Status   StatusEnum
	Score    int
	IsOnline bool
	Tags     Tags
	Counters Counters
}

func (u *User) GetID() int64   { return u.ID }
func (u *User) ComputeFields() {}

var userID = &record.Int64Getter{
	Field: "id",
	Get:   func(item record.Record) int64 { return item.(*User).ID },
}

var userName = &record.StringGetter{
	Field: "name",
	Get:   func(item record.Record) string { return item.(*User).Name },
}

var userStatus = &record.Enum8Getter{
	Field: "status",
	Get:   func(item record.Record) record.Enum8 { return item.(*User).Status },
}

var userScore = &record.IntGetter{
	Field: "score",
	Get:   func(item record.Record) int { return item.(*User).Score },
}

var userIsOnline = &record.BoolGetter{
	Field: "is_online",
	Get:   func(item record.Record) bool { return item.(*User).IsOnline },
}

var userTags = &record.SetGetter{
	Field: "tags",
	Get:   func(item record.Record) record.Set { return item.(*User).Tags },
}

var userCounters = &record.MapGetter{
	Field: "counters",
	Get:   func(item record.Record) record.Map { return item.(*User).Counters },
}

type byID struct{}

func (sorting *byID) CalcIndex(item record.Record) int64 { return item.(*User).ID }

type byOnline struct {
	onlineToUp bool
}

func (sorting *byOnline) CalcIndex(item record.Record) int64 {
	if sorting.onlineToUp == item.(*User).IsOnline {
		return 0
	}
	return 1
}

func Test_FetchAllAndTotal(t *testing.T) {
	store := namespace.CreateNamespace()
	store.AddIndex(hash.NewInt64HashIndex(userID, true))
	store.AddIndex(hash.NewStringHashIndex(userName, false))
	store.AddIndex(hash.NewEnum8HashIndex(userStatus, false))
	store.AddIndex(hash.NewBoolHashIndex(userIsOnline, false))
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
		Query         query.Query
		ExpectedCount int
		ExpectedIDs   []int64
	}{
		{
			Name: "SELECT *, COUNT(*) WHERE status = ACTIVE ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereEnum8(userStatus, where.EQ, StatusActive),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE status = ACTIVE ORDER BY id DESC",
			Query: query.NewBuilder(
				query.WhereEnum8(userStatus, where.EQ, StatusActive),
				query.Sort(sort.ByInt64IndexDesc(&byID{})),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{4, 1},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE status != DISABLED ORDER BY id ASC",
			Query: query.NewBuilder(
				query.Not(),
				query.WhereEnum8(userStatus, where.EQ, StatusDisabled),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereInt(userScore, where.GE, 10),
				query.WhereInt(userScore, where.LT, 20),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 2},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC LIMIT 1",
			Query: query.NewBuilder(
				query.WhereInt(userScore, where.GE, 10),
				query.WhereInt(userScore, where.LT, 20),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
				query.Limit(1),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC OFFSET 1 LIMIT 3",
			Query: query.NewBuilder(
				query.WhereInt(userScore, where.GE, 10),
				query.WhereInt(userScore, where.LT, 20),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
				query.Offset(1),
				query.Limit(3),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{2},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC OFFSET 2 LIMIT 3",
			Query: query.NewBuilder(
				query.WhereInt(userScore, where.GE, 10),
				query.WhereInt(userScore, where.LT, 20),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
				query.Offset(2),
				query.Limit(3),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name = 'Fourth' AND status == ACTIVE",
			Query: query.NewBuilder(
				query.WhereString(userName, where.EQ, "Fourth"),
				query.WhereEnum8(userStatus, where.EQ, StatusActive),
			).Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name = 'Fourth' AND status == DISABLED",
			Query: query.NewBuilder(
				query.WhereString(userName, where.EQ, "Fourth"),
				query.WhereEnum8(userStatus, where.EQ, StatusDisabled),
			).Query(),
			ExpectedCount: 0,
			ExpectedIDs:   []int64{},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name LIKE 'th' ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereString(userName, where.Like, "t"),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR status == DISABLED ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.EQ, 1),
				query.Or(),
				query.WhereEnum8(userStatus, where.EQ, StatusDisabled),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR (NOT status == DISABLED) ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.EQ, 1),
				query.Or(),
				query.OpenBracket(),
				query.Not(),
				query.WhereEnum8(userStatus, where.EQ, StatusDisabled),
				query.CloseBracket(),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR (NOT status == ACTIVE OR NOT is_online = true) ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.EQ, 1),
				query.Or(),
				query.OpenBracket(),
				query.Not(),
				query.WhereEnum8(userStatus, where.EQ, StatusActive),
				query.Or(),
				query.Not(),
				query.WhereBool(userIsOnline, where.EQ, true),
				query.CloseBracket(),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR (status == DISABLED OR is_online = false) ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.EQ, 1),
				query.Or(),
				query.OpenBracket(),
				query.WhereEnum8(userStatus, where.EQ, StatusDisabled),
				query.Or(),
				query.WhereBool(userIsOnline, where.EQ, false),
				query.CloseBracket(),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE (status == DISABLED OR is_online = false) OR id = 1 ORDER BY id ASC",
			Query: query.NewBuilder(
				query.OpenBracket(),
				query.WhereEnum8(userStatus, where.EQ, StatusDisabled),
				query.Or(),
				query.WhereBool(userIsOnline, where.EQ, false),
				query.CloseBracket(),
				query.Or(),
				query.WhereInt64(userID, where.EQ, 1),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 4 OR (status == DISABLED OR is_online = false) ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.EQ, 4),
				query.Or(),
				query.OpenBracket(),
				query.WhereEnum8(userStatus, where.EQ, StatusDisabled),
				query.Or(),
				query.WhereBool(userIsOnline, where.EQ, false),
				query.CloseBracket(),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{1, 2, 3, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 4 AND (status == DISABLED OR is_online = true) ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.EQ, 4),
				query.OpenBracket(),
				query.WhereEnum8(userStatus, where.EQ, StatusDisabled),
				query.Or(),
				query.WhereBool(userIsOnline, where.EQ, true),
				query.CloseBracket(),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE is_online = true ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereBool(userIsOnline, where.EQ, true),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id IN (4, 2) ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereInt64(userID, where.InArray, 4, 2),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{2, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name REGEXP [tT]) ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereStringRegexp(userName, regexp.MustCompile("[tT]")),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 3, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name IN (Second, Third) ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereString(userName, where.InArray, "Second", "Third"),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE ( id = 1 ) AND id IN (1, 2, 3) ORDER BY id ASC",
			Query: query.NewBuilder(
				query.OpenBracket(),
				query.WhereInt64(userID, where.EQ, 1),
				query.CloseBracket(),
				query.WhereInt64(userID, where.InArray, 1, 2, 3),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{1},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE True ORDER BY id ASC",
			Query: query.NewBuilder(
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{1, 2, 3, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE tags HAS confirmed ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereSet(userTags, where.SetHas, TagConfirmed),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE tags HAS free ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereSet(userTags, where.SetHas, TagFree),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE counter MAP_HAS_KEY UnreadMessages ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereMap(userCounters, where.MapHasKey, CounterKeyUnreadMessages),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 2},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE counter MAP_HAS_VALUE HasCounterValueEqual(2) ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereMap(userCounters, where.MapHasValue, HasCounterValueEqual(2)),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{2},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE counter MAP_HAS_VALUE HasCounterValue(1) ORDER BY id ASC",
			Query: query.NewBuilder(
				query.WhereMap(userCounters, where.MapHasValue, HasCounterValueEqual(1)),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE True ORDER BY byOnline ASC id ASC",
			Query: query.NewBuilder(
				query.Sort(sort.ByInt64IndexAsc(&byOnline{onlineToUp: true})),
				query.Sort(sort.ByInt64IndexAsc(&byID{})),
			).Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{4, 1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE True ORDER BY name ASC",
			Query: query.NewBuilder(
				query.Sort(sort.ByStringAsc(userName)),
			).Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{1, 4, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE True ORDER BY name DESC",
			Query: query.NewBuilder(
				query.Sort(sort.ByStringDesc(userName)),
			).Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{3, 2, 4, 1},
		},
	}

	qe := executor.CreateQueryExecutor(store)

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			cursor, count, err := qe.FetchAllAndTotal(ctx, testCase.Query)
			asserts.Success(t, err)
			ids := make([]int64, 0, cursor.Size())
			for cursor.Next(ctx) {
				ids = append(ids, cursor.Item().(*User).ID)
			}
			asserts.Success(t, cursor.Err())
			asserts.Equals(t, testCase.ExpectedIDs, ids, "ids")
			asserts.Equals(t, testCase.ExpectedCount, count, "total count")
		})
	}
}

func Test_Context(t *testing.T) {
	store := namespace.CreateNamespace()
	store.AddIndex(hash.NewInt64HashIndex(userID, true))
	asserts.Success(t, store.Insert(&User{
		ID:     1,
		Name:   "First",
		Status: StatusActive,
		Score:  10,
	}))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := executor.CreateQueryExecutor(store).FetchTotal(ctx, query.NewBuilder().Query())

	asserts.Equals(t, "context canceled", err.Error(), "check error")
	asserts.Equals(t, true, errors.Is(context.Canceled, err), "error is context.Canceled")
}

func Test_CallbackOnIteration(t *testing.T) {
	store := namespace.CreateNamespace()
	store.AddIndex(hash.NewInt64HashIndex(userID, true))
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

	var idsFromCallback []int
	var idsFromCursor []int64
	cur, err := executor.CreateQueryExecutor(store).FetchAll(
		context.Background(),
		query.NewBuilder(
			query.WhereEnum8(userStatus, where.EQ, StatusActive),
			query.Limit(1),
			query.Sort(sort.ByInt64IndexAsc(&byID{})),
			query.OnIteration(func(item record.Record) {
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

func Test_InsertAlreadyExisted(t *testing.T) {
	store := namespace.CreateNamespace()
	store.AddIndex(hash.NewInt64HashIndex(userID, true))
	asserts.Success(t, store.Insert(&User{
		ID:     1,
		Status: StatusActive,
	}))

	err := store.Insert(&User{
		ID:     1,
		Status: StatusDisabled,
	})

	asserts.Equals(t, "simd: record with passed id already exists: ID == 1", err.Error(), "check error")
}

func Test_Upsert(t *testing.T) {
	store := namespace.CreateNamespace()
	store.AddIndex(hash.NewInt64HashIndex(userID, true))
	store.AddIndex(hash.NewEnum8HashIndex(userStatus, true))
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

	err := store.Upsert(&User{
		ID:     2,
		Status: StatusActive,
	})
	asserts.Success(t, err)

	cur, err := executor.CreateQueryExecutor(store).FetchAll(
		context.Background(),
		query.NewBuilder(
			query.WhereInt64(userID, where.EQ, 2),
		).Query(),
	)

	asserts.Success(t, err)
	asserts.Success(t, cur.Err())
	asserts.Equals(t, StatusActive, cur.Item().(*User).Status, "status")
}

func Test_Delete(t *testing.T) {
	store := namespace.CreateNamespace()
	store.AddIndex(hash.NewInt64HashIndex(userID, true))
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

	err := store.Delete(2)
	asserts.Success(t, err)

	var ids []int64
	cur, err := executor.CreateQueryExecutor(store).FetchAll(
		context.Background(),
		query.NewBuilder(
			query.Sort(sort.ByInt64IndexAsc(&byID{})),
		).Query(),
	)
	asserts.Success(t, err)
	for cur.Next(context.Background()) {
		ids = append(ids, cur.Item().GetID())
	}
	asserts.Success(t, cur.Err())
	asserts.Equals(t, []int64{1, 3}, ids, "ids from cursor")
}
