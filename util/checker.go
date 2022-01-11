package util

import (
	"github.com/Mmx233/tool"
	"net/http"
)

type checker struct{}

var Checker checker

// NetOk 网络状况检查
func (a *checker) NetOk(transport *http.Transport) bool {

	h, i, e := tool.HTTP.GetReader(&tool.GetRequest{
		Url:       "https://www.baidu.com/",
		Redirect:  false,
		Transport: transport,
	})
	if e != nil {
		return false
	}
	_ = i.Close()
	return h.Get("Location") == ""
}
