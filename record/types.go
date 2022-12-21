package record

import "strings"

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

type StringsSet struct {
	String string
	Set    map[string]struct{}
}

func NewStringsSet(str string) StringsSet {
	s := make(map[string]struct{})
	if len(str) > 0 {
		for _, item := range strings.Split(str, ",") {
			s[item] = struct{}{}
		}
	}
	return StringsSet{
		String: str,
		Set:    s,
	}
}
