package util

import "github.com/Mmx233/tool"

type checker struct{}

var Checker checker

func (checker) NetOk() bool {
	h, _, e := tool.HTTP.GetBytes(&tool.GetRequest{
		Url:      "https://www.baidu.com/",
		Redirect: false,
	})
	if e != nil || h.Get("Location") != "" {
		return false
	}
	return true
}
