package model

type Link struct {
	Name      string
	Hostname  string
	Path      string
	Proto     string
	IconUrl   string
	Tags      []string
	IsPrivate bool
}

type LinkList []Link

func (list LinkList) ExcludePrivateLinkIfNotLogIn(isLogin bool) LinkList {
	var result []Link
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
