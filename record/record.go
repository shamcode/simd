package record

type Record interface {
	GetID() int64
}

// NewIDGetter create ID getter.
func NewIDGetter[R Record]() ComparableGetter[R, int64] {
	return ComparableGetter[R, int64]{
		Field: field{
			index: 0,
			name:  "ID",
		},
		Get: R.GetID,
	}
}
