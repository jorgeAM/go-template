package collections

func KeyBy[T any, K comparable](collection []T, iteratee func(T) K) map[K]T {
	result := make(map[K]T)
	for _, item := range collection {
		key := iteratee(item)
		result[key] = item
	}
	return result
}
