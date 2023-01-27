package storage

var _ IDIterator = (*MapIDStorage)(nil)

type MapIDStorage map[int64]struct{}

func (s MapIDStorage) Iterate(f func(id int64)) {
	for id := range s {
		f(id)
	}
}

func CreateMapIDStorage() MapIDStorage {
	return make(map[int64]struct{})
}
