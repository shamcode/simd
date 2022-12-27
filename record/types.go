package record

type Enum8 interface {
	Value() uint8
}

type Enum16 interface {
	Value() uint16
}

type Set interface {
	Has(item interface{}) bool
}

type Comparator interface {
	Compare(item interface{}) bool
}

type Map interface {
	HasKey(key interface{}) bool
	HasValue(check Comparator) bool
}
