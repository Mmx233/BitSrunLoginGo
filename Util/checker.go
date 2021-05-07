package Util

import (
	"context"
	"net"
)

type checker struct{}

var Checker checker

func (checker) NetOk() bool {
	r := &net.Resolver{ //指定DNS，防止本地DNS缓存影响
		PreferGo: true,
		Dial:     NetDailEr(),
	}
	if ip, err := r.LookupIP(context.Background(), "ip4", "www.msftconnecttest.com"); err != nil { //通过DNS确认是否在线
		return false
	} else if len(ip) == 0 || ip[0].String() != "13.107.4.52" {
		return false
	}
	return true
}
