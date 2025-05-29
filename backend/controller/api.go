package controller

import (
	"log/slog"

	"github.com/labstack/echo/v4"

	"github.com/ShotaKitazawa/kube-portal/backend/model/port"
	"github.com/ShotaKitazawa/kube-portal/backend/util"
	"github.com/ShotaKitazawa/kube-portal/backend/view"
)

type APIController struct {
	logger    *slog.Logger
	k8sClient port.Kubernetes
	view      view.API
}

func NewAPIController(l *slog.Logger, k8sClient port.Kubernetes, v view.API) *APIController {
	return &APIController{l, k8sClient, v}
}

func (c APIController) ListIngressInfo(ctx echo.Context) error {
	isLoggin := util.GetContextIsLogin(ctx)

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
