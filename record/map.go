package record

type MapValueComparator interface {
	Compare(value interface{}) (bool, error)
}

type Map interface {
	HasKey(key interface{}) bool
	HasValue(check MapValueComparator) (bool, error)
}
