// Contains some generic helper function and datastructure
package utils

// Generric Data Structure
type AssociativeArray[K, V any] []KV[K, V]
type KV[K, V any] struct {
	Key   K
	Value V
}

// Util function for casting
func Cast[T any](in any) (T, bool) {
	var zero T
	switch inter := in.(type) {
	case T:
		return inter, true
	default:
		return zero, false
	}
}
