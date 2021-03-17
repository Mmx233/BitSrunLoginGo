package Util

import (
	"Mmx/Request"
	"fmt"
	"net"
	"time"
)

type checker struct{}

var Checker checker

func (checker) NetOk(url string) bool {
	if ip, err := net.LookupIP("www.msftconnecttest.com"); err != nil { //通过DNS确认是否在线
		return false
	} else if len(ip) == 0 || ip[0].String() != "13.107.4.52" {
		return false
	}

	{ //有些路由器有DNS缓存，故进行进一步确认
		body, err := Request.Get(url, map[string]string{
			"callback": "jQuery1635413",
			"_":        fmt.Sprint(time.Now().UnixNano()),
		})
		ErrHandler(err)
		r, err := GetResult(body)
		if err != nil {
			ErrHandler(err)
		}
		if r != "ok" {
			return false
		}
	}
	return true
}
