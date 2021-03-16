package Util

import (
	"net"
)

type checker struct{}

var Checker checker

func (checker) NetOk() bool {
	if ip, err := net.LookupIP("www.msftconnecttest.com"); err != nil {
		return false
	} else if len(ip) == 0 || ip[0].String() != "13.107.4.52" {
		return false
	}
	return true
}
