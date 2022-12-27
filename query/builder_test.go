package query

import (
	"github.com/shamcode/simd/asserts"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"testing"
)

type testUser struct {
	id int64
}

var userID = &record.Int64Getter{
	Field: "id",
	Get:   func(item interface{}) int64 { return item.(*testUser).id },
}

func TestBuilderErrors(t *testing.T) {
	testCases := []struct {
		query         Query
		expectedError string
	}{
		{
			query: NewBuilder().
				Or().
				WhereInt64(userID, where.EQ, 1).
				WhereInt64(userID, where.EQ, 2).
				Query(),
			expectedError: "1 error occurred:\n\t* .Or() before any condition not supported, add any condition before .Or()\n\n",
		},
		{
			query: NewBuilder().
				Not().
				OpenBracket().
				WhereInt64(userID, where.EQ, 1).
				WhereInt64(userID, where.EQ, 2).
				CloseBracket().
				Query(),
			expectedError: "1 error occurred:\n\t* .Not().OpenBracket() not supported\n\n",
		},
		{
			query: NewBuilder().
				OpenBracket().
				WhereInt64(userID, where.EQ, 1).
				WhereInt64(userID, where.EQ, 2).
				CloseBracket().
				CloseBracket().
				Query(),
			expectedError: "1 error occurred:\n\t* close bracket without open\n\n",
		},
		{
			query: NewBuilder().
				OpenBracket().
				WhereInt64(userID, where.EQ, 1).
				WhereInt64(userID, where.EQ, 2).
				CloseBracket().
				OpenBracket().
				Query(),
			expectedError: "1 error occurred:\n\t* invalid bracket balance: has not closed bracket\n\n",
		},
		{
			query: NewBuilder().
				Not().
				Or().
				OpenBracket().
				WhereInt64(userID, where.EQ, 1).
				WhereInt64(userID, where.EQ, 2).
				CloseBracket().
				OpenBracket().
				Query(),
			expectedError: "3 errors occurred:" +
				"\n\t* .Or() before any condition not supported, add any condition before .Or()" +
				"\n\t* .Not().OpenBracket() not supported" +
				"\n\t* invalid bracket balance: has not closed bracket\n\n",
		},
	}
	for _, testCase := range testCases {
		asserts.Equals(t, testCase.expectedError, testCase.query.Error().Error(), "check expected error")
	}
}
