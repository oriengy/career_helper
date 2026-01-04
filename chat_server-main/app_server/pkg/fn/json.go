package fn

import jsoniter "github.com/json-iterator/go"

func JsonUnmarshalStr[T any](data string) (T, error) {
	var v T
	err := jsoniter.Unmarshal([]byte(data), &v)
	return v, err
}
