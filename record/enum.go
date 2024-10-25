package record

type Enum[T LessComparable] interface {
	Value() T
}
