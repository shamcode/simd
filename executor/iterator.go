package executor

import (
	"context"

	"github.com/shamcode/simd/record"
)

type Iterator[R record.Record] interface {
	Next(ctx context.Context) bool
	Item() R
	Size() int
	Err() error
}

type heapIterator[R record.Record] struct {
	from      int
	index     int
	max       int
	size      int
	heap      *binaryHeap[R]
	lastError error
}

func (i *heapIterator[R]) Next(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		i.lastError = ctx.Err()
		return false
	default:
		result := i.index < i.max
		if result {
			i.index += 1
		}
		return result
	}
}

func (i *heapIterator[R]) Item() R {
	return i.heap.Remove(i.from)
}

func (i *heapIterator[R]) Err() error {
	return i.lastError
}

func (i *heapIterator[R]) Size() int {
	return i.size
}

func newHeapIterator[R record.Record](heap *binaryHeap[R], from, to, size int) Iterator[R] {
	return &heapIterator[R]{ //nolint:exhaustruct
		from:  from,
		index: from,
		max:   to,
		size:  size,
		heap:  heap,
	}
}
