package utils

import "strings"

func ContainsStr(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Split(s, sep string) []string {
	result := strings.Split(s, sep)
	if len(result) == 1 && result[0] == "" {
		return []string{}
	}
	return result
}
