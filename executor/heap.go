//nolint:varnamelen
package executor

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
)

// binaryHeap is a interface{} free "container/heap" with small optimization (https://en.wikipedia.org/wiki/Binary_heap)
type binaryHeap[R record.Record] struct {
	records []R
	sorting []sort.ByWithOrder[R]
}

func (h *binaryHeap[R]) less(i, j int) int8 {
	a := h.records[i]
	b := h.records[j]
	for _, by := range h.sorting {
		if by.Less(a, b) {
			return -1
		} else if by.Less(b, a) {
			return 1
		}
	}
	return 0
}
func (h *binaryHeap[R]) swap(i, j int) { h.records[i], h.records[j] = h.records[j], h.records[i] }

func (h *binaryHeap[R]) Push(item R) {
	h.records = append(h.records, item)
	h.up(len(h.records) - 1)
}

// Remove removes and returns the element at index i from the binaryHeap.
// The complexity is O(log n) where n = h.Len().
func (h *binaryHeap[R]) Remove(i int) R {
	n := len(h.records) - 1
	if n != i {
		h.swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	old := h.records
	x := old[n]
	h.records = old[0:n]
	return x
}

func (h *binaryHeap[R]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || h.less(j, i) > 0 {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *binaryHeap[R]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(j1, j2) > 0 {
			j = j2 // = 2*i + 2  // right child
		}
		if h.less(i, j) < 0 {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > i0
}

func newHeap[R record.Record](sorting []sort.ByWithOrder[R]) *binaryHeap[R] {
	return &binaryHeap[R]{ //nolint:exhaustruct
		sorting: sorting,
	}
}
