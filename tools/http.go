package tools

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/global"
	"net/http"
)

type Http struct {
	Client *http.Client
}

var HttpPack *Http

var httpTools map[string]*Http

func init() {
	if global.Config.Settings.Basic.Interfaces == "" {
		HttpPack = genHttpPack(nil)
	} else {
		httpTools = make(map[string]*Http, 0)
	}
}

func HttpPackSelect(eth *Eth) *Http {
	if HttpPack != nil {
		return HttpPack
	}
	if addrHttp, ok := httpTools[eth.Name]; ok {
		return addrHttp
	} else {
		addrHttp = genHttpPack(eth)
		httpTools[eth.Name] = addrHttp
		return addrHttp
	}
}
