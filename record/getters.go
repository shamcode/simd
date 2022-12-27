package record

type BoolGetter struct {
	Field string
	Get   func(item interface{}) bool
}

type InterfaceGetter struct {
	Field string
	Get   func(item interface{}) interface{}
}

type IntGetter struct {
	Field string
	Get   func(item interface{}) int
}

type Enum8Getter struct {
	Field string
	Get   func(item interface{}) Enum8
}

type Enum16Getter struct {
	Field string
	Get   func(item interface{}) Enum16
}

type Int32Getter struct {
	Field string
	Get   func(item interface{}) int32
}

type Int64Getter struct {
	Field string
	Get   func(item interface{}) int64
}

type StringGetter struct {
	Field string
	Get   func(item interface{}) string
}

type MapGetter struct {
	Field string
	Get   func(item interface{}) Map
}

type SetGetter struct {
	Field string
	Get   func(item interface{}) Set
}
