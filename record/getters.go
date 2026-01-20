package record

type LessComparable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~string
}

type GetterInterface[R Record, T any] interface {
	Field
	GetForRecord(item R) T
}

type (
	Getter[R Record, T any] struct {
		Field

		Get func(item R) T
	}
	BoolGetter[R Record]                         Getter[R, bool]
	ComparableGetter[R Record, T LessComparable] Getter[R, T]
	MapGetter[R Record, K comparable, V any]     Getter[R, Map[K, V]]
	SetGetter[R Record, T comparable]            Getter[R, Set[T]]
)

func (getter Getter[R, T]) GetForRecord(item R) T               { return getter.Get(item) }
func (getter BoolGetter[R]) GetForRecord(item R) bool           { return getter.Get(item) }
func (getter ComparableGetter[R, T]) GetForRecord(item R) T     { return getter.Get(item) }
func (getter MapGetter[R, K, V]) GetForRecord(item R) Map[K, V] { return getter.Get(item) }
func (getter SetGetter[R, T]) GetForRecord(item R) Set[T]       { return getter.Get(item) }

func (getter BoolGetter[R]) Less(a, b R) bool          { return !getter.Get(a) && getter.Get(b) }
func (getter ComparableGetter[R, T]) Less(a, b R) bool { return getter.Get(a) < getter.Get(b) }
