package executor

import (
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/sort"
)

// binaryHeap is a interface{} free "container/heap" with small optimization (https://en.wikipedia.org/wiki/Binary_heap)
type binaryHeap struct {
	records []record.Record
	sorting []sort.By
}

func (h *binaryHeap) less(i, j int) int8 {
	a := h.records[i]
	b := h.records[j]
	for _, by := range h.sorting {
		if by.Equal(a, b) {
			continue
		} else if by.Less(a, b) {
			return -1
		} else {
			return 1
		}
	}
	return 0
}
func (h *binaryHeap) swap(i, j int) { h.records[i], h.records[j] = h.records[j], h.records[i] }

func (h *binaryHeap) Push(item record.Record) {
	h.records = append(h.records, item)
	h.up(len(h.records) - 1)
}

// Remove removes and returns the element at index i from the binaryHeap.
// The complexity is O(log n) where n = h.Len().
func (h *binaryHeap) Remove(i int) record.Record {
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

func (h *binaryHeap) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || h.less(j, i) > 0 {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *binaryHeap) down(i0, n int) bool {
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

func newHeap(sorting []sort.By) *binaryHeap {
	return &binaryHeap{
		sorting: sorting,
	}
}
