package query

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
	"regexp"
)

type LimitOption int

func (o LimitOption) Apply(b Builder) { b.Limit(int(o)) }

func Limit(limitItems int) BuilderOption {
	return LimitOption(limitItems)
}

type OffsetOption int

func (o OffsetOption) Apply(b Builder) { b.Offset(int(o)) }
func Offset(startOffset int) BuilderOption {
	return OffsetOption(startOffset)
}

type OrOption struct{}

func (_ OrOption) Apply(b Builder) { b.Or() }
func Or() BuilderOption {
	return OrOption{}
}

type NotOption struct{}

func (_ NotOption) Apply(b Builder) { b.Not() }
func Not() BuilderOption {
	return NotOption{}
}

type OpenBracketOption struct{}

func (_ OpenBracketOption) Apply(b Builder) { b.OpenBracket() }
func OpenBracket() BuilderOption {
	return OpenBracketOption{}
}

type CloseBracketOption struct{}

func (_ CloseBracketOption) Apply(b Builder) { b.CloseBracket() }
func CloseBracket() BuilderOption {
	return CloseBracketOption{}
}

type SortOption struct {
	by sort.ByWithOrder
}

func (o SortOption) Apply(b Builder) { b.Sort(o.by) }
func Sort(by sort.ByWithOrder) BuilderOption {
	return SortOption{by: by}
}

type AddWhereOption struct {
	Cmp where.FieldComparator
}

func (o AddWhereOption) Apply(b Builder) { b.AddWhere(o.Cmp) }

func Where(getter record.InterfaceGetter, condition where.ComparatorType, values ...interface{}) BuilderOption {
	return AddWhereOption{
		Cmp: comparators.InterfaceFieldComparator{
			Cmp:    condition,
			Getter: getter,
			Value:  values,
		},
	}
}

func WhereInt(getter record.IntGetter, condition where.ComparatorType, value ...int) BuilderOption {
	return AddWhereOption{
		Cmp: comparators.IntFieldComparator{
			Cmp:    condition,
			Getter: getter,
			Value:  value,
		},
	}
}

func WhereInt32(getter record.Int32Getter, condition where.ComparatorType, value ...int32) BuilderOption {
	return AddWhereOption{
		Cmp: comparators.Int32FieldComparator{
			Cmp:    condition,
			Getter: getter,
			Value:  value,
		},
	}
}

func WhereInt64(getter record.Int64Getter, condition where.ComparatorType, value ...int64) BuilderOption {
	return AddWhereOption{
		Cmp: comparators.Int64FieldComparator{
			Cmp:    condition,
			Getter: getter,
			Value:  value,
		},
	}
}

func WhereString(getter record.StringGetter, condition where.ComparatorType, value ...string) BuilderOption {
	return AddWhereOption{
		Cmp: comparators.StringFieldComparator{
			Cmp:    condition,
			Getter: getter,
			Value:  value,
		},
	}
}

func WhereStringRegexp(getter record.StringGetter, value *regexp.Regexp) BuilderOption {
	return AddWhereOption{
		Cmp: comparators.StringFieldRegexpComparator{
			Cmp:    where.Regexp,
			Getter: getter,
			Value:  value,
		},
	}
}

func WhereBool(getter record.BoolGetter, condition where.ComparatorType, value ...bool) BuilderOption {
	return AddWhereOption{
		Cmp: comparators.BoolFieldComparator{
			Cmp:    condition,
			Getter: getter,
			Value:  value,
		},
	}
}

func WhereEnum8(getter record.Enum8Getter, condition where.ComparatorType, value ...record.Enum8) BuilderOption {
	return AddWhereOption{
		Cmp: comparators.Enum8FieldComparator{
			Cmp:    condition,
			Getter: getter,
			Value:  value,
		},
	}
}

func WhereEnum16(getter record.Enum16Getter, condition where.ComparatorType, value ...record.Enum16) BuilderOption {
	return AddWhereOption{
		Cmp: comparators.Enum16FieldComparator{
			Cmp:    condition,
			Getter: getter,
			Value:  value,
		},
	}
}

func WhereMap(getter record.MapGetter, condition where.ComparatorType, value ...interface{}) BuilderOption {
	return AddWhereOption{
		Cmp: comparators.MapFieldComparator{
			Cmp:    condition,
			Getter: getter,
			Value:  value,
		},
	}
}

func WhereSet(getter record.SetGetter, condition where.ComparatorType, value ...interface{}) BuilderOption {
	return AddWhereOption{
		Cmp: comparators.SetFieldComparator{
			Cmp:    condition,
			Getter: getter,
			Value:  value,
		},
	}
}

type onIterationOption func(item record.Record)

func (o onIterationOption) Apply(b Builder) { b.OnIteration(o) }

func OnIteration(cb func(item record.Record)) BuilderOption {
	return onIterationOption(cb)
}
