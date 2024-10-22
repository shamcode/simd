package set

import "math/bits"

const (
	prime1 uint64 = 11400714785074694791
	prime2 uint64 = 14029467366897019727
	prime3 uint64 = 1609587929392839161
	prime4 uint64 = 9650029242287828579
	prime5 uint64 = 2870177450012600261
)

func xxHashQword(key int64) uintptr {
	k1 := uint64(key) * prime2
	k1 = bits.RotateLeft64(k1, 31)
	k1 *= prime1
	hash := (prime5 + 8) ^ k1
	hash = bits.RotateLeft64(hash, 27)*prime1 + prime4

	hash ^= hash >> 33
	hash *= prime2
	hash ^= hash >> 29
	hash *= prime3
	hash ^= hash >> 32

	return uintptr(hash)
}
