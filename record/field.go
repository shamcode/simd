package record

type Field interface {
	Index() uint8
	String() string
}

type field struct {
	index uint8
	name  string
}

func (f field) Index() uint8   { return f.index }
func (f field) String() string { return f.name }

type FieldsConstructor uint8

func (fc *FieldsConstructor) New(name string) Field {
	index := uint8(*fc)
	*fc++
	return field{
		index: index,
		name:  name,
	}
}

func NewFields() *FieldsConstructor {
	var fc = FieldsConstructor(1) // starts with 1, because 0 reserved for ID field
	return &fc
}
