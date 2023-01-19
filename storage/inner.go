package storage

var _ LockableIDStorage = (*innerIDStorage)(nil)

type innerIDStorage struct {
	data map[int64]struct{}
}

func (s *innerIDStorage) RLock()                               {}
func (s *innerIDStorage) RUnlock()                             {}
func (s *innerIDStorage) ThreadUnsafeData() map[int64]struct{} { return s.data }

func createInnerIDStorage() *innerIDStorage {
	return &innerIDStorage{
		data: make(map[int64]struct{}),
	}
}
