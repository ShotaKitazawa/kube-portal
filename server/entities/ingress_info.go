package entities

import "context"

type KubernetesPort interface {
	ListIngressInfo(ctx context.Context) ([]IngressInfo, error)
}

type IngressInfo struct {
	Name      string
	Fqdn      string
	Proto     string
	IconUrl   string
	IsPrivate bool
}

func ExcludePrivateLinkIfNotLogIn(list []IngressInfo, isLogin bool) []IngressInfo {
	var result []IngressInfo
	if !isLogin {
		for _, val := range list {
			if !val.IsPrivate {
				result = append(result, val)
			}
		}
	} else {
		result = list
	}
	return result
}
