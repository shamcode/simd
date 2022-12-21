package namespace

import (
	"context"
	"fmt"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/fields"
	"regexp"
)

var (
	_ Query = (*executor)(nil)
)

type executor struct {
	limitItems   int
	startOffset  int
	withLimit    bool
	withNot      bool
	isOr         bool
	conditionSet bool
	bracketLevel int
	where        where.Conditions
	sortBy       []sort.By
	storage      Namespace
	onIteration  *func(item record.Record)
}

func (q *executor) AddWhere(cmp where.FieldComparator) Query {
	q.where = append(q.where, &where.Condition{
		WithNot:      q.withNot,
		IsOr:         q.isOr,
		BracketLevel: 1 + q.bracketLevel,
		Cmp:          cmp,
	})
	q.withNot = false
	q.isOr = false
	q.conditionSet = true
	return q
}

func (q *executor) MakeCopy() Query {
	cpy := &executor{
		limitItems:   q.limitItems,
		startOffset:  q.startOffset,
		withLimit:    q.withLimit,
		withNot:      q.withNot,
		isOr:         q.isOr,
		bracketLevel: q.bracketLevel,
		where:        make(where.Conditions, len(q.where), cap(q.where)),
		sortBy:       make([]sort.By, len(q.sortBy), cap(q.sortBy)),
		storage:      q.storage,
		onIteration:  q.onIteration,
	}
	for i, item := range q.where {
		cpy.where[i] = item
	}
	for i, item := range q.sortBy {
		cpy.sortBy[i] = item
	}
	return cpy
}

func (q *executor) OnIteration(cb func(item record.Record)) Query {
	q.onIteration = &cb
	return q
}

func (q *executor) Limit(limitItems int) Query {
	q.limitItems = limitItems
	q.withLimit = true
	return q
}

func (q *executor) Offset(startOffset int) Query {
	q.startOffset = startOffset
	return q
}

func (q *executor) Or() Query {
	q.isOr = true
	if !q.conditionSet {
		panic(".Or() before any condition not supported, add any condition before .Or()")
	}
	return q
}

func (q *executor) Not() Query {
	q.withNot = !q.withNot
	return q
}

func (q *executor) OpenBracket() Query {
	if q.withNot {
		panic(".Not().OpenBracket() not supported")
	}
	q.conditionSet = false
	q.bracketLevel += 1
	return q
}

func (q *executor) CloseBracket() Query {
	q.bracketLevel -= 1
	if -1 == q.bracketLevel {
		panic("close bracket without open")
	}
	q.conditionSet = true
	return q
}

func (q *executor) Where(getter *record.InterfaceGetter, condition where.ComparatorType, value ...interface{}) Query {
	return q.AddWhere(&fields.InterfaceFieldComparator{
		BaseFieldComparator: fields.BaseFieldComparator{
			Cmp: condition,
		},
		Getter: getter,
		Value:  value,
	})
}

func (q *executor) WhereInt(getter *record.IntGetter, condition where.ComparatorType, value ...int) Query {
	return q.AddWhere(&fields.IntFieldComparator{
		BaseFieldComparator: fields.BaseFieldComparator{
			Cmp: condition,
		},
		Getter: getter,
		Value:  value,
	})
}

func (q *executor) WhereInt32(getter *record.Int32Getter, condition where.ComparatorType, value ...int32) Query {
	return q.AddWhere(&fields.Int32FieldComparator{
		BaseFieldComparator: fields.BaseFieldComparator{
			Cmp: condition,
		},
		Getter: getter,
		Value:  value,
	})
}

func (q *executor) WhereInt64(getter *record.Int64Getter, condition where.ComparatorType, value ...int64) Query {
	return q.AddWhere(&fields.Int64FieldComparator{
		BaseFieldComparator: fields.BaseFieldComparator{
			Cmp: condition,
		},
		Getter: getter,
		Value:  value,
	})
}

func (q *executor) WhereString(getter *record.StringGetter, condition where.ComparatorType, value ...string) Query {
	return q.AddWhere(&fields.StringFieldComparator{
		BaseFieldComparator: fields.BaseFieldComparator{
			Cmp: condition,
		},
		Getter: getter,
		Value:  value,
	})
}

func (q *executor) WhereStringRegexp(getter *record.StringGetter, value *regexp.Regexp) Query {
	return q.AddWhere(&fields.StringFieldRegexpComparator{
		BaseFieldComparator: fields.BaseFieldComparator{
			Cmp: where.Regexp,
		},
		Getter: getter,
		Value:  value,
	})
}

func (q *executor) WhereBool(getter *record.BoolGetter, condition where.ComparatorType, value ...bool) Query {
	return q.AddWhere(&fields.BoolFieldComparator{
		BaseFieldComparator: fields.BaseFieldComparator{
			Cmp: condition,
		},
		Getter: getter,
		Value:  value,
	})
}

func (q *executor) WhereEnum8(getter *record.Enum8Getter, condition where.ComparatorType, value ...record.Enum8) Query {
	return q.AddWhere(&fields.Enum8FieldComparator{
		BaseFieldComparator: fields.BaseFieldComparator{
			Cmp: condition,
		},
		Getter: getter,
		Value:  value,
	})
}

func (q *executor) WhereEnum16(getter *record.Enum16Getter, condition where.ComparatorType, value ...record.Enum16) Query {
	return q.AddWhere(&fields.Enum16FieldComparator{
		BaseFieldComparator: fields.BaseFieldComparator{
			Cmp: condition,
		},
		Getter: getter,
		Value:  value,
	})
}

func (q *executor) WhereMap(getter *record.MapGetter, condition where.ComparatorType, value ...interface{}) Query {
	return q.AddWhere(&fields.MapFieldComparator{
		BaseFieldComparator: fields.BaseFieldComparator{
			Cmp: condition,
		},
		Getter: getter,
		Value:  value,
	})
}

func (q *executor) WhereSet(getter *record.SetGetter, condition where.ComparatorType, value ...interface{}) Query {
	return q.AddWhere(&fields.SetFieldComparator{
		BaseFieldComparator: fields.BaseFieldComparator{
			Cmp: condition,
		},
		Getter: getter,
		Value:  value,
	})
}

func (q *executor) WhereStringsSet(getter *record.StringsSetGetter, condition where.ComparatorType, value ...string) Query {
	return q.AddWhere(&fields.StringsSetFieldComparator{
		BaseFieldComparator: fields.BaseFieldComparator{
			Cmp: condition,
		},
		Getter: getter,
		Value:  value,
	})
}

func (q *executor) Sort(sortBy sort.By) Query {
	q.sortBy = append(q.sortBy, sortBy)
	return q
}

func (q *executor) FetchTotal(ctx context.Context) (int, error) {
	_, total, err := q.exec(ctx, true)
	if nil != err {
		return 0, wrapErrors(ExecuteQueryErr, err)
	}
	return total, nil
}

func (q *executor) FetchAll(ctx context.Context) (Iterator, error) {
	iter, _, err := q.exec(ctx, false)
	if nil != err {
		return nil, wrapErrors(ExecuteQueryErr, err)
	}
	return iter, nil
}

func (q *executor) FetchAllAndTotal(ctx context.Context) (Iterator, int, error) {
	iter, total, err := q.exec(ctx, false)
	if nil != err {
		return nil, 0, wrapErrors(ExecuteQueryErr, err)
	}
	return iter, total, nil
}

func (q *executor) exec(ctx context.Context, onlyTotal bool) (Iterator, int, error) {
	if q.bracketLevel > 0 {
		return nil, 0, fmt.Errorf("invalid bracet balance: has %d not closed bracket", q.bracketLevel)
	}

	total := 0
	items := newHeap(q.sortBy)
	for _, item := range q.storage.Select(q.where) {
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
			if !q.where.Check(item) {
				continue
			}
			if nil != q.onIteration {
				(*q.onIteration)(item)
			}
			total += 1
			if !onlyTotal {
				items.Push(item)
			}
		}
	}

	if onlyTotal {
		return nil, total, nil
	}

	var last int
	var size int
	itemsCount := total
	if q.withLimit {
		last = q.startOffset + q.limitItems
		if last > itemsCount {
			last = itemsCount
		}
		size = q.limitItems
		if size > itemsCount {
			size = itemsCount
		}
	} else {
		last = itemsCount
		size = itemsCount
	}

	return newHeapIterator(items, q.startOffset, last, size), total, nil
}

func Create(storage Namespace) Query {
	return &executor{
		storage: storage,
	}
}
