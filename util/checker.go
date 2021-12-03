package util

import (
	"github.com/Mmx233/tool"
	"net"
	"time"
)

type checker struct{}

var Checker checker

// NetOk 网络状况检查
func (checker) NetOk(timeout uint, localAddr net.Addr) bool {
	h, i, e := tool.HTTP.GetReader(&tool.GetRequest{
		Url:       "https://www.baidu.com/",
		Redirect:  false,
		Timeout:   time.Duration(timeout) * time.Second,
		LocalAddr: localAddr,
	})
	if e != nil {
		return false
	}
	_ = i.Close()
	return h.Get("Location") == ""
}
