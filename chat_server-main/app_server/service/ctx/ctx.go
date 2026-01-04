package ctx

import (
	"context"
	"strings"

	"connectrpc.com/connect"
)

func CtxInterceptor(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		header := req.Header()
		for k, v := range header {
			if strings.HasPrefix(k, "X-App-") && len(v) > 0 {
				ctx = context.WithValue(ctx, k, v[0])
			}
		}
		return next(ctx, req)
	})
}
