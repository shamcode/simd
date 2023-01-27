package set

import (
	"reflect"
	"strconv"
	"sync/atomic"
	"unsafe"
)

// Based on https://github.com/cornelk/hashmap

// defaultSize is the default size for a set.
const defaultSize = 8

// maxFillRate is the maximum fill rate for the slice before a resize will happen.
const maxFillRate = 50

type Set struct {
	linkedList *List
	store      unsafe.Pointer
	resizing   int64
}

func NewSet() *Set {
	return NewSetSized(defaultSize)
}

func NewSetSized(size uintptr) *Set {
	m := &Set{}
	m.allocate(size)
	return m
}

func (m *Set) Count() int {
	return m.linkedList.Len()
}

func (m *Set) Add(key int64) {
	hash := xxHashQword(key)
	var (
		existed, inserted bool
		element           *ListElement
	)

	for {
		store := (*indexStore)(atomic.LoadPointer(&m.store))
		searchStart := store.item(hash)

		if !inserted { // if retrying after insert during grow, do not add to list again
			element, existed, inserted = m.linkedList.Add(searchStart, hash, key)
			if existed {
				return
			}
			if !inserted {
				continue // a concurrent add did interfere, try again
			}
		}

		count := store.addItem(element)
		currentStore := (*indexStore)(atomic.LoadPointer(&m.store))
		if store != currentStore { // retry insert in case of insert during grow
			continue
		}

		if m.isResizeNeeded(store, count) && atomic.CompareAndSwapInt64(&m.resizing, 0, 1) {
			go m.grow(0, true)
		}
		return
	}
}

func (m *Set) Delete(key int64) {
	hash := xxHashQword(key)
	store := (*indexStore)(atomic.LoadPointer(&m.store))
	element := store.item(hash)
	for ; element != nil; element = element.Next() {
		if element.keyHash == hash && element.key == key {
			m.deleteElement(element)
			m.linkedList.Delete(element)
			return
		}

		if element.keyHash > hash {
			return
		}
	}
}

func (m *Set) Iterate(f func(key int64)) {
	for item := m.linkedList.First(); item != nil; item = item.Next() {
		f(item.key)
	}
}

func (m *Set) allocate(newSize uintptr) {
	m.linkedList = NewList()
	if atomic.CompareAndSwapInt64(&m.resizing, 0, 1) {
		m.grow(newSize, false)
	}
}

func (m *Set) isResizeNeeded(store *indexStore, count uintptr) bool {
	l := uintptr(len(store.index)) // l can't be 0 as it gets initialized in New()
	fillRate := (count * 100) / l
	return fillRate > maxFillRate
}

func (m *Set) deleteElement(element *ListElement) {
	for {
		store := (*indexStore)(atomic.LoadPointer(&m.store))
		index := element.keyHash >> store.keyShifts
		ptr := (*unsafe.Pointer)(unsafe.Pointer(uintptr(store.array) + index*intSizeBytes))

		next := element.Next()
		if next != nil && element.keyHash>>store.keyShifts != index {
			next = nil // do not Set index to next item if it's not the same slice index
		}
		atomic.CompareAndSwapPointer(ptr, unsafe.Pointer(element), unsafe.Pointer(next))

		currentStore := (*indexStore)(atomic.LoadPointer(&m.store))
		if store == currentStore { // check that no resize happened
			break
		}
	}
}

func (m *Set) grow(newSize uintptr, loop bool) {
	defer atomic.CompareAndSwapInt64(&m.resizing, 1, 0)

	for {
		currentStore := (*indexStore)(atomic.LoadPointer(&m.store))
		if newSize == 0 {
			newSize = uintptr(len(currentStore.index)) << 1
		} else {
			newSize = roundUpPower2(newSize)
		}

		index := make([]*ListElement, newSize)
		header := (*reflect.SliceHeader)(unsafe.Pointer(&index))

		newStore := &indexStore{
			keyShifts: strconv.IntSize - log2(newSize),
			array:     unsafe.Pointer(header.Data), // use address of slice data storage
			index:     index,
		}

		m.fillIndexItems(newStore) // initialize new index slice with longer keys

		atomic.StorePointer(&m.store, unsafe.Pointer(newStore))

		m.fillIndexItems(newStore) // make sure that the new index is up-to-date with the current state of the linked list

		if !loop {
			return
		}

		// check if a new resize needs to be done already
		count := uintptr(m.Count())
		if !m.isResizeNeeded(newStore, count) {
			return
		}
		newSize = 0 // 0 means double the current size
	}
}

func (m *Set) fillIndexItems(store *indexStore) {
	first := m.linkedList.First()
	item := first
	lastIndex := uintptr(0)

	for item != nil {
		index := item.keyHash >> store.keyShifts
		if item == first || index != lastIndex { // store item with smallest hash key for every index
			store.addItem(item)
			lastIndex = index
		}
		item = item.Next()
	}
}

// roundUpPower2 rounds a number to the next power of 2.
func roundUpPower2(i uintptr) uintptr {
	i--
	i |= i >> 1
	i |= i >> 2
	i |= i >> 4
	i |= i >> 8
	i |= i >> 16
	i |= i >> 32
	i++
	return i
}

// log2 computes the binary logarithm of x, rounded up to the next integer.
func log2(i uintptr) uintptr {
	var n, p uintptr
	for p = 1; p < i; p += p {
		n++
	}
	return n
}
