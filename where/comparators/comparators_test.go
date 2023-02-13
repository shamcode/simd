package comparators

import (
	"github.com/shamcode/simd/asserts"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"testing"
)

type user struct {
	bool   bool
	enum8  enum8
	enum16 enum16
	int    int
	int32  int32
	int64  int64
	iface  interface{}
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

func (m mp) HasKey(key interface{}) bool {
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

type mapValueComparator func(item interface{}) (bool, error)

func (e mapValueComparator) Compare(item interface{}) (bool, error) {
	return e(item)
}

type set map[int]struct{}

func (s set) Has(item interface{}) bool {
	intValue, ok := item.(int)
	if !ok {
		return false
	}
	_, ok = s[intValue]
	return ok
}

var fields = record.NewFields()

var boolGetter = record.BoolGetter{
	Field: fields.New("bool"),
	Get:   func(item record.Record) bool { return item.(*user).bool },
}

var enum8Getter = record.Enum8Getter{
	Field: fields.New("enum8"),
	Get:   func(item record.Record) record.Enum8 { return item.(*user).enum8 },
}

var enum16Getter = record.Enum16Getter{
	Field: fields.New("enum16"),
	Get:   func(item record.Record) record.Enum16 { return item.(*user).enum16 },
}

var intGetter = record.IntGetter{
	Field: fields.New("int"),
	Get:   func(item record.Record) int { return item.(*user).int },
}

var int32Getter = record.Int32Getter{
	Field: fields.New("int32"),
	Get:   func(item record.Record) int32 { return item.(*user).int32 },
}

var int64Getter = record.Int64Getter{
	Field: fields.New("int64"),
	Get:   func(item record.Record) int64 { return item.(*user).int64 },
}

func TestComparators(t *testing.T) {
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
		comparator     where.FieldComparator
		expectedResult bool
		expectedError  error
		expectedCmp    where.ComparatorType
		expectedField  string
		expectedValues []interface{}
	}

	checkTestCases := func(t *testing.T, testCases []testCase) {
		for _, test := range testCases {
			t.Run(test.name, func(t *testing.T) {
				res, err := test.comparator.Compare(item)
				asserts.Equals(t, test.expectedResult, res, "result")
				asserts.Equals(t, test.expectedError, err, "error")
				asserts.Equals(t, test.comparator.GetType(), test.expectedCmp, "comparator type")
				asserts.Equals(t, test.comparator.GetField().String(), test.expectedField, "field")
				var values []interface{}
				for i := 0; i < test.comparator.ValuesCount(); i++ {
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
				comparator: BoolFieldComparator{
					Cmp:    where.EQ,
					Getter: boolGetter,
					Value:  []bool{true},
				},
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "bool",
				expectedValues: []interface{}{true},
			},
			{
				name: "true = false",
				comparator: BoolFieldComparator{
					Cmp:    where.EQ,
					Getter: boolGetter,
					Value:  []bool{false},
				},
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "bool",
				expectedValues: []interface{}{false},
			},
			{
				name: "true ? true",
				comparator: BoolFieldComparator{
					Cmp:    0,
					Getter: boolGetter,
					Value:  []bool{true},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(boolGetter.Field, 0),
				expectedCmp:    0,
				expectedField:  "bool",
				expectedValues: []interface{}{true},
			},
		})
	})

	t.Run("enum8", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "2 = 2",
				comparator: Enum8FieldComparator{
					Cmp:    where.EQ,
					Getter: enum8Getter,
					Value:  []record.Enum8{enum8(2)},
				},
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "enum8",
				expectedValues: []interface{}{enum8(2)},
			},
			{
				name: "2 = 3",
				comparator: Enum8FieldComparator{
					Cmp:    where.EQ,
					Getter: enum8Getter,
					Value:  []record.Enum8{enum8(3)},
				},
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "enum8",
				expectedValues: []interface{}{enum8(3)},
			},
			{
				name: "2 IN (1, 2)",
				comparator: Enum8FieldComparator{
					Cmp:    where.InArray,
					Getter: enum8Getter,
					Value:  []record.Enum8{enum8(1), enum8(2)},
				},
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "enum8",
				expectedValues: []interface{}{enum8(1), enum8(2)},
			},
			{
				name: "2 IN (1, 3)",
				comparator: Enum8FieldComparator{
					Cmp:    where.InArray,
					Getter: enum8Getter,
					Value:  []record.Enum8{enum8(1), enum8(3)},
				},
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "enum8",
				expectedValues: []interface{}{enum8(1), enum8(3)},
			},
			{
				name: "2 ? 2",
				comparator: Enum8FieldComparator{
					Cmp:    0,
					Getter: enum8Getter,
					Value:  []record.Enum8{enum8(2)},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(enum8Getter.Field, 0),
				expectedCmp:    0,
				expectedField:  "enum8",
				expectedValues: []interface{}{enum8(2)},
			},
		})
	})

	t.Run("enum16", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "2 = 2",
				comparator: Enum16FieldComparator{
					Cmp:    where.EQ,
					Getter: enum16Getter,
					Value:  []record.Enum16{enum16(2)},
				},
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "enum16",
				expectedValues: []interface{}{enum16(2)},
			},
			{
				name: "2 = 3",
				comparator: Enum16FieldComparator{
					Cmp:    where.EQ,
					Getter: enum16Getter,
					Value:  []record.Enum16{enum16(3)},
				},
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "enum16",
				expectedValues: []interface{}{enum16(3)},
			},
			{
				name: "2 IN (1, 2)",
				comparator: Enum16FieldComparator{
					Cmp:    where.InArray,
					Getter: enum16Getter,
					Value:  []record.Enum16{enum16(1), enum16(2)},
				},
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "enum16",
				expectedValues: []interface{}{enum16(1), enum16(2)},
			},
			{
				name: "2 IN (1, 3)",
				comparator: Enum16FieldComparator{
					Cmp:    where.InArray,
					Getter: enum16Getter,
					Value:  []record.Enum16{enum16(1), enum16(3)},
				},
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "enum16",
				expectedValues: []interface{}{enum16(1), enum16(3)},
			},
			{
				name: "2 ? 2",
				comparator: Enum16FieldComparator{
					Cmp:    0,
					Getter: enum16Getter,
					Value:  []record.Enum16{enum16(2)},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(enum16Getter.Field, 0),
				expectedCmp:    0,
				expectedField:  "enum16",
				expectedValues: []interface{}{enum16(2)},
			},
		})
	})

	t.Run("int", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "10 = 10",
				comparator: IntFieldComparator{
					Cmp:    where.EQ,
					Getter: intGetter,
					Value:  []int{10},
				},
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "int",
				expectedValues: []interface{}{10},
			},
			{
				name: "10 = 3",
				comparator: IntFieldComparator{
					Cmp:    where.EQ,
					Getter: intGetter,
					Value:  []int{3},
				},
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "int",
				expectedValues: []interface{}{3},
			},
			{
				name: "10 > 3",
				comparator: IntFieldComparator{
					Cmp:    where.GT,
					Getter: intGetter,
					Value:  []int{3},
				},
				expectedResult: true,
				expectedCmp:    where.GT,
				expectedField:  "int",
				expectedValues: []interface{}{3},
			},
			{
				name: "10 > 30",
				comparator: IntFieldComparator{
					Cmp:    where.GT,
					Getter: intGetter,
					Value:  []int{30},
				},
				expectedResult: false,
				expectedCmp:    where.GT,
				expectedField:  "int",
				expectedValues: []interface{}{30},
			},
			{
				name: "10 >= 3",
				comparator: IntFieldComparator{
					Cmp:    where.GE,
					Getter: intGetter,
					Value:  []int{3},
				},
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int",
				expectedValues: []interface{}{3},
			},
			{
				name: "10 >= 30",
				comparator: IntFieldComparator{
					Cmp:    where.GE,
					Getter: intGetter,
					Value:  []int{30},
				},
				expectedResult: false,
				expectedCmp:    where.GE,
				expectedField:  "int",
				expectedValues: []interface{}{30},
			},
			{
				name: "10 >= 10",
				comparator: IntFieldComparator{
					Cmp:    where.GE,
					Getter: intGetter,
					Value:  []int{10},
				},
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int",
				expectedValues: []interface{}{10},
			},
			{
				name: "10 < 3",
				comparator: IntFieldComparator{
					Cmp:    where.LT,
					Getter: intGetter,
					Value:  []int{3},
				},
				expectedResult: false,
				expectedCmp:    where.LT,
				expectedField:  "int",
				expectedValues: []interface{}{3},
			},
			{
				name: "10 < 30",
				comparator: IntFieldComparator{
					Cmp:    where.LT,
					Getter: intGetter,
					Value:  []int{30},
				},
				expectedResult: true,
				expectedCmp:    where.LT,
				expectedField:  "int",
				expectedValues: []interface{}{30},
			},
			{
				name: "10 <= 3",
				comparator: IntFieldComparator{
					Cmp:    where.LE,
					Getter: intGetter,
					Value:  []int{3},
				},
				expectedResult: false,
				expectedCmp:    where.LE,
				expectedField:  "int",
				expectedValues: []interface{}{3},
			},
			{
				name: "10 <= 30",
				comparator: IntFieldComparator{
					Cmp:    where.LE,
					Getter: intGetter,
					Value:  []int{30},
				},
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int",
				expectedValues: []interface{}{30},
			},
			{
				name: "10 <= 10",
				comparator: IntFieldComparator{
					Cmp:    where.LE,
					Getter: intGetter,
					Value:  []int{10},
				},
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int",
				expectedValues: []interface{}{10},
			},
			{
				name: "10 IN (1, 2, 10)",
				comparator: IntFieldComparator{
					Cmp:    where.InArray,
					Getter: intGetter,
					Value:  []int{1, 2, 10},
				},
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "int",
				expectedValues: []interface{}{1, 2, 10},
			},
			{
				name: "10 IN (1, 3)",
				comparator: IntFieldComparator{
					Cmp:    where.InArray,
					Getter: intGetter,
					Value:  []int{1, 3},
				},
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "int",
				expectedValues: []interface{}{1, 3},
			},
			{
				name: "10 ? 10",
				comparator: IntFieldComparator{
					Cmp:    0,
					Getter: intGetter,
					Value:  []int{10},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(intGetter.Field, 0),
				expectedCmp:    0,
				expectedField:  "int",
				expectedValues: []interface{}{10},
			},
		})
	})

	t.Run("int32", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "10 = 10",
				comparator: Int32FieldComparator{
					Cmp:    where.EQ,
					Getter: int32Getter,
					Value:  []int32{10},
				},
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "int32",
				expectedValues: []interface{}{int32(10)},
			},
			{
				name: "10 = 3",
				comparator: Int32FieldComparator{
					Cmp:    where.EQ,
					Getter: int32Getter,
					Value:  []int32{3},
				},
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "int32",
				expectedValues: []interface{}{int32(3)},
			},
			{
				name: "10 > 3",
				comparator: Int32FieldComparator{
					Cmp:    where.GT,
					Getter: int32Getter,
					Value:  []int32{3},
				},
				expectedResult: true,
				expectedCmp:    where.GT,
				expectedField:  "int32",
				expectedValues: []interface{}{int32(3)},
			},
			{
				name: "10 > 30",
				comparator: Int32FieldComparator{
					Cmp:    where.GT,
					Getter: int32Getter,
					Value:  []int32{30},
				},
				expectedResult: false,
				expectedCmp:    where.GT,
				expectedField:  "int32",
				expectedValues: []interface{}{int32(30)},
			},
			{
				name: "10 >= 3",
				comparator: Int32FieldComparator{
					Cmp:    where.GE,
					Getter: int32Getter,
					Value:  []int32{3},
				},
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int32",
				expectedValues: []interface{}{int32(3)},
			},
			{
				name: "10 >= 30",
				comparator: Int32FieldComparator{
					Cmp:    where.GE,
					Getter: int32Getter,
					Value:  []int32{30},
				},
				expectedResult: false,
				expectedCmp:    where.GE,
				expectedField:  "int32",
				expectedValues: []interface{}{int32(30)},
			},
			{
				name: "10 >= 10",
				comparator: Int32FieldComparator{
					Cmp:    where.GE,
					Getter: int32Getter,
					Value:  []int32{10},
				},
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int32",
				expectedValues: []interface{}{int32(10)},
			},
			{
				name: "10 < 3",
				comparator: Int32FieldComparator{
					Cmp:    where.LT,
					Getter: int32Getter,
					Value:  []int32{3},
				},
				expectedResult: false,
				expectedCmp:    where.LT,
				expectedField:  "int32",
				expectedValues: []interface{}{int32(3)},
			},
			{
				name: "10 < 30",
				comparator: Int32FieldComparator{
					Cmp:    where.LT,
					Getter: int32Getter,
					Value:  []int32{30},
				},
				expectedResult: true,
				expectedCmp:    where.LT,
				expectedField:  "int32",
				expectedValues: []interface{}{int32(30)},
			},
			{
				name: "10 <= 3",
				comparator: Int32FieldComparator{
					Cmp:    where.LE,
					Getter: int32Getter,
					Value:  []int32{3},
				},
				expectedResult: false,
				expectedCmp:    where.LE,
				expectedField:  "int32",
				expectedValues: []interface{}{int32(3)},
			},
			{
				name: "10 <= 30",
				comparator: Int32FieldComparator{
					Cmp:    where.LE,
					Getter: int32Getter,
					Value:  []int32{30},
				},
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int32",
				expectedValues: []interface{}{int32(30)},
			},
			{
				name: "10 <= 10",
				comparator: Int32FieldComparator{
					Cmp:    where.LE,
					Getter: int32Getter,
					Value:  []int32{10},
				},
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int32",
				expectedValues: []interface{}{int32(10)},
			},
			{
				name: "10 IN (1, 2, 10)",
				comparator: Int32FieldComparator{
					Cmp:    where.InArray,
					Getter: int32Getter,
					Value:  []int32{1, 2, 10},
				},
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "int32",
				expectedValues: []interface{}{int32(1), int32(2), int32(10)},
			},
			{
				name: "10 IN (1, 3)",
				comparator: Int32FieldComparator{
					Cmp:    where.InArray,
					Getter: int32Getter,
					Value:  []int32{1, 3},
				},
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "int32",
				expectedValues: []interface{}{int32(1), int32(3)},
			},
			{
				name: "10 ? 10",
				comparator: Int32FieldComparator{
					Cmp:    0,
					Getter: int32Getter,
					Value:  []int32{10},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(int32Getter.Field, 0),
				expectedCmp:    0,
				expectedField:  "int32",
				expectedValues: []interface{}{int32(10)},
			},
		})
	})

	t.Run("int64", func(t *testing.T) {
		checkTestCases(t, []testCase{
			{
				name: "10 = 10",
				comparator: Int64FieldComparator{
					Cmp:    where.EQ,
					Getter: int64Getter,
					Value:  []int64{10},
				},
				expectedResult: true,
				expectedCmp:    where.EQ,
				expectedField:  "int64",
				expectedValues: []interface{}{int64(10)},
			},
			{
				name: "10 = 3",
				comparator: Int64FieldComparator{
					Cmp:    where.EQ,
					Getter: int64Getter,
					Value:  []int64{3},
				},
				expectedResult: false,
				expectedCmp:    where.EQ,
				expectedField:  "int64",
				expectedValues: []interface{}{int64(3)},
			},
			{
				name: "10 > 3",
				comparator: Int64FieldComparator{
					Cmp:    where.GT,
					Getter: int64Getter,
					Value:  []int64{3},
				},
				expectedResult: true,
				expectedCmp:    where.GT,
				expectedField:  "int64",
				expectedValues: []interface{}{int64(3)},
			},
			{
				name: "10 > 30",
				comparator: Int64FieldComparator{
					Cmp:    where.GT,
					Getter: int64Getter,
					Value:  []int64{30},
				},
				expectedResult: false,
				expectedCmp:    where.GT,
				expectedField:  "int64",
				expectedValues: []interface{}{int64(30)},
			},
			{
				name: "10 >= 3",
				comparator: Int64FieldComparator{
					Cmp:    where.GE,
					Getter: int64Getter,
					Value:  []int64{3},
				},
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int64",
				expectedValues: []interface{}{int64(3)},
			},
			{
				name: "10 >= 30",
				comparator: Int64FieldComparator{
					Cmp:    where.GE,
					Getter: int64Getter,
					Value:  []int64{30},
				},
				expectedResult: false,
				expectedCmp:    where.GE,
				expectedField:  "int64",
				expectedValues: []interface{}{int64(30)},
			},
			{
				name: "10 >= 10",
				comparator: Int64FieldComparator{
					Cmp:    where.GE,
					Getter: int64Getter,
					Value:  []int64{10},
				},
				expectedResult: true,
				expectedCmp:    where.GE,
				expectedField:  "int64",
				expectedValues: []interface{}{int64(10)},
			},
			{
				name: "10 < 3",
				comparator: Int64FieldComparator{
					Cmp:    where.LT,
					Getter: int64Getter,
					Value:  []int64{3},
				},
				expectedResult: false,
				expectedCmp:    where.LT,
				expectedField:  "int64",
				expectedValues: []interface{}{int64(3)},
			},
			{
				name: "10 < 30",
				comparator: Int64FieldComparator{
					Cmp:    where.LT,
					Getter: int64Getter,
					Value:  []int64{30},
				},
				expectedResult: true,
				expectedCmp:    where.LT,
				expectedField:  "int64",
				expectedValues: []interface{}{int64(30)},
			},
			{
				name: "10 <= 3",
				comparator: Int64FieldComparator{
					Cmp:    where.LE,
					Getter: int64Getter,
					Value:  []int64{3},
				},
				expectedResult: false,
				expectedCmp:    where.LE,
				expectedField:  "int64",
				expectedValues: []interface{}{int64(3)},
			},
			{
				name: "10 <= 30",
				comparator: Int64FieldComparator{
					Cmp:    where.LE,
					Getter: int64Getter,
					Value:  []int64{30},
				},
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int64",
				expectedValues: []interface{}{int64(30)},
			},
			{
				name: "10 <= 10",
				comparator: Int64FieldComparator{
					Cmp:    where.LE,
					Getter: int64Getter,
					Value:  []int64{10},
				},
				expectedResult: true,
				expectedCmp:    where.LE,
				expectedField:  "int64",
				expectedValues: []interface{}{int64(10)},
			},
			{
				name: "10 IN (1, 2, 10)",
				comparator: Int64FieldComparator{
					Cmp:    where.InArray,
					Getter: int64Getter,
					Value:  []int64{1, 2, 10},
				},
				expectedResult: true,
				expectedCmp:    where.InArray,
				expectedField:  "int64",
				expectedValues: []interface{}{int64(1), int64(2), int64(10)},
			},
			{
				name: "10 IN (1, 3)",
				comparator: Int64FieldComparator{
					Cmp:    where.InArray,
					Getter: int64Getter,
					Value:  []int64{1, 3},
				},
				expectedResult: false,
				expectedCmp:    where.InArray,
				expectedField:  "int64",
				expectedValues: []interface{}{int64(1), int64(3)},
			},
			{
				name: "10 ? 10",
				comparator: Int64FieldComparator{
					Cmp:    0,
					Getter: int64Getter,
					Value:  []int64{10},
				},
				expectedResult: false,
				expectedError:  NewNotImplementComparatorError(int64Getter.Field, 0),
				expectedCmp:    0,
				expectedField:  "int64",
				expectedValues: []interface{}{int64(10)},
			},
		})
	})
}
