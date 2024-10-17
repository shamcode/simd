package query

import (
	"testing"

	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

type builderOptionFunction func(b Builder)

func (fn builderOptionFunction) Apply(b Builder) {
	fn(b)
}

type whereStructStack struct {
	cmp where.FieldComparator
}

func (w whereStructStack) Apply(b Builder) {
	b.AddWhere(w.cmp)
}

type whereStructHeap struct {
	cmp where.FieldComparator
}

func (w *whereStructHeap) Apply(b Builder) {
	b.AddWhere(w.cmp)
}

func Benchmark_OptionsStructAndFunction(b *testing.B) {
	var whereInt64Fn = func(getter record.Int64Getter, condition where.ComparatorType, value ...int64) BuilderOption {
		return builderOptionFunction(func(b Builder) {
			b.AddWhere(comparators.Int64FieldComparator{
				Cmp:    condition,
				Getter: getter,
				Value:  value,
			})
		})
	}

	var whereInt64StructStack = func(getter record.Int64Getter, condition where.ComparatorType, value ...int64) BuilderOption {
		return whereStructStack{
			cmp: comparators.Int64FieldComparator{
				Cmp:    condition,
				Getter: getter,
				Value:  value,
			},
		}
	}

	var whereInt64StructHeap = func(getter record.Int64Getter, condition where.ComparatorType, value ...int64) BuilderOption {
		return &whereStructHeap{
			cmp: comparators.Int64FieldComparator{
				Cmp:    condition,
				Getter: getter,
				Value:  value,
			},
		}
	}

	b.Run("struct stack", func(b *testing.B) {
		qb := NewBuilder()
		for i := 0; i < b.N; i++ {
			qb.Append(whereInt64StructStack(record.ID, where.EQ, 1))
		}
	})

	b.Run("struct heap", func(b *testing.B) {
		qb := NewBuilder()
		for i := 0; i < b.N; i++ {
			qb.Append(whereInt64StructHeap(record.ID, where.EQ, 1))
		}
	})

	b.Run("func", func(b *testing.B) {
		qb := NewBuilder()
		for i := 0; i < b.N; i++ {
			qb.Append(whereInt64Fn(record.ID, where.EQ, 1))
		}
	})
}
