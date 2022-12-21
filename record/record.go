package record

type Record interface {
	GetID() int64
	ComputeFields()
}
