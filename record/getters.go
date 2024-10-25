package record

type LessComparable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

type Getter[R Record, T any] struct {
	Field
	Get func(item R) T
}

type (
	BoolGetter[R Record]                         Getter[R, bool]
	InterfaceGetter[R Record]                    Getter[R, any]
	ComparableGetter[R Record, T LessComparable] Getter[R, T]
	StringGetter[R Record]                       Getter[R, string]
	EnumGetter[R Record, T LessComparable]       Getter[R, Enum[T]]
	MapGetter[R Record]                          Getter[R, Map]
	SetGetter[R Record]                          Getter[R, Set]
)

func (getter BoolGetter[R]) Less(a, b R) bool { return !getter.Get(a) && getter.Get(b) }

func (getter ComparableGetter[R, T]) Less(a, b R) bool { return getter.Get(a) < getter.Get(b) }

func (getter StringGetter[R]) Less(a, b R) bool { return getter.Get(a) < getter.Get(b) }

func (getter EnumGetter[R, T]) Less(a, b R) bool {
	return getter.Get(a).Value() < getter.Get(b).Value()
}
