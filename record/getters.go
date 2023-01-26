package record

type BoolGetter struct {
	Field Field
	Get   func(item Record) bool
}

type InterfaceGetter struct {
	Field Field
	Get   func(item Record) interface{}
}

type IntGetter struct {
	Field Field
	Get   func(item Record) int
}

type Enum8Getter struct {
	Field Field
	Get   func(item Record) Enum8
}

type Enum16Getter struct {
	Field Field
	Get   func(item Record) Enum16
}

type Int32Getter struct {
	Field Field
	Get   func(item Record) int32
}

type Int64Getter struct {
	Field Field
	Get   func(item Record) int64
}

type StringGetter struct {
	Field Field
	Get   func(item Record) string
}

type MapGetter struct {
	Field Field
	Get   func(item Record) Map
}

type SetGetter struct {
	Field Field
	Get   func(item Record) Set
}
