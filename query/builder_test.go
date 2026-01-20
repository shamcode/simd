package query

import (
	"testing"

	asserts "github.com/shamcode/assert"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

func TestBuilderErrors(t *testing.T) {
	_id := record.NewIDGetter[record.Record]()

	testCases := []struct {
		builderOptions []BuilderOption
		expectedError  string
	}{
		{
			builderOptions: []BuilderOption{
				Or(),
				Where(_id, where.EQ, 1),
				Where(_id, where.EQ, 2),
			},
			expectedError: ".Or() before any condition not supported, add any condition before .Or()",
		},
		{
			builderOptions: []BuilderOption{
				Not(),
				OpenBracket(),
				Where(_id, where.EQ, 1),
				Where(_id, where.EQ, 2),
				CloseBracket(),
			},
			expectedError: ".Not().OpenBracket() not supported",
		},
		{
			builderOptions: []BuilderOption{
				OpenBracket(),
				Where(_id, where.EQ, 1),
				Where(_id, where.EQ, 2),
				CloseBracket(),
				CloseBracket(),
			},
			expectedError: "close bracket without open",
		},
		{
			builderOptions: []BuilderOption{
				OpenBracket(),
				Where(_id, where.EQ, 1),
				Where(_id, where.EQ, 2),
				CloseBracket(),
				OpenBracket(),
			},
			expectedError: "invalid bracket balance: has not closed bracket",
		},
		{
			builderOptions: []BuilderOption{
				Not(),
				Or(),
				OpenBracket(),
				Where(_id, where.EQ, 1),
				Where(_id, where.EQ, 2),
				CloseBracket(),
				OpenBracket(),
			},
			expectedError: ".Or() before any condition not supported, add any condition before .Or()\n" +
				".Not().OpenBracket() not supported\n" +
				"invalid bracket balance: has not closed bracket",
		},
	}
	for _, testCase := range testCases {
		err := NewBuilder[record.Record](testCase.builderOptions...).Query().Error().Error()
		asserts.Equals(t, testCase.expectedError, err, "check expected error")
	}
}
