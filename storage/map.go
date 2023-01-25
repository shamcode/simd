package storage

var _ LockableIDStorage = (*MapIDStorage)(nil)

type MapIDStorage map[int64]struct{}

func (s MapIDStorage) RLock()                               {}
func (s MapIDStorage) RUnlock()                             {}
func (s MapIDStorage) ThreadUnsafeData() map[int64]struct{} { return s }

func CreateMapIDStorage() MapIDStorage {
	return make(map[int64]struct{})
}
