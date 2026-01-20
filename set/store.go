package set

import (
	"strconv"
	"sync/atomic"
	"unsafe"
)

const intSizeBytes = strconv.IntSize >> 3

type indexStore struct {
	keyShifts uintptr        // Pointer size - log2 of array size, to be used as index in the data array
	count     uintptr        // count of filled elements in the slice
	array     unsafe.Pointer // pointer to slice data array
	index     []*ListElement // storage for the slice for the garbage collector to not clean it up
}

func (s *indexStore) item(hashedKey uintptr) *ListElement {
	index := hashedKey >> s.keyShifts
	ptr := (*unsafe.Pointer)(unsafe.Add(s.array, index*intSizeBytes))
	item := (*ListElement)(atomic.LoadPointer(ptr))

	return item
}

func (s *indexStore) addItem(item *ListElement) uintptr {
	index := item.keyHash >> s.keyShifts

	ptr := (*unsafe.Pointer)(unsafe.Add(s.array, index*intSizeBytes))
	for { // loop until the smallest key hash is in the index
		element := (*ListElement)(atomic.LoadPointer(ptr)) // get the current item in the index
		if element == nil {                                // no item yet at this index
			if atomic.CompareAndSwapPointer(ptr, nil, unsafe.Pointer(item)) {
				return atomic.LoadUintptr(&s.count)
			}

			continue // a new item was inserted concurrently, retry
		}

		if item.keyHash < element.keyHash {
			// the new item is the smallest for this index?
			if !atomic.CompareAndSwapPointer(ptr, unsafe.Pointer(element), unsafe.Pointer(item)) {
				continue // a new item was inserted concurrently, retry
			}
		}

		return 0
	}
}
