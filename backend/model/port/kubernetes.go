package port

import (
	"context"

	"github.com/ShotaKitazawa/kube-portal/backend/model"
)

type Kubernetes interface {
	ListIngress(ctx context.Context) (model.LinkList, error)
	ListExternalLink(ctx context.Context) (model.LinkList, error)
}
