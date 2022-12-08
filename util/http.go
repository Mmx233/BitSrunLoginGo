package util

import (
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/tool"
	"net"
	"net/http"
)

type Http struct {
	Client *http.Client
}

var HttpPack *Http

var httpTools map[net.Addr]*Http

func init() {
	if global.Config.Settings.Basic.Interfaces == "" {
		HttpPack = genHttpPack(nil)
	} else {
		httpTools = make(map[net.Addr]*Http, 0)
	}
}

func HttpPackSelect(addr net.Addr) *Http {
	if HttpPack != nil {
		return HttpPack
	}
	if addrHttp, ok := httpTools[addr]; ok {
		return addrHttp
	} else {
		addrHttp = genHttpPack(addr)
		httpTools[addr] = addrHttp
		return addrHttp
	}
}

func genHttpPack(addr net.Addr) *Http {
	tr := tool.GenHttpTransport(&tool.HttpTransportOptions{
		Timeout:           global.Timeout,
		LocalAddr:         addr,
		SkipSslCertVerify: global.Config.Settings.Basic.SkipCertVerify,
	})
	tr.Proxy = http.ProxyFromEnvironment
	return &Http{
		Client: tool.GenHttpClient(&tool.HttpClientOptions{
			Transport: tr,
			Timeout:   global.Timeout,
		}),
	}
}
