package record

type MapValueComparator[V any] interface {
	Compare(value V) (bool, error)
}

type Map[K comparable, V any] interface {
	HasKey(key K) bool
	HasValue(check MapValueComparator[V]) (bool, error)
}
