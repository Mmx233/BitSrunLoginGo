//go:build !linux

package tools

import (
	"crypto/tls"
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
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   config.Timeout,
			LocalAddr: addr,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.Settings.Basic.SkipCertVerify,
		},
		Proxy: http.ProxyFromEnvironment,
	}
	return &Http{
		Client: tool.GenHttpClient(&tool.HttpClientOptions{
			Transport: tr,
			Timeout:   config.Timeout,
		}),
	}
}
