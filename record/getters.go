package record

type BoolGetter struct {
	Field
	Get func(item Record) bool
}

func (getter BoolGetter) Less(a, b Record) bool { return !getter.Get(a) && getter.Get(b) }

type InterfaceGetter struct {
	Field
	Get func(item Record) interface{}
}

type IntGetter struct {
	Field
	Get func(item Record) int
}

func (getter IntGetter) Less(a, b Record) bool { return getter.Get(a) < getter.Get(b) }

type Enum8Getter struct {
	Field
	Get func(item Record) Enum8
}

func (getter Enum8Getter) Less(a, b Record) bool {
	return getter.Get(a).Value() < getter.Get(b).Value()
}

type Enum16Getter struct {
	Field
	Get func(item Record) Enum16
}

func (getter Enum16Getter) Less(a, b Record) bool {
	return getter.Get(a).Value() < getter.Get(b).Value()
}

type Int32Getter struct {
	Field
	Get func(item Record) int32
}

func (getter Int32Getter) Less(a, b Record) bool { return getter.Get(a) < getter.Get(b) }

type Int64Getter struct {
	Field
	Get func(item Record) int64
}

func (getter Int64Getter) Less(a, b Record) bool { return getter.Get(a) < getter.Get(b) }

type StringGetter struct {
	Field
	Get func(item Record) string
}

func (getter StringGetter) Less(a, b Record) bool { return getter.Get(a) < getter.Get(b) }

type MapGetter struct {
	Field
	Get func(item Record) Map
}

type SetGetter struct {
	Field
	Get func(item Record) Set
}
