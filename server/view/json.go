package view

import (
	"github.com/ShotaKitazawa/kube-portal/server/entities"
	"github.com/labstack/echo/v4"
)

type JsonView struct{}

type IngressInfo struct {
	Name      string `json:"name"`
	Url       string `json:"url"`
	IconUrl   string `json:"icon_url"`
	IsPrivate bool   `json:"is_private"`
}

func (v JsonView) ListIngressInfo(ctx echo.Context, list []entities.IngressInfo) error {
	var res []IngressInfo
	for _, val := range list {
		res = append(res, IngressInfo{
			Name:      val.Name,
			Url:       val.Proto + "://" + val.Fqdn,
			IconUrl:   val.IconUrl,
			IsPrivate: val.IsPrivate,
		})
	}
	return ctx.JSON(200, res)
}
