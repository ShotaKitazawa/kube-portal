package models

type IngressInfo struct {
	Name      string
	Fqdn      string
	Path      string
	Proto     string
	IconUrl   string
	Tags      []string
	IsPrivate bool
}

type IngressInfoList []IngressInfo

func (list IngressInfoList) ExcludePrivateLinkIfNotLogIn(isLogin bool) IngressInfoList {
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
