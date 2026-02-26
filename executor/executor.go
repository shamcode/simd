//nolint:varnamelen
package executor

import (
	"context"

	"github.com/shamcode/simd/record"

	"github.com/shamcode/simd/query"
)

type QueryExecutor[R record.Record] interface {
	FetchTotal(ctx context.Context, q query.Query[R]) (int, error)
	FetchAll(ctx context.Context, q query.Query[R]) (Iterator[R], error)
	FetchAllAndTotal(ctx context.Context, q query.Query[R]) (Iterator[R], int, error)
}

type executor[R record.Record] struct {
	selector Selector[R]
}

func (e *executor[R]) FetchTotal(ctx context.Context, q query.Query[R]) (int, error) {
	_, total, err := e.exec(ctx, q, true)
	return total, err
}

func (e *executor[R]) FetchAll(ctx context.Context, q query.Query[R]) (Iterator[R], error) {
	iter, _, err := e.exec(ctx, q, false)
	return iter, err
}

func (e *executor[R]) FetchAllAndTotal(ctx context.Context, q query.Query[R]) (Iterator[R], int, error) {
	return e.exec(ctx, q, false)
}

func (e *executor[R]) exec( //nolint:cyclop,funlen
	ctx context.Context,
	q query.Query[R],
	onlyTotal bool,
) (Iterator[R], int, error) {
	if err := q.Error(); err != nil {
		return nil, 0, NewValidateQueryError(err)
	}

	total := 0
	items := newHeap(q.Sorting())
	callback := q.OnIterationCallback()
	conditions := q.Conditions()

	itemsForCheck, err := e.selector.PreselectForExecutor(conditions)
	if err != nil {
		return nil, 0, NewExecuteQueryError(err)
	}

	for _, item := range itemsForCheck {
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
			res, err := conditions.Check(item)
			if err != nil {
				return nil, 0, NewExecuteQueryError(err)
			}

			if !res {
				continue
			}

			if nil != callback {
				(*callback)(item)
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

	var (
		last int
		size int
	)

	itemsCount := total

	if limit, withLimit := q.Limit(); withLimit {
		last = min(q.Offset()+limit, itemsCount)

		size = min(limit, itemsCount)
	} else {
		last = itemsCount
		size = itemsCount
	}

	return newHeapIterator(items, q.Offset(), last, size), total, nil
}

func CreateQueryExecutor[R record.Record](selector Selector[R]) QueryExecutor[R] {
	return &executor[R]{
		selector: selector,
	}
}
