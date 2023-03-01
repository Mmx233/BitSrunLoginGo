package dnsUtil

import (
	"errors"
	"strings"
)

func DecodeDomain(domain string) (subStr string, rootDomain string, e error) {
	t := strings.Split(domain, ".")
	if len(t) == 1 {
		return "", "", errors.New("域名不合法")
	} else if len(t) == 2 {
		return "@", domain, nil
	}

	l := len(t)
	return strings.Join(t[:l-2], "."), strings.Join(t[l-2:l], "."), nil
}
