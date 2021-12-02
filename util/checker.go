package util

import (
	"github.com/Mmx233/tool"
	"time"
)

type checker struct{}

var Checker checker

// NetOk 网络状况检查
func (checker) NetOk(timeout uint) bool {
	h, _, e := tool.HTTP.GetBytes(&tool.GetRequest{
		Url:      "https://www.baidu.com/",
		Redirect: false,
		Timeout:  time.Duration(timeout) * time.Second,
	})
	if e != nil || h.Get("Location") != "" {
		return false
	}
	return true
}
