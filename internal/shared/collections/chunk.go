package collections

func Chunks[T any](collection []T, size int) [][]T {
	if len(collection) == 0 {
		return nil
	}

	divided := make([][]T, (len(collection)+size-1)/size)
	prev := 0
	i := 0
	till := len(collection) - size

	for prev < till {
		next := prev + size
		divided[i] = collection[prev:next]
		prev = next
		i++
	}

	divided[i] = collection[prev:]

	return divided
}
