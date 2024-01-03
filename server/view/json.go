package view

import (
	"github.com/ShotaKitazawa/kube-portal/server/models"
	"github.com/labstack/echo/v4"
)

type JSON struct{}

var _ Iface = (*JSON)(nil)

type IngressInfo struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	IconUrl string `json:"icon_url"`
}

func (v JSON) ListIngressInfo(ctx echo.Context, list []models.IngressInfo) error {
	var res []IngressInfo
	for _, val := range list {
		res = append(res, IngressInfo{
			Name:    val.Name,
			Url:     val.Proto + "://" + val.Fqdn,
			IconUrl: val.IconUrl,
		})
	}
	return ctx.JSON(200, res)
}
