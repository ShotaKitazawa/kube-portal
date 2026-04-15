package handler

import "context"

type contextKey string

const contextKeyIsLogin contextKey = "is_login"

func setIsLogin(ctx context.Context, v bool) context.Context {
	return context.WithValue(ctx, contextKeyIsLogin, v)
}

func getIsLogin(ctx context.Context) bool {
	v, _ := ctx.Value(contextKeyIsLogin).(bool)
	return v
}
