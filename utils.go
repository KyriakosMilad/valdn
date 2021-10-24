package validation

import "strings"

func getRuleValue(rule string) string {
	if strings.ContainsRune(rule, ':') {
		return strings.Split(rule, ":")[1]
	}
	return ""
}
