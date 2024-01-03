package view

import (
	"github.com/ShotaKitazawa/kube-portal/server/models"
	"github.com/labstack/echo/v4"
)

type Iface interface {
	ListIngressInfo(ctx echo.Context, list []models.IngressInfo) error
}
