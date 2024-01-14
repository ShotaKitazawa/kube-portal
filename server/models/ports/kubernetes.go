package ports

import (
	"context"

	"github.com/ShotaKitazawa/kube-portal/server/models"
)

type Kubernetes interface {
	ListIngress(ctx context.Context) (models.IngressInfoList, error)
	ListExternalLink(ctx context.Context) (models.IngressInfoList, error)
}
