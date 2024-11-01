package query

import (
	"testing"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

type builderOptionFunction[R record.Record] func(b BuilderGeneric[R])

func (fn builderOptionFunction[R]) Apply(b any) {
	fn(b.(BuilderGeneric[R]))
}

type whereStructStack[R record.Record] struct {
	cmp where.FieldComparator[R]
}

func (w whereStructStack[R]) Apply(b any) {
	b.(BuilderGeneric[R]).AddWhere(w.cmp)
}

type whereStructHeap[R record.Record] struct {
	cmp where.FieldComparator[R]
}

func (w *whereStructHeap[R]) Apply(b any) {
	b.(BuilderGeneric[R]).AddWhere(w.cmp)
}

func Benchmark_OptionsStructAndFunction(b *testing.B) {
	_id := record.NewIDGetter[record.Record]()
	var whereInt64Fn = func(
		getter record.GetterInterface[record.Record, int64],
		condition where.ComparatorType,
		value ...int64,
	) BuilderOption {
		return builderOptionFunction[record.Record](func(b BuilderGeneric[record.Record]) {
			b.AddWhere(comparators.ComparableFieldComparator[record.Record, int64]{
				EqualComparator: comparators.EqualComparator[record.Record, int64]{
					Cmp:    condition,
					Getter: getter,
					Value:  value,
				},
			})
		})
	}

	var whereInt64StructStack = func(
		getter record.GetterInterface[record.Record, int64],
		condition where.ComparatorType,
		value ...int64,
	) BuilderOption {
		return whereStructStack[record.Record]{
			cmp: comparators.ComparableFieldComparator[record.Record, int64]{
				EqualComparator: comparators.EqualComparator[record.Record, int64]{
					Cmp:    condition,
					Getter: getter,
					Value:  value,
				},
			},
		}
	}

	var whereInt64StructHeap = func(
		getter record.GetterInterface[record.Record, int64],
		condition where.ComparatorType,
		value ...int64,
	) BuilderOption {
		return &whereStructHeap[record.Record]{
			cmp: comparators.ComparableFieldComparator[record.Record, int64]{
				EqualComparator: comparators.EqualComparator[record.Record, int64]{
					Cmp:    condition,
					Getter: getter,
					Value:  value,
				},
			},
		}
	}

	b.Run("struct stack", func(b *testing.B) {
		qb := NewBuilder[record.Record]()
		for i := 0; i < b.N; i++ {
			qb.Append(whereInt64StructStack(_id, where.EQ, 1))
		}
	})

	b.Run("struct heap", func(b *testing.B) {
		qb := NewBuilder[record.Record]()
		for i := 0; i < b.N; i++ {
			qb.Append(whereInt64StructHeap(_id, where.EQ, 1))
		}
	})

	b.Run("func", func(b *testing.B) {
		qb := NewBuilder[record.Record]()
		for i := 0; i < b.N; i++ {
			qb.Append(whereInt64Fn(_id, where.EQ, 1))
		}
	})
}
