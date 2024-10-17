package record

type Record interface {
	GetID() int64
}

// ID is a common getter for all types Records.
var ID = Int64Getter{
	Field: field{
		index: 0,
		name:  "ID",
	},
	Get: func(item Record) int64 {
		return item.GetID()
	},
}
