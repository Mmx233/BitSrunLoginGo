//go:build !linux

package tools

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/global"
	"github.com/Mmx233/tool"
	"net"
	"net/http"
)

func genHttpPack(eth *Eth) *Http {
	var addr net.Addr
	if eth != nil {
		addr = eth.Addr
	}
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
