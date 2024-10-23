package storage

import "sync/atomic"

var _ UniqueIDStorage = (*uniqID)(nil)

type uniqID int64

func (u *uniqID) RLock()   {}
func (u *uniqID) RUnlock() {}
func (u *uniqID) Iterate(f func(id int64)) {
	if id := atomic.LoadInt64((*int64)(u)); id != 0 {
		f(id)
	}
}

func (u *uniqID) Count() int {
	if id := atomic.LoadInt64((*int64)(u)); id == 0 {
		return 0
	}
	return 1
}

func (u *uniqID) Add(id int64) {
	atomic.StoreInt64((*int64)(u), id)
}

func (u *uniqID) Delete(_ int64) {
	atomic.StoreInt64((*int64)(u), 0)
}

func (u *uniqID) ID() int64 {
	return atomic.LoadInt64((*int64)(u))
}

func CreateUniqueIDStorage() UniqueIDStorage {
	id := uniqID(0)
	return &id
}
