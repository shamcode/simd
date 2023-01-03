package record

type Record interface {
	GetID() int64

	// ComputeFields is a special hook for optimize slow computing fields.
	// ComputeFields call on insert or update record.
	ComputeFields()
}
