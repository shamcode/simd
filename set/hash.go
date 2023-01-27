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
	h := (prime5 + 8) ^ k1
	h = bits.RotateLeft64(h, 27)*prime1 + prime4

	h ^= h >> 33
	h *= prime2
	h ^= h >> 29
	h *= prime3
	h ^= h >> 32

	return uintptr(h)
}
