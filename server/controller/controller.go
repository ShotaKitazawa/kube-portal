package controller

import (
	"github.com/labstack/echo/v4"

	"github.com/ShotaKitazawa/kube-portal/server/entities"
)

type Controller struct {
	k8sClient entities.KubernetesPort
}

func New(k8sClient entities.KubernetesPort) *Controller {
	return &Controller{k8sClient}
}

func (c Controller) ListIngressInfo(ctx echo.Context) error {
	list, err := c.k8sClient.ListIngressInfo(ctx.Request().Context())
	if err != nil {
		return err
	}
	return ctx.JSON(200, list)
}
