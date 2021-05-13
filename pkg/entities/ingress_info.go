package entities

import "context"

type KubernetesPort interface {
	ListIngressInfo(ctx context.Context) ([]IngressInfo, error)
}

type IngressInfo struct {
	Fqdn      string
	Proto     string
	IconUrl   string
	IsPrivate bool
}
