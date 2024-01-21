package controller

import (
	"fmt"
	"log/slog"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"github.com/ShotaKitazawa/kube-portal/server/models/ports"
	"github.com/ShotaKitazawa/kube-portal/server/view"
)

var (
	contextKeyIsLogin = "is_login"
)

type K8sController struct {
	logger    *slog.Logger
	k8sClient ports.Kubernetes
	view      view.Iface

	jwtSecret string
}

func NewK8sController(l *slog.Logger, k8sClient ports.Kubernetes, jwtSecret string, showUntaggedLinks bool) *K8sController {
	v := &view.JSON{ShowUntaggedLinks: showUntaggedLinks}
	return &K8sController{l, k8sClient, v, jwtSecret}
}

func (c K8sController) ListIngressInfo(ctx echo.Context) error {
	isLoggin := ctx.Get(contextKeyIsLogin).(bool)

	// logic
	res1, err := c.k8sClient.ListIngress(ctx.Request().Context())
	if err != nil {
		return err
	}
	res2, err := c.k8sClient.ListExternalLink(ctx.Request().Context())
	if err != nil {
		return err
	}
	res := append(res1, res2...).ExcludePrivateLinkIfNotLogIn(isLoggin)

	// view
	return c.view.ListIngressInfo(ctx, res)
}

func (c K8sController) JwtAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var flagIsLogin bool
		if cookie, err := ctx.Cookie("jwt"); err == nil {
			reqJwt := cookie.Value
			if _, err := jwt.Parse(reqJwt, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(c.jwtSecret), nil
			}); err == nil {
				flagIsLogin = true
			}
		}
		ctx.Set(contextKeyIsLogin, flagIsLogin)
		return next(ctx)
	}
}
