package ref

func Ref[T any](t T) *T {
	return &t
}
