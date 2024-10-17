package executor

import (
	"context"

	"github.com/shamcode/simd/record"
)

type Iterator interface {
	Next(ctx context.Context) bool
	Item() record.Record
	Size() int
	Err() error
}

var _ Iterator = (*heapIterator)(nil)

type heapIterator struct {
	from      int
	index     int
	max       int
	size      int
	heap      *binaryHeap
	lastError error
}

func (i *heapIterator) Next(ctx context.Context) bool {
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

func (i *heapIterator) Item() record.Record {
	return i.heap.Remove(i.from)
}

func (i *heapIterator) Err() error {
	return i.lastError
}

func (i *heapIterator) Size() int {
	return i.size
}

func newHeapIterator(heap *binaryHeap, from, to, size int) Iterator {
	return &heapIterator{
		from:  from,
		index: from,
		max:   to,
		size:  size,
		heap:  heap,
	}
}
