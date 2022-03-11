package util

import (
	"github.com/Mmx233/tool"
	"net/http"
)

type checker struct {
	url string
	set bool
}

var Checker = checker{
	url: "https://www.baidu.com/",
}

func (a *checker) SetUrl(url string) {
	if a.set {
		return
	}
	a.url = url
	a.set = true
}

// NetOk 网络状况检查
func (a *checker) NetOk(transport *http.Transport) bool {
	Log.Debug("GET ", a.url)
	res, e := tool.HTTP.GetReader(&tool.GetRequest{
		Url:       a.url,
		Redirect:  false,
		Transport: transport,
	})
	if e != nil {
		Log.Debug(e)
		return false
	}
	_ = res.Body.Close()
	return res.Header.Get("Location") == ""
}
