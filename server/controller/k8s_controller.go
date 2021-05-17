package controller

import (
	"github.com/labstack/echo/v4"

	"github.com/ShotaKitazawa/kube-portal/server/entities"
	"github.com/ShotaKitazawa/kube-portal/server/view"
)

type K8sController struct {
	k8sClient entities.KubernetesPort
	view      view.ViewPort
}

func NewK8sController(k8sClient entities.KubernetesPort) *K8sController {
	return &K8sController{k8sClient, &view.JsonView{}}
}

func (c K8sController) ListIngressInfo(ctx echo.Context) error {
	// logic
	list, err := c.k8sClient.ListIngressInfo(ctx.Request().Context())
	if err != nil {
		return err
	}

	// view
	return c.view.ListIngressInfo(ctx, list)
}
