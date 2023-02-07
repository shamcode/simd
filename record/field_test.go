package record

import (
	"github.com/shamcode/simd/asserts"
	"testing"
)

func TestField(t *testing.T) {
	fields := NewFields()

	testCases := []struct {
		field         Field
		expectedName  string
		expectedIndex uint8
	}{
		{
			field:         ID.Field,
			expectedName:  "ID",
			expectedIndex: 0,
		},
		{
			field:         fields.New("name"),
			expectedName:  "name",
			expectedIndex: 1,
		},
		{
			field:         fields.New("score"),
			expectedName:  "score",
			expectedIndex: 2,
		},
	}

	for _, testCase := range testCases {
		asserts.Equals(t, testCase.expectedName, testCase.field.String(), "name")
		asserts.Equals(t, testCase.expectedIndex, testCase.field.Index(), "index")
	}
}
