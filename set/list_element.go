package set

import (
	"sync/atomic"
	"unsafe"
)

type ListElement struct {
	keyHash uintptr
	key     int64
	deleted int64
	next    unsafe.Pointer
}

func (e *ListElement) Next() *ListElement {
	for next := (*ListElement)(atomic.LoadPointer(&e.next)); next != nil; {
		if atomic.LoadInt64(&next.deleted) == 0 {
			return next
		}

		following := next.Next()
		if atomic.CompareAndSwapPointer(&e.next, unsafe.Pointer(next), unsafe.Pointer(following)) {
			next = following
		} else {
			next = next.Next()
		}
	}

	return nil
}
