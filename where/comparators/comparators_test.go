//nolint:exhaustruct
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

type enum16 uint16

type mp map[int]int

func (m mp) HasKey(key int) bool {
	_, ok := m[key]
	return ok
}
func (m mp) HasValue(check record.MapValueComparator[int]) (bool, error) {
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

type mapValueComparator func(item int) (bool, error)

func (e mapValueComparator) Compare(item int) (bool, error) {
	return e(item)
}

type set map[int]struct{}

func (s set) Has(item int) bool {
	_, ok := s[item]
	return ok
}

var fields = record.NewFields()

var boolGetter = record.BoolGetter[*user]{
	Field: fields.New("bool"),
	Get:   func(item *user) bool { return item.bool },
}

var enum8Getter = record.ComparableGetter[*user, enum8]{
	Field: fields.New("enum8"),
	Get:   func(item *user) enum8 { return item.enum8 },
}

var enum16Getter = record.ComparableGetter[*user, enum16]{
	Field: fields.New("enum16"),
	Get:   func(item *user) enum16 { return item.enum16 },
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

var ifaceGetter = record.Getter[*user, any]{
	Field: fields.New("iface"),
	Get:   func(item *user) any { return item.iface },
}

var mapGetter = record.MapGetter[*user, int, int]{
	Field: fields.New("map"),
	Get:   func(item *user) record.Map[int, int] { return item.mp },
}

var setGetter = record.SetGetter[*user, int]{
	Field: fields.New("set"),
	Get:   func(item *user) record.Set[int] { return item.set },
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

				values := make([]any, 0, test.comparator.ValuesCount())
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
				comparator: EqualComparator[*user, bool]{
					Cmp:    where.EQ,
					Getter: record.Getter[*user, bool](boolGetter),
					Value:  []bool{true},
				},
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "bool",
				expectedValues: []any{true},
			},
			{
				name: "true = false",
				comparator: EqualComparator[*user, bool]{
					Cmp:    where.EQ,
					Getter: record.Getter[*user, bool](boolGetter),
					Value:  []bool{false},
				},
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "bool",
				expectedValues: []any{false},
			},
			{
				name: "true ? true",
				comparator: EqualComparator[*user, bool]{
					Cmp:    0,
					Getter: record.Getter[*user, bool](boolGetter),
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
				name:           "2 = 2",
				comparator:     NewComparableFieldComparator[*user, enum8](where.EQ, enum8Getter, enum8(2)),
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "enum8",
				expectedValues: []any{enum8(2)},
			},
			{
				name:           "2 = 3",
				comparator:     NewComparableFieldComparator[*user, enum8](where.EQ, enum8Getter, enum8(3)),
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "enum8",
				expectedValues: []any{enum8(3)},
			},
			{
				name:           "2 IN (1, 2)",
				comparator:     NewComparableFieldComparator[*user, enum8](where.InArray, enum8Getter, enum8(1), enum8(2)),
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "enum8",
				expectedValues: []any{enum8(1), enum8(2)},
			},
			{
				name:           "2 IN (1, 3)",
				comparator:     NewComparableFieldComparator[*user, enum8](where.InArray, enum8Getter, enum8(1), enum8(3)),
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "enum8",
				expectedValues: []any{enum8(1), enum8(3)},
			},
			{
				name:           "2 ? 2",
				comparator:     NewComparableFieldComparator[*user, enum8](0, enum8Getter, enum8(2)),
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
				name:           "2 = 2",
				comparator:     NewComparableFieldComparator[*user, enum16](where.EQ, enum16Getter, enum16(2)),
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "enum16",
				expectedValues: []any{enum16(2)},
			},
			{
				name:           "2 = 3",
				comparator:     NewComparableFieldComparator[*user, enum16](where.EQ, enum16Getter, enum16(3)),
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "enum16",
				expectedValues: []any{enum16(3)},
			},
			{
				name:           "2 IN (1, 2)",
				comparator:     NewComparableFieldComparator[*user, enum16](where.InArray, enum16Getter, enum16(1), enum16(2)),
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "enum16",
				expectedValues: []any{enum16(1), enum16(2)},
			},
			{
				name:           "2 IN (1, 3)",
				comparator:     NewComparableFieldComparator[*user, enum16](where.InArray, enum16Getter, enum16(1), enum16(3)),
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "enum16",
				expectedValues: []any{enum16(1), enum16(3)},
			},
			{
				name:           "2 ? 2",
				comparator:     NewComparableFieldComparator[*user, enum16](0, enum16Getter, enum16(2)),
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
				name:           "10 = 10",
				comparator:     NewComparableFieldComparator[*user, int](where.EQ, intGetter, 10),
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "int",
				expectedValues: []any{10},
			},
			{
				name:           "10 = 3",
				comparator:     NewComparableFieldComparator[*user, int](where.EQ, intGetter, 3),
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "int",
				expectedValues: []any{3},
			},
			{
				name:           "10 > 3",
				comparator:     NewComparableFieldComparator[*user, int](where.GT, intGetter, 3),
				expectedResult: true,
				expectedCmp:    where.GT,
				expectedField:  "int",
				expectedValues: []any{3},
			},
			{
				name:           "10 > 30",
				comparator:     NewComparableFieldComparator[*user, int](where.GT, intGetter, 30),
				expectedResult: false,
				expectedCmp:    where.GT,
				expectedField:  "int",
				expectedValues: []any{30},
			},
			{
				name:           "10 >= 3",
				comparator:     NewComparableFieldComparator[*user, int](where.GE, intGetter, 3),
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int",
				expectedValues: []any{3},
			},
			{
				name:           "10 >= 30",
				comparator:     NewComparableFieldComparator(where.GE, intGetter, 30),
				expectedResult: false,
				expectedCmp:    where.GE,
				expectedField:  "int",
				expectedValues: []any{30},
			},
			{
				name:           "10 >= 10",
				comparator:     NewComparableFieldComparator[*user, int](where.GE, intGetter, 10),
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int",
				expectedValues: []any{10},
			},
			{
				name:           "10 < 3",
				comparator:     NewComparableFieldComparator[*user, int](where.LT, intGetter, 3),
				expectedResult: false,
				expectedCmp:    where.LT,
				expectedField:  "int",
				expectedValues: []any{3},
			},
			{
				name:           "10 < 30",
				comparator:     NewComparableFieldComparator[*user, int](where.LT, intGetter, 30),
				expectedResult: true,
				expectedCmp:    where.LT,
				expectedField:  "int",
				expectedValues: []any{30},
			},
			{
				name:           "10 <= 3",
				comparator:     NewComparableFieldComparator[*user, int](where.LE, intGetter, 3),
				expectedResult: false,
				expectedCmp:    where.LE,
				expectedField:  "int",
				expectedValues: []any{3},
			},
			{
				name:           "10 <= 30",
				comparator:     NewComparableFieldComparator[*user, int](where.LE, intGetter, 30),
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int",
				expectedValues: []any{30},
			},
			{
				name:           "10 <= 10",
				comparator:     NewComparableFieldComparator[*user, int](where.LE, intGetter, 10),
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int",
				expectedValues: []any{10},
			},
			{
				name:           "10 IN (1, 2, 10)",
				comparator:     NewComparableFieldComparator[*user, int](where.InArray, intGetter, 1, 2, 10),
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "int",
				expectedValues: []any{1, 2, 10},
			},
			{
				name:           "10 IN (1, 3)",
				comparator:     NewComparableFieldComparator[*user, int](where.InArray, intGetter, 1, 3),
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "int",
				expectedValues: []any{1, 3},
			},
			{
				name:           "10 ? 10",
				comparator:     NewComparableFieldComparator[*user, int](0, intGetter, 10),
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
				name:           "10 = 10",
				comparator:     NewComparableFieldComparator[*user, int32](where.EQ, int32Getter, 10),
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "int32",
				expectedValues: []any{int32(10)},
			},
			{
				name:           "10 = 3",
				comparator:     NewComparableFieldComparator[*user, int32](where.EQ, int32Getter, 3),
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "int32",
				expectedValues: []any{int32(3)},
			},
			{
				name:           "10 > 3",
				comparator:     NewComparableFieldComparator[*user, int32](where.GT, int32Getter, 3),
				expectedResult: true,
				expectedCmp:    where.GT,
				expectedField:  "int32",
				expectedValues: []any{int32(3)},
			},
			{
				name:           "10 > 30",
				comparator:     NewComparableFieldComparator[*user, int32](where.GT, int32Getter, 30),
				expectedResult: false,
				expectedCmp:    where.GT,
				expectedField:  "int32",
				expectedValues: []any{int32(30)},
			},
			{
				name:           "10 >= 3",
				comparator:     NewComparableFieldComparator[*user, int32](where.GE, int32Getter, 3),
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int32",
				expectedValues: []any{int32(3)},
			},
			{
				name:           "10 >= 30",
				comparator:     NewComparableFieldComparator[*user, int32](where.GE, int32Getter, 30),
				expectedResult: false,
				expectedCmp:    where.GE,
				expectedField:  "int32",
				expectedValues: []any{int32(30)},
			},
			{
				name:           "10 >= 10",
				comparator:     NewComparableFieldComparator[*user, int32](where.GE, int32Getter, 10),
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int32",
				expectedValues: []any{int32(10)},
			},
			{
				name:           "10 < 3",
				comparator:     NewComparableFieldComparator[*user, int32](where.LT, int32Getter, 3),
				expectedResult: false,
				expectedCmp:    where.LT,
				expectedField:  "int32",
				expectedValues: []any{int32(3)},
			},
			{
				name:           "10 < 30",
				comparator:     NewComparableFieldComparator[*user, int32](where.LT, int32Getter, 30),
				expectedResult: true,
				expectedCmp:    where.LT,
				expectedField:  "int32",
				expectedValues: []any{int32(30)},
			},
			{
				name:           "10 <= 3",
				comparator:     NewComparableFieldComparator[*user, int32](where.LE, int32Getter, 3),
				expectedResult: false,
				expectedCmp:    where.LE,
				expectedField:  "int32",
				expectedValues: []any{int32(3)},
			},
			{
				name:           "10 <= 30",
				comparator:     NewComparableFieldComparator[*user, int32](where.LE, int32Getter, 30),
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int32",
				expectedValues: []any{int32(30)},
			},
			{
				name:           "10 <= 10",
				comparator:     NewComparableFieldComparator[*user, int32](where.LE, int32Getter, 10),
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int32",
				expectedValues: []any{int32(10)},
			},
			{
				name:           "10 IN (1, 2, 10)",
				comparator:     NewComparableFieldComparator[*user, int32](where.InArray, int32Getter, 1, 2, 10),
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "int32",
				expectedValues: []any{int32(1), int32(2), int32(10)},
			},
			{
				name:           "10 IN (1, 3)",
				comparator:     NewComparableFieldComparator[*user, int32](where.InArray, int32Getter, 1, 3),
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "int32",
				expectedValues: []any{int32(1), int32(3)},
			},
			{
				name:           "10 ? 10",
				comparator:     NewComparableFieldComparator[*user, int32](0, int32Getter, 10),
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
				name:           "10 = 10",
				comparator:     NewComparableFieldComparator[*user, int64](where.EQ, int64Getter, 10),
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "int64",
				expectedValues: []any{int64(10)},
			},
			{
				name:           "10 = 3",
				comparator:     NewComparableFieldComparator[*user, int64](where.EQ, int64Getter, 3),
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "int64",
				expectedValues: []any{int64(3)},
			},
			{
				name:           "10 > 3",
				comparator:     NewComparableFieldComparator[*user, int64](where.GT, int64Getter, 3),
				expectedResult: true,
				expectedCmp:    where.GT,
				expectedField:  "int64",
				expectedValues: []any{int64(3)},
			},
			{
				name:           "10 > 30",
				comparator:     NewComparableFieldComparator[*user, int64](where.GT, int64Getter, 30),
				expectedResult: false,
				expectedCmp:    where.GT,
				expectedField:  "int64",
				expectedValues: []any{int64(30)},
			},
			{
				name:           "10 >= 3",
				comparator:     NewComparableFieldComparator[*user, int64](where.GE, int64Getter, 3),
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int64",
				expectedValues: []any{int64(3)},
			},
			{
				name:           "10 >= 30",
				comparator:     NewComparableFieldComparator[*user, int64](where.GE, int64Getter, 30),
				expectedResult: false,
				expectedCmp:    where.GE,
				expectedField:  "int64",
				expectedValues: []any{int64(30)},
			},
			{
				name:           "10 >= 10",
				comparator:     NewComparableFieldComparator[*user, int64](where.GE, int64Getter, 10),
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int64",
				expectedValues: []any{int64(10)},
			},
			{
				name:           "10 < 3",
				comparator:     NewComparableFieldComparator[*user, int64](where.LT, int64Getter, 3),
				expectedResult: false,
				expectedCmp:    where.LT,
				expectedField:  "int64",
				expectedValues: []any{int64(3)},
			},
			{
				name:           "10 < 30",
				comparator:     NewComparableFieldComparator[*user, int64](where.LT, int64Getter, 30),
				expectedResult: true,
				expectedCmp:    where.LT,
				expectedField:  "int64",
				expectedValues: []any{int64(30)},
			},
			{
				name:           "10 <= 3",
				comparator:     NewComparableFieldComparator[*user, int64](where.LE, int64Getter, 3),
				expectedResult: false,
				expectedCmp:    where.LE,
				expectedField:  "int64",
				expectedValues: []any{int64(3)},
			},
			{
				name:           "10 <= 30",
				comparator:     NewComparableFieldComparator[*user, int64](where.LE, int64Getter, 30),
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int64",
				expectedValues: []any{int64(30)},
			},
			{
				name:           "10 <= 10",
				comparator:     NewComparableFieldComparator[*user, int64](where.LE, int64Getter, 10),
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int64",
				expectedValues: []any{int64(10)},
			},
			{
				name:           "10 IN (1, 2, 10)",
				comparator:     NewComparableFieldComparator[*user, int64](where.InArray, int64Getter, 1, 2, 10),
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "int64",
				expectedValues: []any{int64(1), int64(2), int64(10)},
			},
			{
				name:           "10 IN (1, 3)",
				comparator:     NewComparableFieldComparator[*user, int64](where.InArray, int64Getter, 1, 3),
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "int64",
				expectedValues: []any{int64(1), int64(3)},
			},
			{
				name:           "10 ? 10",
				comparator:     NewComparableFieldComparator[*user, int64](0, int64Getter, 10),
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
				comparator: EqualComparator[*user, any]{
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
				comparator: EqualComparator[*user, any]{
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
				comparator: EqualComparator[*user, any]{
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
				comparator: EqualComparator[*user, any]{
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
				comparator: EqualComparator[*user, any]{
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
				name:           "MapHasKey 2",
				comparator:     NewMapFieldComparator[*user](where.MapHasKey, mapGetter, 2),
				expectedResult: true,
				expectedCmp:    where.MapHasKey,
				expectedField:  "map",
				expectedValues: []any{2},
			},
			{
				name:           "MapHasKey 4",
				comparator:     NewMapFieldComparator[*user](where.MapHasKey, mapGetter, 4),
				expectedResult: false,
				expectedCmp:    where.MapHasKey,
				expectedField:  "map",
				expectedValues: []any{4},
			},
			{
				name: "MapHasValue 8",
				comparator: NewMapFieldComparator[*user](
					where.MapHasValue,
					mapGetter,
					mapValueComparator(func(item int) (bool, error) {
						return item == 8, nil
					}),
				),
				expectedResult: true,
				expectedCmp:    where.MapHasValue,
				expectedField:  "map",
				expectedValues: []any{mapValueComparator(func(item int) (bool, error) {
					return item == 8, nil
				})},
			},
			{
				name: "MapHasValue 10",
				comparator: NewMapFieldComparator[*user](
					where.MapHasValue,
					mapGetter,
					mapValueComparator(func(item int) (bool, error) {
						return item == 10, nil
					}),
				),
				expectedResult: false,
				expectedCmp:    where.MapHasValue,
				expectedField:  "map",
				expectedValues: []any{mapValueComparator(func(item int) (bool, error) {
					return item == 10, nil
				})},
			},
			{
				name:           "MapHasValue cast error",
				comparator:     NewMapFieldComparator[*user](where.MapHasValue, mapGetter, 42),
				expectedResult: false,
				expectedError:  NewFailCastTypeError(mapGetter.Field, where.MapHasValue, 42, "record.MapValueComparator"),
				expectedCmp:    where.MapHasValue,
				expectedField:  "map",
				expectedValues: []any{42},
			},
			{
				name: "MapHasValue error",
				comparator: NewMapFieldComparator[*user](
					where.MapHasValue, mapGetter,
					mapValueComparator(func(item int) (bool, error) {
						return false, errors.New("comparator error")
					}),
				),
				expectedResult: false,
				expectedError:  errors.New("comparator error"),
				expectedCmp:    where.MapHasValue,
				expectedField:  "map",
				expectedValues: []any{mapValueComparator(func(item int) (bool, error) {
					return false, errors.New("comparator error")
				})},
			},
			{
				name:           "? 2",
				comparator:     NewMapFieldComparator[*user](0, mapGetter, 2),
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
				name:           "SetHas 2",
				comparator:     NewSetFieldComparator[*user, int](where.SetHas, setGetter, 2),
				expectedResult: true,
				expectedCmp:    where.SetHas,
				expectedField:  "set",
				expectedValues: []any{2},
			},
			{
				name:           "SetHas 3",
				comparator:     NewSetFieldComparator[*user, int](where.SetHas, setGetter, 3),
				expectedResult: false,
				expectedCmp:    where.SetHas,
				expectedField:  "set",
				expectedValues: []any{3},
			},
			{
				name:           "? 2",
				comparator:     NewSetFieldComparator[*user, int](0, setGetter, 2),
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
				name:           "SetHas 2",
				comparator:     NewSetFieldComparator[*user, int](where.SetHas, setGetter, 2),
				expectedResult: true,
				expectedCmp:    where.SetHas,
				expectedField:  "set",
				expectedValues: []any{2},
			},
			{
				name:           "SetHas 3",
				comparator:     NewSetFieldComparator[*user, int](where.SetHas, setGetter, 3),
				expectedResult: false,
				expectedCmp:    where.SetHas,
				expectedField:  "set",
				expectedValues: []any{3},
			},
			{
				name:           "? 2",
				comparator:     NewSetFieldComparator[*user, int](0, setGetter, 2),
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
				name:           "foo = foo",
				comparator:     NewStringFieldComparator[*user](where.EQ, stringGetter, "foo"),
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "string",
				expectedValues: []any{"foo"},
			},
			{
				name:           "foo = bar",
				comparator:     NewStringFieldComparator[*user](where.EQ, stringGetter, "bar"),
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "string",
				expectedValues: []any{"bar"},
			},
			{
				name:           "foo > bar",
				comparator:     NewStringFieldComparator[*user](where.GT, stringGetter, "bar"),
				expectedResult: true,
				expectedCmp:    where.GT,
				expectedField:  "string",
				expectedValues: []any{"bar"},
			},
			{
				name:           "foo > zzz",
				comparator:     NewStringFieldComparator[*user](where.GT, stringGetter, "zzz"),
				expectedResult: false,
				expectedCmp:    where.GT,
				expectedField:  "string",
				expectedValues: []any{"zzz"},
			},
			{
				name:           "foo >= bar",
				comparator:     NewStringFieldComparator[*user](where.GE, stringGetter, "bar"),
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "string",
				expectedValues: []any{"bar"},
			},
			{
				name:           "foo >= zzz",
				comparator:     NewStringFieldComparator[*user](where.GE, stringGetter, "zzz"),
				expectedResult: false,
				expectedCmp:    where.GE,
				expectedField:  "string",
				expectedValues: []any{"zzz"},
			},
			{
				name:           "foo >= foo",
				comparator:     NewStringFieldComparator[*user](where.GE, stringGetter, "foo"),
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "string",
				expectedValues: []any{"foo"},
			},
			{
				name:           "foo < bar",
				comparator:     NewStringFieldComparator[*user](where.LT, stringGetter, "bar"),
				expectedResult: false,
				expectedCmp:    where.LT,
				expectedField:  "string",
				expectedValues: []any{"bar"},
			},
			{
				name:           "foo < zzz",
				comparator:     NewStringFieldComparator[*user](where.LT, stringGetter, "zzz"),
				expectedResult: true,
				expectedCmp:    where.LT,
				expectedField:  "string",
				expectedValues: []any{"zzz"},
			},
			{
				name:           "foo <= bar",
				comparator:     NewStringFieldComparator[*user](where.LE, stringGetter, "bar"),
				expectedResult: false,
				expectedCmp:    where.LE,
				expectedField:  "string",
				expectedValues: []any{"bar"},
			},
			{
				name:           "foo <= zzz",
				comparator:     NewStringFieldComparator[*user](where.LE, stringGetter, "zzz"),
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "string",
				expectedValues: []any{"zzz"},
			},
			{
				name:           "foo <= foo",
				comparator:     NewStringFieldComparator[*user](where.LE, stringGetter, "foo"),
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "string",
				expectedValues: []any{"foo"},
			},
			{
				name:           "foo IN (bar, foo)",
				comparator:     NewStringFieldComparator[*user](where.InArray, stringGetter, "bar", "foo"),
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "string",
				expectedValues: []any{"bar", "foo"},
			},
			{
				name:           "foo IN (bar, zzz)",
				comparator:     NewStringFieldComparator[*user](where.InArray, stringGetter, "bar", "zzz"),
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "string",
				expectedValues: []any{"bar", "zzz"},
			},
			{
				name:           "foo LIKE oo",
				comparator:     NewStringFieldComparator[*user](where.Like, stringGetter, "oo"),
				expectedResult: true,
				expectedCmp:    where.Like,
				expectedField:  "string",
				expectedValues: []any{"oo"},
			},
			{
				name:           "foo LIKE ff",
				comparator:     NewStringFieldComparator[*user](where.Like, stringGetter, "ff"),
				expectedResult: false,
				expectedCmp:    where.Like,
				expectedField:  "string",
				expectedValues: []any{"ff"},
			},
			{
				name:           "foo ? bar",
				comparator:     NewStringFieldComparator[*user](0, stringGetter, "bar"),
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
				comparator: NewStringFieldRegexpComparator[*user](
					where.Regexp,
					stringGetter,
					regexp.MustCompile(`fo+`),
				),
				expectedResult: true,
				expectedCmp:    where.Regexp,
				expectedField:  "string",
				expectedValues: []any{regexp.MustCompile(`fo+`)},
			},
			{
				name: "foo Regexp /\\d+/",
				comparator: NewStringFieldRegexpComparator[*user](
					where.Regexp,
					stringGetter,
					regexp.MustCompile(`\d+`),
				),
				expectedResult: false,
				expectedCmp:    where.Regexp,
				expectedField:  "string",
				expectedValues: []any{regexp.MustCompile(`\d+`)},
			},
			{
				name: "foo ? fo+",
				comparator: NewStringFieldRegexpComparator[*user](
					0,
					stringGetter,
					regexp.MustCompile("fo+"),
				),
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(stringGetter.Field, 0),
				expectedCmp:    0,
				expectedField:  "string",
				expectedValues: []any{regexp.MustCompile("fo+")},
			},
		})
	})
}
