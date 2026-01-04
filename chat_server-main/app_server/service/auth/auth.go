package auth

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"app_server/pkg/fn"
	"app_server/pkg/jwt"

	"connectrpc.com/connect"
)

func AuthInterceptor(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		authToken := req.Header().Get("Authorization")

		userID, err := ParseUserID(authToken)
		if err != nil {
			return nil, err
		}

		// 设置用户ID到上下文
		ctx = context.WithValue(ctx, userIDKey, userID)

		return next(ctx, req)
	})
}

func ParseUserID(authToken string) (uint, error) {
	authToken = strings.TrimPrefix(authToken, "Bearer ")
	if authToken == "" {
		return 0, connect.NewError(connect.CodeUnauthenticated, errors.New("access token is required"))
	}

	// 解析token
	userIDStr, err := jwt.Get().ParseToken(authToken)
	if err != nil {
		slog.Error("parse token error", "error", err, "authToken", authToken)
		return 0, connect.NewError(connect.CodeUnauthenticated, err)
	}

	// 验证用户ID
	userIDInt := fn.Atoi[uint](userIDStr)
	if userIDInt == 0 {
		return 0, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid user id"))
	}

	return userIDInt, nil
}

func GetUserID(ctx context.Context) uint {
	return ctx.Value(userIDKey).(uint)
}

// SetUserIDToContext 设置用户ID到上下文（仅用于测试）
func SetUserIDToContext(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

type ctxKey int

const (
	userIDKey ctxKey = iota
	userIDStrKey
)
