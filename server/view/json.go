package view

import (
	"net/url"

	"github.com/ShotaKitazawa/kube-portal/server/models"
	"github.com/labstack/echo/v4"
)

type Iface interface {
	ListIngressInfo(ctx echo.Context, list []models.IngressInfo) error
}

type JSON struct {
	ShowUntaggedLinks bool
}

var _ Iface = (*JSON)(nil)

type IngressInfo struct {
	Name    string   `json:"name"`
	Url     string   `json:"url"`
	IconUrl string   `json:"icon_url"`
	Tags    []string `json:"tags"`
}

func (v JSON) ListIngressInfo(ctx echo.Context, list []models.IngressInfo) error {
	if v.ShowUntaggedLinks {
		for idx := range list {
			list[idx].Tags = append(list[idx].Tags, "")
		}
	}

	var res []IngressInfo
	for _, val := range list {
		u, err := url.Parse(val.Proto + "://" + val.Hostname + "/" + val.Path)
		if err != nil {
			return err
		}
		res = append(res, IngressInfo{
			Name:    val.Name,
			Url:     u.String(),
			IconUrl: val.IconUrl,
			Tags:    val.Tags,
		})
	}
	return ctx.JSON(200, res)
}
