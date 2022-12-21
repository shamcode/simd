package namespace

import (
	"context"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"regexp"
)

type Namespace interface {
	Get(id int64) record.Record
	Insert(item record.Record) error
	Delete(id int64) error
	Upsert(item record.Record) error
	Query() Query
	Select(conditions where.Conditions) []record.Record
}

type Query interface {
	MakeCopy() Query

	Limit(limitItems int) Query
	Offset(startOffset int) Query

	Not() Query
	Or() Query

	OpenBracket() Query
	CloseBracket() Query

	AddWhere(cmp where.FieldComparator) Query

	Where(getter *record.InterfaceGetter, condition where.ComparatorType, values ...interface{}) Query
	WhereInt(getter *record.IntGetter, condition where.ComparatorType, values ...int) Query
	WhereInt32(getter *record.Int32Getter, condition where.ComparatorType, values ...int32) Query
	WhereInt64(getter *record.Int64Getter, condition where.ComparatorType, values ...int64) Query
	WhereString(getter *record.StringGetter, condition where.ComparatorType, values ...string) Query
	WhereStringRegexp(getter *record.StringGetter, value *regexp.Regexp) Query
	WhereBool(getter *record.BoolGetter, condition where.ComparatorType, values ...bool) Query
	WhereEnum8(getter *record.Enum8Getter, condition where.ComparatorType, values ...record.Enum8) Query
	WhereEnum16(getter *record.Enum16Getter, condition where.ComparatorType, values ...record.Enum16) Query
	WhereMap(getter *record.MapGetter, condition where.ComparatorType, values ...interface{}) Query
	WhereSet(getter *record.SetGetter, condition where.ComparatorType, values ...interface{}) Query
	WhereStringsSet(getter *record.StringsSetGetter, condition where.ComparatorType, values ...string) Query

	Sort(by sort.By) Query

	OnIteration(cb func(item record.Record)) Query

	FetchTotal(ctx context.Context) (int, error)
	FetchAll(ctx context.Context) (Iterator, error)
	FetchAllAndTotal(ctx context.Context) (Iterator, int, error)
}

type Iterator interface {
	Next(ctx context.Context) bool
	Item() record.Record
	Size() int
	Err() error
}
