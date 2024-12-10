package record

type Set[T comparable] interface {
	Has(item T) bool
}
