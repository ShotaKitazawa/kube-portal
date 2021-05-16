package controller

import (
	"github.com/labstack/echo/v4"

	"github.com/ShotaKitazawa/kube-portal/server/entities"
	"github.com/ShotaKitazawa/kube-portal/server/view"
)

type Controller struct {
	k8sClient entities.KubernetesPort
	view      view.ViewPort
}

func New(k8sClient entities.KubernetesPort) *Controller {
	return &Controller{k8sClient, &view.JsonView{}}
}

func (c Controller) ListIngressInfo(ctx echo.Context) error {
	// logic
	list, err := c.k8sClient.ListIngressInfo(ctx.Request().Context())
	if err != nil {
		return err
	}

	// view
	return c.view.ListIngressInfo(ctx, list)
}
