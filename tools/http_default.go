//go:build !linux

package tools

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
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
		Timeout:           config.Timeout,
		LocalAddr:         addr,
		SkipSslCertVerify: config.Settings.Basic.SkipCertVerify,
	})
	tr.Proxy = http.ProxyFromEnvironment
	return &Http{
		Client: tool.GenHttpClient(&tool.HttpClientOptions{
			Transport: tr,
			Timeout:   config.Timeout,
		}),
	}
}
