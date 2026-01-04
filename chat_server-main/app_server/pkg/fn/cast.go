package fn

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

func Atoi[T constraints.Integer](s string) (v T) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}
	return T(i)
}

func Itoa[T constraints.Integer](v T) string {
	return strconv.FormatInt(int64(v), 10)
}

func CastNumber[A, B constraints.Integer | constraints.Float](a A) (v B) {
	return B(a)
}

func CastNumbers[A, B constraints.Integer | constraints.Float](a []A) (v []B) {
	v = make([]B, len(a))
	for i, a := range a {
		v[i] = B(a)
	}
	return
}
