package util

import (
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/tool"
	"github.com/corpix/uarand"
	"net"
	"net/http"
)

type Http struct {
	Client *http.Client
	Header http.Header
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
	var header = make(http.Header, 2)
	header.Add("User-Agent", uarand.GetRandom())
	header.Set("X-Requested-With", "XMLHttpRequest")

	return &Http{
		Client: tool.GenHttpClient(&tool.HttpClientOptions{
			Transport: tool.GenHttpTransport(&tool.HttpTransportOptions{
				Timeout:           global.Timeout,
				LocalAddr:         addr,
				SkipSslCertVerify: global.Config.Settings.Basic.SkipCertVerify,
			}),
			Timeout: global.Timeout,
		}),
		Header: header,
	}
}
