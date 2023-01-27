package set

import (
	"sync"
	"testing"
)

type mapRWMutex struct {
	sync.RWMutex
	data map[int64]struct{}
}

func (m *mapRWMutex) Iterate(f func(id int64)) {
	m.RLock()
	for id := range m.data {
		f(id)
	}
	m.RUnlock()
}

func Benchmark_GoMapSet(b *testing.B) {
	const (
		concurrent = 100
		records    = 10_000
	)

	b.Run("map_rwmutex", func(b *testing.B) {
		var wg sync.WaitGroup
		sem := make(chan struct{}, concurrent)
		store := &mapRWMutex{
			data: make(map[int64]struct{}),
		}
		for i := 1; i < records; i++ {
			store.data[int64(i)] = struct{}{}
		}

		b.ResetTimer()

		wg.Add(1)
		go func() {
			for i := 1; i < records; i++ {
				store.Lock()
				store.data[int64(i)] = struct{}{}
				store.Unlock()
			}
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			for i := 1; i < records; i++ {
				store.Lock()
				delete(store.data, int64(i))
				store.Unlock()
			}
			wg.Done()
		}()

		wg.Add(b.N)
		for i := 0; i < b.N; i++ {
			go func() {
				sem <- struct{}{}
				store.Iterate(func(_ int64) {})
				<-sem
				wg.Done()
			}()
		}

		wg.Wait()
	})

	b.Run("set", func(b *testing.B) {
		var wg sync.WaitGroup
		sem := make(chan struct{}, concurrent)
		store := NewSet()
		for i := 1; i < records; i++ {
			store.Add(int64(i))
		}

		b.ResetTimer()

		wg.Add(1)
		go func() {
			for i := 1; i < records; i++ {
				store.Add(int64(i))
			}
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			for i := 1; i < records; i++ {
				store.Delete(int64(i))
			}
			wg.Done()
		}()

		wg.Add(b.N)
		for i := 0; i < b.N; i++ {
			go func() {
				sem <- struct{}{}
				store.Iterate(func(_ int64) {})
				<-sem
				wg.Done()
			}()
		}

		wg.Wait()
	})
}
