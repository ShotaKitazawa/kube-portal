package util

import (
	"github.com/labstack/echo/v4"
)

var contextKeyIsLogin = "is_login"

func SetContextIsLogin(ctx echo.Context, isLogin bool) {
	ctx.Set(contextKeyIsLogin, isLogin)
}

func GetContextIsLogin(ctx echo.Context) bool {
	res := ctx.Get(contextKeyIsLogin)
	if res == nil {
		return false
	}
	return res.(bool)
}
