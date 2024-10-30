//nolint:exhaustruct,err113
package comparators

import (
	"errors"
	"regexp"
	"testing"

	asserts "github.com/shamcode/assert"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
)

type user struct {
	bool   bool
	enum8  enum8
	enum16 enum16
	int    int
	int32  int32
	int64  int64
	iface  any
	mp     mp
	set    set
	string string
}

func (u *user) GetID() int64 { return u.int64 }

type enum8 uint8

func (e enum8) Value() uint8 { return uint8(e) }

type enum16 uint16

func (e enum16) Value() uint16 { return uint16(e) }

type mp map[int]int

func (m mp) HasKey(key any) bool {
	intKey, ok := key.(int)
	if !ok {
		return false
	}
	_, ok = m[intKey]
	return ok
}
func (m mp) HasValue(check record.MapValueComparator) (bool, error) {
	for _, value := range m {
		res, err := check.Compare(value)
		if nil != err {
			return false, err
		}
		if res {
			return true, nil
		}
	}
	return false, nil
}

type mapValueComparator func(item any) (bool, error)

func (e mapValueComparator) Compare(item any) (bool, error) {
	return e(item)
}

type set map[int]struct{}

func (s set) Has(item any) bool {
	intValue, ok := item.(int)
	if !ok {
		return false
	}
	_, ok = s[intValue]
	return ok
}

var fields = record.NewFields()

var boolGetter = record.BoolGetter[*user]{
	Field: fields.New("bool"),
	Get:   func(item *user) bool { return item.bool },
}

var enum8Getter = record.EnumGetter[*user, uint8]{
	Field: fields.New("enum8"),
	Get:   func(item *user) record.Enum[uint8] { return item.enum8 },
}

var enum16Getter = record.EnumGetter[*user, uint16]{
	Field: fields.New("enum16"),
	Get:   func(item *user) record.Enum[uint16] { return item.enum16 },
}

var intGetter = record.ComparableGetter[*user, int]{
	Field: fields.New("int"),
	Get:   func(item *user) int { return item.int },
}

var int32Getter = record.ComparableGetter[*user, int32]{
	Field: fields.New("int32"),
	Get:   func(item *user) int32 { return item.int32 },
}

var int64Getter = record.ComparableGetter[*user, int64]{
	Field: fields.New("int64"),
	Get:   func(item *user) int64 { return item.int64 },
}

var ifaceGetter = record.InterfaceGetter[*user]{
	Field: fields.New("iface"),
	Get:   func(item *user) any { return item.iface },
}

var mapGetter = record.MapGetter[*user]{
	Field: fields.New("map"),
	Get:   func(item *user) record.Map { return item.mp },
}

var setGetter = record.SetGetter[*user]{
	Field: fields.New("set"),
	Get:   func(item *user) record.Set { return item.set },
}

var stringGetter = record.ComparableGetter[*user, string]{
	Field: fields.New("string"),
	Get:   func(item *user) string { return item.string },
}

func TestComparators(t *testing.T) { //nolint:maintidx
	item := &user{
		bool:   true,
		enum8:  2,
		enum16: 2,
		int:    10,
		int32:  10,
		int64:  10,
		iface:  42,
		mp: map[int]int{
			1: 1,
			2: 4,
			3: 8,
		},
		set: map[int]struct{}{
			1: {},
			2: {},
		},
		string: "foo",
	}

	type testCase struct {
		name           string
		comparator     where.FieldComparator[*user]
		expectedResult bool
		expectedError  error
		expectedCmp    where.ComparatorType
		expectedField  string
		expectedValues []any
	}

	checkTestCases := func(t *testing.T, testCases []testCase) { //nolint:thelper
		for _, test := range testCases {
			t.Run(test.name, func(t *testing.T) {
				res, err := test.comparator.Compare(item)
				asserts.Equals(t, test.expectedResult, res, "result")
				asserts.Equals(t, test.expectedError, err, "error")
				asserts.Equals(t, test.comparator.GetType(), test.expectedCmp, "comparator type")
				asserts.Equals(t, test.comparator.GetField().String(), test.expectedField, "field")
				var values []any
				for i := range test.comparator.ValuesCount() {
					values = append(values, test.comparator.ValueAt(i))
				}
				asserts.Equals(t, values, test.expectedValues, "values")
			})
		}
	}

	t.Run("bool", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "true = true",
				comparator: BoolFieldComparator[*user]{
					Cmp:    where.EQ,
					Getter: boolGetter,
					Value:  []bool{true},
				},
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "bool",
				expectedValues: []any{true},
			},
			{
				name: "true = false",
				comparator: BoolFieldComparator[*user]{
					Cmp:    where.EQ,
					Getter: boolGetter,
					Value:  []bool{false},
				},
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "bool",
				expectedValues: []any{false},
			},
			{
				name: "true ? true",
				comparator: BoolFieldComparator[*user]{
					Cmp:    0,
					Getter: boolGetter,
					Value:  []bool{true},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(boolGetter.Field, 0),
				expectedCmp:    0,
				expectedField:  "bool",
				expectedValues: []any{true},
			},
		})
	})

	t.Run("enum8", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "2 = 2",
				comparator: EnumFieldComparator[*user, uint8]{
					Cmp:    where.EQ,
					Getter: enum8Getter,
					Value:  []record.Enum[uint8]{enum8(2)},
				},
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "enum8",
				expectedValues: []any{enum8(2)},
			},
			{
				name: "2 = 3",
				comparator: EnumFieldComparator[*user, uint8]{
					Cmp:    where.EQ,
					Getter: enum8Getter,
					Value:  []record.Enum[uint8]{enum8(3)},
				},
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "enum8",
				expectedValues: []any{enum8(3)},
			},
			{
				name: "2 IN (1, 2)",
				comparator: EnumFieldComparator[*user, uint8]{
					Cmp:    where.InArray,
					Getter: enum8Getter,
					Value:  []record.Enum[uint8]{enum8(1), enum8(2)},
				},
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "enum8",
				expectedValues: []any{enum8(1), enum8(2)},
			},
			{
				name: "2 IN (1, 3)",
				comparator: EnumFieldComparator[*user, uint8]{
					Cmp:    where.InArray,
					Getter: enum8Getter,
					Value:  []record.Enum[uint8]{enum8(1), enum8(3)},
				},
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "enum8",
				expectedValues: []any{enum8(1), enum8(3)},
			},
			{
				name: "2 ? 2",
				comparator: EnumFieldComparator[*user, uint8]{
					Cmp:    0,
					Getter: enum8Getter,
					Value:  []record.Enum[uint8]{enum8(2)},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(enum8Getter.Field, 0),
				expectedCmp:    0,
				expectedField:  "enum8",
				expectedValues: []any{enum8(2)},
			},
		})
	})

	t.Run("enum16", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "2 = 2",
				comparator: EnumFieldComparator[*user, uint16]{
					Cmp:    where.EQ,
					Getter: enum16Getter,
					Value:  []record.Enum[uint16]{enum16(2)},
				},
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "enum16",
				expectedValues: []any{enum16(2)},
			},
			{
				name: "2 = 3",
				comparator: EnumFieldComparator[*user, uint16]{
					Cmp:    where.EQ,
					Getter: enum16Getter,
					Value:  []record.Enum[uint16]{enum16(3)},
				},
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "enum16",
				expectedValues: []any{enum16(3)},
			},
			{
				name: "2 IN (1, 2)",
				comparator: EnumFieldComparator[*user, uint16]{
					Cmp:    where.InArray,
					Getter: enum16Getter,
					Value:  []record.Enum[uint16]{enum16(1), enum16(2)},
				},
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "enum16",
				expectedValues: []any{enum16(1), enum16(2)},
			},
			{
				name: "2 IN (1, 3)",
				comparator: EnumFieldComparator[*user, uint16]{
					Cmp:    where.InArray,
					Getter: enum16Getter,
					Value:  []record.Enum[uint16]{enum16(1), enum16(3)},
				},
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "enum16",
				expectedValues: []any{enum16(1), enum16(3)},
			},
			{
				name: "2 ? 2",
				comparator: EnumFieldComparator[*user, uint16]{
					Cmp:    0,
					Getter: enum16Getter,
					Value:  []record.Enum[uint16]{enum16(2)},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(enum16Getter.Field, 0),
				expectedCmp:    0,
				expectedField:  "enum16",
				expectedValues: []any{enum16(2)},
			},
		})
	})

	t.Run("int", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "10 = 10",
				comparator: ComparableFieldComparator[*user, int]{
					Cmp:    where.EQ,
					Getter: intGetter,
					Value:  []int{10},
				},
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "int",
				expectedValues: []any{10},
			},
			{
				name: "10 = 3",
				comparator: ComparableFieldComparator[*user, int]{
					Cmp:    where.EQ,
					Getter: intGetter,
					Value:  []int{3},
				},
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "int",
				expectedValues: []any{3},
			},
			{
				name: "10 > 3",
				comparator: ComparableFieldComparator[*user, int]{
					Cmp:    where.GT,
					Getter: intGetter,
					Value:  []int{3},
				},
				expectedResult: true,
				expectedCmp:    where.GT,
				expectedField:  "int",
				expectedValues: []any{3},
			},
			{
				name: "10 > 30",
				comparator: ComparableFieldComparator[*user, int]{
					Cmp:    where.GT,
					Getter: intGetter,
					Value:  []int{30},
				},
				expectedResult: false,
				expectedCmp:    where.GT,
				expectedField:  "int",
				expectedValues: []any{30},
			},
			{
				name: "10 >= 3",
				comparator: ComparableFieldComparator[*user, int]{
					Cmp:    where.GE,
					Getter: intGetter,
					Value:  []int{3},
				},
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int",
				expectedValues: []any{3},
			},
			{
				name: "10 >= 30",
				comparator: ComparableFieldComparator[*user, int]{
					Cmp:    where.GE,
					Getter: intGetter,
					Value:  []int{30},
				},
				expectedResult: false,
				expectedCmp:    where.GE,
				expectedField:  "int",
				expectedValues: []any{30},
			},
			{
				name: "10 >= 10",
				comparator: ComparableFieldComparator[*user, int]{
					Cmp:    where.GE,
					Getter: intGetter,
					Value:  []int{10},
				},
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int",
				expectedValues: []any{10},
			},
			{
				name: "10 < 3",
				comparator: ComparableFieldComparator[*user, int]{
					Cmp:    where.LT,
					Getter: intGetter,
					Value:  []int{3},
				},
				expectedResult: false,
				expectedCmp:    where.LT,
				expectedField:  "int",
				expectedValues: []any{3},
			},
			{
				name: "10 < 30",
				comparator: ComparableFieldComparator[*user, int]{
					Cmp:    where.LT,
					Getter: intGetter,
					Value:  []int{30},
				},
				expectedResult: true,
				expectedCmp:    where.LT,
				expectedField:  "int",
				expectedValues: []any{30},
			},
			{
				name: "10 <= 3",
				comparator: ComparableFieldComparator[*user, int]{
					Cmp:    where.LE,
					Getter: intGetter,
					Value:  []int{3},
				},
				expectedResult: false,
				expectedCmp:    where.LE,
				expectedField:  "int",
				expectedValues: []any{3},
			},
			{
				name: "10 <= 30",
				comparator: ComparableFieldComparator[*user, int]{
					Cmp:    where.LE,
					Getter: intGetter,
					Value:  []int{30},
				},
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int",
				expectedValues: []any{30},
			},
			{
				name: "10 <= 10",
				comparator: ComparableFieldComparator[*user, int]{
					Cmp:    where.LE,
					Getter: intGetter,
					Value:  []int{10},
				},
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int",
				expectedValues: []any{10},
			},
			{
				name: "10 IN (1, 2, 10)",
				comparator: ComparableFieldComparator[*user, int]{
					Cmp:    where.InArray,
					Getter: intGetter,
					Value:  []int{1, 2, 10},
				},
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "int",
				expectedValues: []any{1, 2, 10},
			},
			{
				name: "10 IN (1, 3)",
				comparator: ComparableFieldComparator[*user, int]{
					Cmp:    where.InArray,
					Getter: intGetter,
					Value:  []int{1, 3},
				},
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "int",
				expectedValues: []any{1, 3},
			},
			{
				name: "10 ? 10",
				comparator: ComparableFieldComparator[*user, int]{
					Cmp:    0,
					Getter: intGetter,
					Value:  []int{10},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(intGetter.Field, 0),
				expectedCmp:    0,
				expectedField:  "int",
				expectedValues: []any{10},
			},
		})
	})

	t.Run("int32", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "10 = 10",
				comparator: ComparableFieldComparator[*user, int32]{
					Cmp:    where.EQ,
					Getter: int32Getter,
					Value:  []int32{10},
				},
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "int32",
				expectedValues: []any{int32(10)},
			},
			{
				name: "10 = 3",
				comparator: ComparableFieldComparator[*user, int32]{
					Cmp:    where.EQ,
					Getter: int32Getter,
					Value:  []int32{3},
				},
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "int32",
				expectedValues: []any{int32(3)},
			},
			{
				name: "10 > 3",
				comparator: ComparableFieldComparator[*user, int32]{
					Cmp:    where.GT,
					Getter: int32Getter,
					Value:  []int32{3},
				},
				expectedResult: true,
				expectedCmp:    where.GT,
				expectedField:  "int32",
				expectedValues: []any{int32(3)},
			},
			{
				name: "10 > 30",
				comparator: ComparableFieldComparator[*user, int32]{
					Cmp:    where.GT,
					Getter: int32Getter,
					Value:  []int32{30},
				},
				expectedResult: false,
				expectedCmp:    where.GT,
				expectedField:  "int32",
				expectedValues: []any{int32(30)},
			},
			{
				name: "10 >= 3",
				comparator: ComparableFieldComparator[*user, int32]{
					Cmp:    where.GE,
					Getter: int32Getter,
					Value:  []int32{3},
				},
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int32",
				expectedValues: []any{int32(3)},
			},
			{
				name: "10 >= 30",
				comparator: ComparableFieldComparator[*user, int32]{
					Cmp:    where.GE,
					Getter: int32Getter,
					Value:  []int32{30},
				},
				expectedResult: false,
				expectedCmp:    where.GE,
				expectedField:  "int32",
				expectedValues: []any{int32(30)},
			},
			{
				name: "10 >= 10",
				comparator: ComparableFieldComparator[*user, int32]{
					Cmp:    where.GE,
					Getter: int32Getter,
					Value:  []int32{10},
				},
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int32",
				expectedValues: []any{int32(10)},
			},
			{
				name: "10 < 3",
				comparator: ComparableFieldComparator[*user, int32]{
					Cmp:    where.LT,
					Getter: int32Getter,
					Value:  []int32{3},
				},
				expectedResult: false,
				expectedCmp:    where.LT,
				expectedField:  "int32",
				expectedValues: []any{int32(3)},
			},
			{
				name: "10 < 30",
				comparator: ComparableFieldComparator[*user, int32]{
					Cmp:    where.LT,
					Getter: int32Getter,
					Value:  []int32{30},
				},
				expectedResult: true,
				expectedCmp:    where.LT,
				expectedField:  "int32",
				expectedValues: []any{int32(30)},
			},
			{
				name: "10 <= 3",
				comparator: ComparableFieldComparator[*user, int32]{
					Cmp:    where.LE,
					Getter: int32Getter,
					Value:  []int32{3},
				},
				expectedResult: false,
				expectedCmp:    where.LE,
				expectedField:  "int32",
				expectedValues: []any{int32(3)},
			},
			{
				name: "10 <= 30",
				comparator: ComparableFieldComparator[*user, int32]{
					Cmp:    where.LE,
					Getter: int32Getter,
					Value:  []int32{30},
				},
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int32",
				expectedValues: []any{int32(30)},
			},
			{
				name: "10 <= 10",
				comparator: ComparableFieldComparator[*user, int32]{
					Cmp:    where.LE,
					Getter: int32Getter,
					Value:  []int32{10},
				},
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int32",
				expectedValues: []any{int32(10)},
			},
			{
				name: "10 IN (1, 2, 10)",
				comparator: ComparableFieldComparator[*user, int32]{
					Cmp:    where.InArray,
					Getter: int32Getter,
					Value:  []int32{1, 2, 10},
				},
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "int32",
				expectedValues: []any{int32(1), int32(2), int32(10)},
			},
			{
				name: "10 IN (1, 3)",
				comparator: ComparableFieldComparator[*user, int32]{
					Cmp:    where.InArray,
					Getter: int32Getter,
					Value:  []int32{1, 3},
				},
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "int32",
				expectedValues: []any{int32(1), int32(3)},
			},
			{
				name: "10 ? 10",
				comparator: ComparableFieldComparator[*user, int32]{
					Cmp:    0,
					Getter: int32Getter,
					Value:  []int32{10},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(int32Getter.Field, 0),
				expectedCmp:    0,
				expectedField:  "int32",
				expectedValues: []any{int32(10)},
			},
		})
	})

	t.Run("int64", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "10 = 10",
				comparator: ComparableFieldComparator[*user, int64]{
					Cmp:    where.EQ,
					Getter: int64Getter,
					Value:  []int64{10},
				},
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "int64",
				expectedValues: []any{int64(10)},
			},
			{
				name: "10 = 3",
				comparator: ComparableFieldComparator[*user, int64]{
					Cmp:    where.EQ,
					Getter: int64Getter,
					Value:  []int64{3},
				},
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "int64",
				expectedValues: []any{int64(3)},
			},
			{
				name: "10 > 3",
				comparator: ComparableFieldComparator[*user, int64]{
					Cmp:    where.GT,
					Getter: int64Getter,
					Value:  []int64{3},
				},
				expectedResult: true,
				expectedCmp:    where.GT,
				expectedField:  "int64",
				expectedValues: []any{int64(3)},
			},
			{
				name: "10 > 30",
				comparator: ComparableFieldComparator[*user, int64]{
					Cmp:    where.GT,
					Getter: int64Getter,
					Value:  []int64{30},
				},
				expectedResult: false,
				expectedCmp:    where.GT,
				expectedField:  "int64",
				expectedValues: []any{int64(30)},
			},
			{
				name: "10 >= 3",
				comparator: ComparableFieldComparator[*user, int64]{
					Cmp:    where.GE,
					Getter: int64Getter,
					Value:  []int64{3},
				},
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int64",
				expectedValues: []any{int64(3)},
			},
			{
				name: "10 >= 30",
				comparator: ComparableFieldComparator[*user, int64]{
					Cmp:    where.GE,
					Getter: int64Getter,
					Value:  []int64{30},
				},
				expectedResult: false,
				expectedCmp:    where.GE,
				expectedField:  "int64",
				expectedValues: []any{int64(30)},
			},
			{
				name: "10 >= 10",
				comparator: ComparableFieldComparator[*user, int64]{
					Cmp:    where.GE,
					Getter: int64Getter,
					Value:  []int64{10},
				},
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int64",
				expectedValues: []any{int64(10)},
			},
			{
				name: "10 < 3",
				comparator: ComparableFieldComparator[*user, int64]{
					Cmp:    where.LT,
					Getter: int64Getter,
					Value:  []int64{3},
				},
				expectedResult: false,
				expectedCmp:    where.LT,
				expectedField:  "int64",
				expectedValues: []any{int64(3)},
			},
			{
				name: "10 < 30",
				comparator: ComparableFieldComparator[*user, int64]{
					Cmp:    where.LT,
					Getter: int64Getter,
					Value:  []int64{30},
				},
				expectedResult: true,
				expectedCmp:    where.LT,
				expectedField:  "int64",
				expectedValues: []any{int64(30)},
			},
			{
				name: "10 <= 3",
				comparator: ComparableFieldComparator[*user, int64]{
					Cmp:    where.LE,
					Getter: int64Getter,
					Value:  []int64{3},
				},
				expectedResult: false,
				expectedCmp:    where.LE,
				expectedField:  "int64",
				expectedValues: []any{int64(3)},
			},
			{
				name: "10 <= 30",
				comparator: ComparableFieldComparator[*user, int64]{
					Cmp:    where.LE,
					Getter: int64Getter,
					Value:  []int64{30},
				},
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int64",
				expectedValues: []any{int64(30)},
			},
			{
				name: "10 <= 10",
				comparator: ComparableFieldComparator[*user, int64]{
					Cmp:    where.LE,
					Getter: int64Getter,
					Value:  []int64{10},
				},
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int64",
				expectedValues: []any{int64(10)},
			},
			{
				name: "10 IN (1, 2, 10)",
				comparator: ComparableFieldComparator[*user, int64]{
					Cmp:    where.InArray,
					Getter: int64Getter,
					Value:  []int64{1, 2, 10},
				},
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "int64",
				expectedValues: []any{int64(1), int64(2), int64(10)},
			},
			{
				name: "10 IN (1, 3)",
				comparator: ComparableFieldComparator[*user, int64]{
					Cmp:    where.InArray,
					Getter: int64Getter,
					Value:  []int64{1, 3},
				},
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "int64",
				expectedValues: []any{int64(1), int64(3)},
			},
			{
				name: "10 ? 10",
				comparator: ComparableFieldComparator[*user, int64]{
					Cmp:    0,
					Getter: int64Getter,
					Value:  []int64{10},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(int64Getter.Field, 0),
				expectedCmp:    0,
				expectedField:  "int64",
				expectedValues: []any{int64(10)},
			},
		})
	})

	t.Run("any", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "42 = 42",
				comparator: InterfaceFieldComparator[*user]{
					Cmp:    where.EQ,
					Getter: ifaceGetter,
					Value:  []any{42},
				},
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "iface",
				expectedValues: []any{42},
			},
			{
				name: "42 = 10",
				comparator: InterfaceFieldComparator[*user]{
					Cmp:    where.EQ,
					Getter: ifaceGetter,
					Value:  []any{10},
				},
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "iface",
				expectedValues: []any{10},
			},
			{
				name: "42 IN (10, 42)",
				comparator: InterfaceFieldComparator[*user]{
					Cmp:    where.InArray,
					Getter: ifaceGetter,
					Value:  []any{10, 42},
				},
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "iface",
				expectedValues: []any{10, 42},
			},
			{
				name: "42 IN (10, 4)",
				comparator: InterfaceFieldComparator[*user]{
					Cmp:    where.InArray,
					Getter: ifaceGetter,
					Value:  []any{10, 4},
				},
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "iface",
				expectedValues: []any{10, 4},
			},
			{
				name: "42 ? 2",
				comparator: InterfaceFieldComparator[*user]{
					Cmp:    0,
					Getter: ifaceGetter,
					Value:  []any{2},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(ifaceGetter.Field, 0),
				expectedCmp:    0,
				expectedField:  "iface",
				expectedValues: []any{2},
			},
		})
	})

	t.Run("map", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "MapHasKey 2",
				comparator: MapFieldComparator[*user]{
					Cmp:    where.MapHasKey,
					Getter: mapGetter,
					Value:  []any{2},
				},
				expectedResult: true,
				expectedCmp:    where.MapHasKey,
				expectedField:  "map",
				expectedValues: []any{2},
			},
			{
				name: "MapHasKey 4",
				comparator: MapFieldComparator[*user]{
					Cmp:    where.MapHasKey,
					Getter: mapGetter,
					Value:  []any{4},
				},
				expectedResult: false,
				expectedCmp:    where.MapHasKey,
				expectedField:  "map",
				expectedValues: []any{4},
			},
			{
				name: "MapHasValue 8",
				comparator: MapFieldComparator[*user]{
					Cmp:    where.MapHasValue,
					Getter: mapGetter,
					Value: []any{mapValueComparator(func(item any) (bool, error) {
						return item.(int) == 8, nil
					})},
				},
				expectedResult: true,
				expectedCmp:    where.MapHasValue,
				expectedField:  "map",
				expectedValues: []any{mapValueComparator(func(item any) (bool, error) {
					return item.(int) == 8, nil
				})},
			},
			{
				name: "MapHasValue 10",
				comparator: MapFieldComparator[*user]{
					Cmp:    where.MapHasValue,
					Getter: mapGetter,
					Value: []any{mapValueComparator(func(item any) (bool, error) {
						return item.(int) == 10, nil
					})},
				},
				expectedResult: false,
				expectedCmp:    where.MapHasValue,
				expectedField:  "map",
				expectedValues: []any{mapValueComparator(func(item any) (bool, error) {
					return item.(int) == 10, nil
				})},
			},
			{
				name: "MapHasValue cast error",
				comparator: MapFieldComparator[*user]{
					Cmp:    where.MapHasValue,
					Getter: mapGetter,
					Value:  []any{42},
				},
				expectedResult: false,
				expectedError:  NewFailCastTypeError(mapGetter.Field, where.MapHasValue, 42, "record.MapValueComparator"),
				expectedCmp:    where.MapHasValue,
				expectedField:  "map",
				expectedValues: []any{42},
			},
			{
				name: "MapHasValue error",
				comparator: MapFieldComparator[*user]{
					Cmp:    where.MapHasValue,
					Getter: mapGetter,
					Value: []any{mapValueComparator(func(item any) (bool, error) {
						return false, errors.New("comparator error")
					})},
				},
				expectedResult: false,
				expectedError:  errors.New("comparator error"),
				expectedCmp:    where.MapHasValue,
				expectedField:  "map",
				expectedValues: []any{mapValueComparator(func(item any) (bool, error) {
					return false, errors.New("comparator error")
				})},
			},
			{
				name: "? 2",
				comparator: MapFieldComparator[*user]{
					Cmp:    0,
					Getter: mapGetter,
					Value:  []any{2},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(mapGetter.Field, 0),
				expectedCmp:    0,
				expectedField:  "map",
				expectedValues: []any{2},
			},
		})
	})

	t.Run("set", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "SetHas 2",
				comparator: SetFieldComparator[*user]{
					Cmp:    where.SetHas,
					Getter: setGetter,
					Value:  []any{2},
				},
				expectedResult: true,
				expectedCmp:    where.SetHas,
				expectedField:  "set",
				expectedValues: []any{2},
			},
			{
				name: "SetHas 3",
				comparator: SetFieldComparator[*user]{
					Cmp:    where.SetHas,
					Getter: setGetter,
					Value:  []any{3},
				},
				expectedResult: false,
				expectedCmp:    where.SetHas,
				expectedField:  "set",
				expectedValues: []any{3},
			},
			{
				name: "? 2",
				comparator: SetFieldComparator[*user]{
					Cmp:    0,
					Getter: setGetter,
					Value:  []any{2},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(setGetter.Field, 0),
				expectedCmp:    0,
				expectedField:  "set",
				expectedValues: []any{2},
			},
		})
	})

	t.Run("set", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "SetHas 2",
				comparator: SetFieldComparator[*user]{
					Cmp:    where.SetHas,
					Getter: setGetter,
					Value:  []any{2},
				},
				expectedResult: true,
				expectedCmp:    where.SetHas,
				expectedField:  "set",
				expectedValues: []any{2},
			},
			{
				name: "SetHas 3",
				comparator: SetFieldComparator[*user]{
					Cmp:    where.SetHas,
					Getter: setGetter,
					Value:  []any{3},
				},
				expectedResult: false,
				expectedCmp:    where.SetHas,
				expectedField:  "set",
				expectedValues: []any{3},
			},
			{
				name: "? 2",
				comparator: SetFieldComparator[*user]{
					Cmp:    0,
					Getter: setGetter,
					Value:  []any{2},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(setGetter.Field, 0),
				expectedCmp:    0,
				expectedField:  "set",
				expectedValues: []any{2},
			},
		})
	})

	t.Run("string", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "foo = foo",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.EQ,
						Getter: stringGetter,
						Value:  []string{"foo"},
					},
				},
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "string",
				expectedValues: []any{"foo"},
			},
			{
				name: "foo = bar",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.EQ,
						Getter: stringGetter,
						Value:  []string{"bar"},
					},
				},
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "string",
				expectedValues: []any{"bar"},
			},
			{
				name: "foo > bar",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.GT,
						Getter: stringGetter,
						Value:  []string{"bar"},
					},
				},
				expectedResult: true,
				expectedCmp:    where.GT,
				expectedField:  "string",
				expectedValues: []any{"bar"},
			},
			{
				name: "foo > zzz",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.GT,
						Getter: stringGetter,
						Value:  []string{"zzz"},
					},
				},
				expectedResult: false,
				expectedCmp:    where.GT,
				expectedField:  "string",
				expectedValues: []any{"zzz"},
			},
			{
				name: "foo >= bar",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.GE,
						Getter: stringGetter,
						Value:  []string{"bar"},
					},
				},
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "string",
				expectedValues: []any{"bar"},
			},
			{
				name: "foo >= zzz",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.GE,
						Getter: stringGetter,
						Value:  []string{"zzz"},
					},
				},
				expectedResult: false,
				expectedCmp:    where.GE,
				expectedField:  "string",
				expectedValues: []any{"zzz"},
			},
			{
				name: "foo >= foo",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.GE,
						Getter: stringGetter,
						Value:  []string{"foo"},
					},
				},
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "string",
				expectedValues: []any{"foo"},
			},
			{
				name: "foo < bar",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.LT,
						Getter: stringGetter,
						Value:  []string{"bar"},
					},
				},
				expectedResult: false,
				expectedCmp:    where.LT,
				expectedField:  "string",
				expectedValues: []any{"bar"},
			},
			{
				name: "foo < zzz",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.LT,
						Getter: stringGetter,
						Value:  []string{"zzz"},
					},
				},
				expectedResult: true,
				expectedCmp:    where.LT,
				expectedField:  "string",
				expectedValues: []any{"zzz"},
			},
			{
				name: "foo <= bar",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.LE,
						Getter: stringGetter,
						Value:  []string{"bar"},
					},
				},
				expectedResult: false,
				expectedCmp:    where.LE,
				expectedField:  "string",
				expectedValues: []any{"bar"},
			},
			{
				name: "foo <= zzz",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.LE,
						Getter: stringGetter,
						Value:  []string{"zzz"},
					},
				},
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "string",
				expectedValues: []any{"zzz"},
			},
			{
				name: "foo <= foo",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.LE,
						Getter: stringGetter,
						Value:  []string{"foo"},
					},
				},
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "string",
				expectedValues: []any{"foo"},
			},
			{
				name: "foo IN (bar, foo)",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.InArray,
						Getter: stringGetter,
						Value:  []string{"bar", "foo"},
					},
				},
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "string",
				expectedValues: []any{"bar", "foo"},
			},
			{
				name: "foo IN (bar, zzz)",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.InArray,
						Getter: stringGetter,
						Value:  []string{"bar", "zzz"},
					},
				},
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "string",
				expectedValues: []any{"bar", "zzz"},
			},
			{
				name: "foo LIKE oo",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.Like,
						Getter: stringGetter,
						Value:  []string{"oo"},
					},
				},
				expectedResult: true,
				expectedCmp:    where.Like,
				expectedField:  "string",
				expectedValues: []any{"oo"},
			},
			{
				name: "foo LIKE ff",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    where.Like,
						Getter: stringGetter,
						Value:  []string{"ff"},
					},
				},
				expectedResult: false,
				expectedCmp:    where.Like,
				expectedField:  "string",
				expectedValues: []any{"ff"},
			},
			{
				name: "foo ? bar",
				comparator: StringFieldComparator[*user]{
					ComparableFieldComparator: ComparableFieldComparator[*user, string]{
						Cmp:    0,
						Getter: stringGetter,
						Value:  []string{"bar"},
					},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(stringGetter.Field, 0),
				expectedCmp:    0,
				expectedField:  "string",
				expectedValues: []any{"bar"},
			},
		})
	})

	t.Run("string regexp", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "foo Regexp /fo+/",
				comparator: StringFieldRegexpComparator[*user]{
					Cmp:    where.Regexp,
					Getter: stringGetter,
					Value:  regexp.MustCompile(`fo+`),
				},
				expectedResult: true,
				expectedCmp:    where.Regexp,
				expectedField:  "string",
				expectedValues: []any{regexp.MustCompile(`fo+`)},
			},
			{
				name: "foo Regexp /\\d+/",
				comparator: StringFieldRegexpComparator[*user]{
					Cmp:    where.Regexp,
					Getter: stringGetter,
					Value:  regexp.MustCompile(`\d+`),
				},
				expectedResult: false,
				expectedCmp:    where.Regexp,
				expectedField:  "string",
				expectedValues: []any{regexp.MustCompile(`\d+`)},
			},
			{
				name: "foo ? fo+",
				comparator: StringFieldRegexpComparator[*user]{
					Cmp:    0,
					Getter: stringGetter,
					Value:  regexp.MustCompile("fo+"),
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(stringGetter.Field, 0),
				expectedCmp:    0,
				expectedField:  "string",
				expectedValues: []any{regexp.MustCompile("fo+")},
			},
		})
	})
}
