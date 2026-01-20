package set

import (
	"sync/atomic"
	"unsafe"
)

type List struct {
	count int64
	head  *ListElement
}

func NewList() *List {
	return &List{ //nolint:exhaustruct
		head: &ListElement{}, //nolint:exhaustruct
	}
}

func (list *List) Len() int {
	return int(atomic.LoadInt64(&list.count))
}

func (list *List) First() *ListElement {
	return list.head.Next()
}

func (list *List) Add( //nolint:nonamedreturns
	searchStart *ListElement,
	hash uintptr,
	key int64,
) (element *ListElement, existed bool, inserted bool) {
	left, found, right := list.search(searchStart, hash, key)
	if found != nil { // existing item found
		return found, true, false
	}

	element = &ListElement{ //nolint:exhaustruct
		keyHash: hash,
		key:     key,
	}

	return element, false, list.insertAt(element, left, right)
}

func (list *List) Delete(element *ListElement) {
	if !atomic.CompareAndSwapInt64(&element.deleted, 0, 1) {
		return
	}

	right := element.Next()
	atomic.CompareAndSwapPointer(&list.head.next, unsafe.Pointer(element), unsafe.Pointer(right))
	atomic.AddInt64(&list.count, -1)
}

func (list *List) search( //nolint:cyclop,nonamedreturns
	searchStart *ListElement,
	hash uintptr,
	key int64,
) (left, found, right *ListElement) {
	if searchStart != nil && hash < searchStart.keyHash {
		searchStart = nil
	}

	if searchStart == nil { // start search at head?
		left = list.head

		found = left.Next()
		if found == nil { // no items beside head?
			return nil, nil, nil
		}
	} else {
		found = searchStart
	}

	for {
		if hash == found.keyHash && key == found.key {
			return nil, found, nil
		}

		if hash < found.keyHash { // new item needs to be inserted before the found value
			if list.head == left {
				return nil, nil, found
			}

			return left, nil, found
		}

		// go to next element in sorted linked list
		left = found

		found = left.Next()
		if found == nil { // no more items on the right
			return left, nil, nil
		}
	}
}

func (list *List) insertAt(element, left, right *ListElement) bool {
	if left == nil {
		left = list.head
	}

	atomic.StorePointer(&element.next, unsafe.Pointer(right))

	if !atomic.CompareAndSwapPointer(&left.next, unsafe.Pointer(right), unsafe.Pointer(element)) {
		return false
	}

	atomic.AddInt64(&list.count, 1)

	return true
}
