package controller

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"github.com/ShotaKitazawa/kube-portal/server/entities"
	"github.com/ShotaKitazawa/kube-portal/server/view"
)

var (
	contextKeyIsLogin = "is_login"
)

type K8sController struct {
	k8sClient entities.KubernetesPort
	view      view.ViewPort

	jwtSecret string
}

func NewK8sController(k8sClient entities.KubernetesPort, jwtSecret string) *K8sController {
	return &K8sController{k8sClient, &view.JsonView{}, jwtSecret}
}

func (c K8sController) ListIngressInfo(ctx echo.Context) error {
	isLoggin := ctx.Get(contextKeyIsLogin).(bool)

	// logic
	list, err := c.k8sClient.ListIngressInfo(ctx.Request().Context())
	if err != nil {
		return err
	}
	list = entities.ExcludePrivateLinkIfNotLogIn(list, isLoggin)

	// view
	return c.view.ListIngressInfo(ctx, list)
}

func (c K8sController) JwtAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var flagIsLogin bool
		if cookie, err := ctx.Cookie("jwt"); err == nil {
			reqJwt := cookie.Value
			// TODO Authn jwt
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
