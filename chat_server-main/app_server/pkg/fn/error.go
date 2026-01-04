package fn

func NoErr[T any](v T, err error) T {
	return v
}
