package ports

import (
	"context"

	"github.com/ShotaKitazawa/kube-portal/server/models"
)

type Kubernetes interface {
	ListIngressInfo(ctx context.Context) (models.IngressInfoList, error)
}
