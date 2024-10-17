package query

import (
	"testing"

	asserts "github.com/shamcode/assert"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

func TestBuilderErrors(t *testing.T) {
	testCases := []struct {
		builderOptions []BuilderOption
		expectedError  string
	}{
		{
			builderOptions: []BuilderOption{
				Or(),
				WhereInt64(record.ID, where.EQ, 1),
				WhereInt64(record.ID, where.EQ, 2),
			},
			expectedError: "1 error occurred:\n\t* .Or() before any condition not supported, add any condition before .Or()\n\n",
		},
		{
			builderOptions: []BuilderOption{
				Not(),
				OpenBracket(),
				WhereInt64(record.ID, where.EQ, 1),
				WhereInt64(record.ID, where.EQ, 2),
				CloseBracket(),
			},
			expectedError: "1 error occurred:\n\t* .Not().OpenBracket() not supported\n\n",
		},
		{
			builderOptions: []BuilderOption{
				OpenBracket(),
				WhereInt64(record.ID, where.EQ, 1),
				WhereInt64(record.ID, where.EQ, 2),
				CloseBracket(),
				CloseBracket(),
			},
			expectedError: "1 error occurred:\n\t* close bracket without open\n\n",
		},
		{
			builderOptions: []BuilderOption{
				OpenBracket(),
				WhereInt64(record.ID, where.EQ, 1),
				WhereInt64(record.ID, where.EQ, 2),
				CloseBracket(),
				OpenBracket(),
			},
			expectedError: "1 error occurred:\n\t* invalid bracket balance: has not closed bracket\n\n",
		},
		{
			builderOptions: []BuilderOption{
				Not(),
				Or(),
				OpenBracket(),
				WhereInt64(record.ID, where.EQ, 1),
				WhereInt64(record.ID, where.EQ, 2),
				CloseBracket(),
				OpenBracket(),
			},
			expectedError: "3 errors occurred:" +
				"\n\t* .Or() before any condition not supported, add any condition before .Or()" +
				"\n\t* .Not().OpenBracket() not supported" +
				"\n\t* invalid bracket balance: has not closed bracket\n\n",
		},
	}
	for _, testCase := range testCases {
		err := NewBuilder(testCase.builderOptions...).Query().Error().Error()
		asserts.Equals(t, testCase.expectedError, err, "check expected error")
	}
}
