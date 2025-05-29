package view

import (
	"net/url"

	"github.com/labstack/echo/v4"

	"github.com/ShotaKitazawa/kube-portal/backend/model"
)

type API interface {
	ListIngressInfo(ctx echo.Context, list []model.Link) error
}

type JSONForAPI struct {
	ShowUntaggedLinks bool
}

var _ API = (*JSONForAPI)(nil)

type IngressInfo struct {
	Name    string   `json:"name"`
	Url     string   `json:"url"`
	IconUrl string   `json:"icon_url"`
	Tags    []string `json:"tags"`
}

func (v JSONForAPI) ListIngressInfo(ctx echo.Context, list []model.Link) error {
	if v.ShowUntaggedLinks {
		for idx := range list {
			list[idx].Tags = append(list[idx].Tags, "")
		}
	}

	var res []IngressInfo
	for _, val := range list {
		u := &url.URL{
			Scheme: val.Proto,
			Host:   val.Hostname,
			Path:   val.Path,
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
