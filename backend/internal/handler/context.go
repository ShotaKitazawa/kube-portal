package handler

import "context"

type contextKey string

const (
	contextKeyIsLogin contextKey = "is_login"
	contextKeyToken   contextKey = "token"
)

func setIsLogin(ctx context.Context, v bool) context.Context {
	return context.WithValue(ctx, contextKeyIsLogin, v)
}

func getIsLogin(ctx context.Context) bool {
	v, _ := ctx.Value(contextKeyIsLogin).(bool)
	return v
}

func setToken(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, contextKeyToken, v)
}

func getToken(ctx context.Context) string {
	v, _ := ctx.Value(contextKeyToken).(string)
	return v
}
