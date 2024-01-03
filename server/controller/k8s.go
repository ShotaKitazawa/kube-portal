package controller

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/kube-portal/server/models/ports"
	"github.com/ShotaKitazawa/kube-portal/server/view"
)

var (
	contextKeyIsLogin = "is_login"
)

type K8sController struct {
	logger    *logrus.Logger
	k8sClient ports.Kubernetes
	view      view.Iface

	jwtSecret string
}

func NewK8sController(l *logrus.Logger, k8sClient ports.Kubernetes, jwtSecret string) *K8sController {
	return &K8sController{l, k8sClient, &view.JSON{}, jwtSecret}
}

func (c K8sController) ListIngressInfo(ctx echo.Context) error {
	isLoggin := ctx.Get(contextKeyIsLogin).(bool)

	// logic
	l, err := c.k8sClient.ListIngressInfo(ctx.Request().Context())
	if err != nil {
		return err
	}
	l = l.ExcludePrivateLinkIfNotLogIn(isLoggin)

	// view
	return c.view.ListIngressInfo(ctx, l)
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
