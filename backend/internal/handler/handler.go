package handler

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ShotaKitazawa/kube-portal/internal/api"
	"github.com/ShotaKitazawa/kube-portal/internal/model"
	"github.com/ShotaKitazawa/kube-portal/internal/model/port"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OIDCConfig struct {
	Issuer   string
	ClientID string
	Audience string
}

type Handler struct {
	k8sClient         port.Kubernetes
	showUntaggedLinks bool
	disableOIDC       bool
	oidcConfig        OIDCConfig
	provider          *oidc.Provider
}

var _ api.Handler = (*Handler)(nil)

func NewHandler(k8sClient port.Kubernetes, showUntaggedLinks bool, disableOIDC bool, oidcConfig OIDCConfig, provider *oidc.Provider) *Handler {
	return &Handler{
		k8sClient:         k8sClient,
		showUntaggedLinks: showUntaggedLinks,
		disableOIDC:       disableOIDC,
		oidcConfig:        oidcConfig,
		provider:          provider,
	}
}

func (h *Handler) GetOidcConfig(_ context.Context) (*api.OIDCConfig, error) {
	if h.disableOIDC {
		return &api.OIDCConfig{Enabled: false}, nil
	}
	return &api.OIDCConfig{
		Enabled:  true,
		Issuer:   api.NewOptString(h.oidcConfig.Issuer),
		ClientID: api.NewOptString(h.oidcConfig.ClientID),
		Audience: api.NewOptString(h.oidcConfig.Audience),
	}, nil
}

func (h *Handler) GetUserinfo(ctx context.Context) (api.GetUserinfoRes, error) {
	if h.disableOIDC {
		return &api.UserInfo{
			Sub:   "dev-user",
			Name:  api.NewOptString("Dev User"),
			Email: api.NewOptString("dev@example.com"),
		}, nil
	}

	token := getToken(ctx)
	if token == "" {
		return &api.GetUserinfoUnauthorized{}, nil
	}

	userInfo, err := h.provider.UserInfo(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	if err != nil {
		return nil, fmt.Errorf("failed to get userinfo: %w", err)
	}

	var claims struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}
	if err := userInfo.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse userinfo claims: %w", err)
	}

	result := &api.UserInfo{Sub: userInfo.Subject}
	if claims.Name != "" {
		result.Name = api.NewOptString(claims.Name)
	}
	if claims.Email != "" {
		result.Email = api.NewOptString(claims.Email)
	}
	if claims.Picture != "" {
		result.Picture = api.NewOptString(claims.Picture)
	}
	return result, nil
}

func (h *Handler) ListIngressInfo(ctx context.Context) ([]api.IngressInfo, error) {
	isLogin := getIsLogin(ctx) || h.disableOIDC

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
