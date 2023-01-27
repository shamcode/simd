package storage

import (
	"github.com/shamcode/simd/set"
)

var _ IDStorage = (*setStorage)(nil)

type setStorage set.Set

func (s *setStorage) Count() int {
	return (*set.Set)(s).Count()
}

func (s *setStorage) Add(id int64) {
	(*set.Set)(s).Add(id)
}

func (s *setStorage) Delete(id int64) {
	(*set.Set)(s).Delete(id)
}

func (s *setStorage) Iterate(f func(id int64)) {
	(*set.Set)(s).Iterate(f)
}

func CreateSetIDStorage() IDStorage {
	return (*setStorage)(set.NewSet())
}
