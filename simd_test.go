package simd

import (
	"context"
	"fmt"
	"github.com/shamcode/simd/asserts"
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/indexes/fields"
	"github.com/shamcode/simd/query"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"regexp"
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

type User struct {
	ID       int64
	Name     string
	Status   StatusEnum
	Score    int
	IsOnline bool
}

func (u *User) GetID() int64   { return u.ID }
func (u *User) ComputeFields() {}

var userID = &record.Int64Getter{
	Field: "id",
	Get:   func(item interface{}) int64 { return item.(*User).ID },
}

var userName = &record.StringGetter{
	Field: "name",
	Get:   func(item interface{}) string { return item.(*User).Name },
}

var userStatus = &record.Enum8Getter{
	Field: "status",
	Get:   func(item interface{}) record.Enum8 { return item.(*User).Status },
}

var userScore = &record.IntGetter{
	Field: "score",
	Get:   func(item interface{}) int { return item.(*User).Score },
}

var userIsOnline = &record.BoolGetter{
	Field: "is_online",
	Get:   func(item interface{}) bool { return item.(*User).IsOnline },
}

type byIDAsc struct{}

func (sorting *byIDAsc) CalcIndex(item record.Record) int64 { return item.(*User).ID }

type byIDDesc struct{}

func (sorting *byIDDesc) CalcIndex(item record.Record) int64 {
	return sort.Int64IndexDesc(item.(*User).ID)
}

func Test_FetchAllAndTotal(t *testing.T) {
	store := indexes.CreateNamespace()
	store.AddIndex(fields.NewInt64Index(userID))
	store.AddIndex(fields.NewStringIndex(userName))
	store.AddIndex(fields.NewEnum8Index(userStatus))
	store.AddIndex(fields.NewBoolIndex(userIsOnline))
	asserts.Success(t, store.Insert(&User{
		ID:     1,
		Name:   "First",
		Status: StatusActive,
		Score:  10,
	}))
	asserts.Success(t, store.Insert(&User{
		ID:     2,
		Name:   "Second",
		Status: StatusDisabled,
		Score:  15,
	}))
	asserts.Success(t, store.Insert(&User{
		ID:     3,
		Name:   "Third",
		Status: StatusDisabled,
		Score:  20,
	}))
	asserts.Success(t, store.Insert(&User{
		ID:       4,
		Name:     "Fourth",
		Status:   StatusActive,
		Score:    25,
		IsOnline: true,
	}))

	testCases := []struct {
		Name          string
		Query         query.Query
		ExpectedCount int
		ExpectedIDs   []int64
	}{
		{
			Name: "SELECT *, COUNT(*) WHERE status = ACTIVE ORDER BY id ASC",
			Query: query.NewBuilder().
				WhereEnum8(userStatus, where.EQ, StatusActive).
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE status = ACTIVE ORDER BY id DESC",
			Query: query.NewBuilder().
				WhereEnum8(userStatus, where.EQ, StatusActive).
				Sort(sort.ByInt64Index(&byIDDesc{})).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{4, 1},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE status != DISABLED ORDER BY id ASC",
			Query: query.NewBuilder().
				Not().
				WhereEnum8(userStatus, where.EQ, StatusDisabled).
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC",
			Query: query.NewBuilder().
				WhereInt(userScore, where.GE, 10).
				WhereInt(userScore, where.LT, 20).
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 2},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC LIMIT 1",
			Query: query.NewBuilder().
				WhereInt(userScore, where.GE, 10).
				WhereInt(userScore, where.LT, 20).
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Limit(1).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC OFFSET 1 LIMIT 3",
			Query: query.NewBuilder().
				WhereInt(userScore, where.GE, 10).
				WhereInt(userScore, where.LT, 20).
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Offset(1).
				Limit(3).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{2},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE score >= 10 AND score < 20 ORDER BY id ASC OFFSET 2 LIMIT 3",
			Query: query.NewBuilder().
				WhereInt(userScore, where.GE, 10).
				WhereInt(userScore, where.LT, 20).
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Offset(2).
				Limit(3).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name = 'Fourth' AND status == ACTIVE",
			Query: query.NewBuilder().
				WhereString(userName, where.EQ, "Fourth").
				WhereEnum8(userStatus, where.EQ, StatusActive).
				Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name = 'Fourth' AND status == DISABLED",
			Query: query.NewBuilder().
				WhereString(userName, where.EQ, "Fourth").
				WhereEnum8(userStatus, where.EQ, StatusDisabled).
				Query(),
			ExpectedCount: 0,
			ExpectedIDs:   []int64{},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name LIKE 'th' ORDER BY id ASC",
			Query: query.NewBuilder().
				WhereString(userName, where.Like, "t").
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR status == DISABLED ORDER BY id ASC",
			Query: query.NewBuilder().
				WhereInt64(userID, where.EQ, 1).
				Or().
				WhereEnum8(userStatus, where.EQ, StatusDisabled).
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR (NOT status == DISABLED) ORDER BY id ASC",
			Query: query.NewBuilder().
				WhereInt64(userID, where.EQ, 1).
				Or().
				OpenBracket().
				Not().
				WhereEnum8(userStatus, where.EQ, StatusDisabled).
				CloseBracket().
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{1, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR (NOT status == ACTIVE OR NOT is_online = true) ORDER BY id ASC",
			Query: query.NewBuilder().
				WhereInt64(userID, where.EQ, 1).
				Or().
				OpenBracket().
				Not().
				WhereEnum8(userStatus, where.EQ, StatusActive).
				Or().
				Not().
				WhereBool(userIsOnline, where.EQ, true).
				CloseBracket().
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 1 OR (status == DISABLED OR is_online = false) ORDER BY id ASC",
			Query: query.NewBuilder().
				WhereInt64(userID, where.EQ, 1).
				Or().
				OpenBracket().
				WhereEnum8(userStatus, where.EQ, StatusDisabled).
				Or().
				WhereBool(userIsOnline, where.EQ, false).
				CloseBracket().
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE (status == DISABLED OR is_online = false) OR id = 1 ORDER BY id ASC",
			Query: query.NewBuilder().
				OpenBracket().
				WhereEnum8(userStatus, where.EQ, StatusDisabled).
				Or().
				WhereBool(userIsOnline, where.EQ, false).
				CloseBracket().
				Or().
				WhereInt64(userID, where.EQ, 1).
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 4 OR (status == DISABLED OR is_online = false) ORDER BY id ASC",
			Query: query.NewBuilder().
				WhereInt64(userID, where.EQ, 4).
				Or().
				OpenBracket().
				WhereEnum8(userStatus, where.EQ, StatusDisabled).
				Or().
				WhereBool(userIsOnline, where.EQ, false).
				CloseBracket().
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{1, 2, 3, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id = 4 AND (status == DISABLED OR is_online = true) ORDER BY id ASC",
			Query: query.NewBuilder().
				WhereInt64(userID, where.EQ, 4).
				OpenBracket().
				WhereEnum8(userStatus, where.EQ, StatusDisabled).
				Or().
				WhereBool(userIsOnline, where.EQ, true).
				CloseBracket().
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE is_online = true ORDER BY id ASC",
			Query: query.NewBuilder().
				WhereBool(userIsOnline, where.EQ, true).
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE id IN (4, 2) ORDER BY id ASC",
			Query: query.NewBuilder().
				WhereInt64(userID, where.InArray, 4, 2).
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{2, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name REGEXP [tT]) ORDER BY id ASC",
			Query: query.NewBuilder().
				WhereStringRegexp(userName, regexp.MustCompile("[tT]")).
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 3,
			ExpectedIDs:   []int64{1, 3, 4},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE name IN (Second, Third) ORDER BY id ASC",
			Query: query.NewBuilder().
				WhereString(userName, where.InArray, "Second", "Third").
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 2,
			ExpectedIDs:   []int64{2, 3},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE ( id = 1 ) AND id IN (1, 2, 3) ORDER BY id ASC",
			Query: query.NewBuilder().
				OpenBracket().
				WhereInt64(userID, where.EQ, 1).
				CloseBracket().
				WhereInt64(userID, where.InArray, 1, 2, 3).
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 1,
			ExpectedIDs:   []int64{1},
		},
		{
			Name: "SELECT *, COUNT(*) WHERE True id ASC",
			Query: query.NewBuilder().
				Sort(sort.ByInt64Index(&byIDAsc{})).
				Query(),
			ExpectedCount: 4,
			ExpectedIDs:   []int64{1, 2, 3, 4},
		},
	}

	for _, testCase := range testCases {
		ctx := context.Background()
		cursor, count, err := store.QueryExecutor().FetchAllAndTotal(ctx, testCase.Query)
		asserts.Equals(t, nil, err, fmt.Sprintf("%s: nil == err", testCase.Name))
		ids := make([]int64, 0, cursor.Size())
		for cursor.Next(ctx) {
			ids = append(ids, cursor.Item().(*User).ID)
		}
		asserts.Equals(t, nil, cursor.Err(), fmt.Sprintf("%s: nil == cursor.Err", testCase.Name))
		asserts.Equals(t, testCase.ExpectedIDs, ids, fmt.Sprintf("%s: ids", testCase.Name))
		asserts.Equals(t, testCase.ExpectedCount, count, fmt.Sprintf("%s: total count", testCase.Name))
	}
}
