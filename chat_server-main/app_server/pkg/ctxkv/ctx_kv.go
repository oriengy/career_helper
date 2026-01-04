package ctxkv

import "context"

func GetCtxKvString(ctx context.Context, key string) string {
	v, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}
	return v
}

func GetCtxKvInt(ctx context.Context, key string) int {
	switch v := ctx.Value(key).(type) {
	case int:
		return v
	case int64:
		return int(v)
	}
	return 0
}
