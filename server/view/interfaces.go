package view

import (
	"github.com/ShotaKitazawa/kube-portal/server/entities"
	"github.com/labstack/echo/v4"
)

type ViewPort interface {
	ListIngressInfo(ctx echo.Context, list []entities.IngressInfo) error
}
