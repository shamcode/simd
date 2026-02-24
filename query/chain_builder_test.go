package query

import (
	"testing"

	asserts "github.com/shamcode/assert"
	"github.com/shamcode/simd/record"
	. "github.com/shamcode/simd/where"
)

func TestChainBuilderErrors(t *testing.T) {
	_id := record.NewIDGetter[record.Record]()

	testCases := []struct {
		query         Query[record.Record]
		expectedError string
	}{
		{
			query: NewChainBuilder(NewBuilder[record.Record]()).
				Or().
				AddWhere(Where(_id, EQ, 1)).
				AddWhere(Where(_id, EQ, 1)).
				Query(),
			expectedError: ".Or() before any condition not supported, add any condition before .Or()",
		},
		{
			query: NewChainBuilder(NewBuilder[record.Record]()).
				Not().
				OpenBracket().
				AddWhere(Where(_id, EQ, 1)).
				AddWhere(Where(_id, EQ, 2)).
				CloseBracket().
				Query(),
			expectedError: ".Not().OpenBracket() not supported",
		},
		{
			query: NewChainBuilder(NewBuilder[record.Record]()).
				OpenBracket().
				AddWhere(Where(_id, EQ, 1)).
				AddWhere(Where(_id, EQ, 2)).
				CloseBracket().
				CloseBracket().
				Query(),
			expectedError: "close bracket without open",
		},
		{
			query: NewChainBuilder(NewBuilder[record.Record]()).
				OpenBracket().
				AddWhere(Where(_id, EQ, 1)).
				AddWhere(Where(_id, EQ, 2)).
				CloseBracket().
				OpenBracket().
				Query(),

			expectedError: "invalid bracket balance: has not closed bracket",
		},
		{
			query: NewChainBuilder(NewBuilder[record.Record]()).
				Not().
				Or().
				OpenBracket().
				AddWhere(Where(_id, EQ, 1)).
				AddWhere(Where(_id, EQ, 2)).
				CloseBracket().
				OpenBracket().
				Query(),
			expectedError: ".Or() before any condition not supported, add any condition before .Or()\n" +
				".Not().OpenBracket() not supported\n" +
				"invalid bracket balance: has not closed bracket",
		},
	}
	for _, testCase := range testCases {
		err := testCase.query.Error().Error()
		asserts.Equals(t, testCase.expectedError, err, "check expected error")
	}
}
