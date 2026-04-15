package handler

import (
	"context"
	"net/url"

	"github.com/ShotaKitazawa/kube-portal/internal/api"
	"github.com/ShotaKitazawa/kube-portal/internal/model"
	"github.com/ShotaKitazawa/kube-portal/internal/model/port"
)

type Handler struct {
	k8sClient         port.Kubernetes
	showUntaggedLinks bool
}

var _ api.Handler = (*Handler)(nil)

func NewHandler(k8sClient port.Kubernetes, showUntaggedLinks bool) *Handler {
	return &Handler{k8sClient: k8sClient, showUntaggedLinks: showUntaggedLinks}
}

func (h *Handler) ListIngressInfo(ctx context.Context) ([]api.IngressInfo, error) {
	isLogin := getIsLogin(ctx)

	res1, err := h.k8sClient.ListIngress(ctx)
	if err != nil {
		return nil, err
	}
	res2, err := h.k8sClient.ListExternalLink(ctx)
	if err != nil {
		return nil, err
	}
	links := append(res1, res2...).ExcludePrivateLinkIfNotLogIn(isLogin)

	return toIngressInfoList(links, h.showUntaggedLinks), nil
}

func toIngressInfoList(links model.LinkList, showUntaggedLinks bool) []api.IngressInfo {
	var result []api.IngressInfo
	for _, link := range links {
		tags := link.Tags
		if showUntaggedLinks {
			tags = append(tags, "")
		}
		u := &url.URL{
			Scheme: link.Proto,
			Host:   link.Hostname,
			Path:   link.Path,
		}
		result = append(result, api.IngressInfo{
			Name:    link.Name,
			URL:     u.String(),
			IconURL: link.IconUrl,
			Tags:    tags,
		})
	}
	return result
}
