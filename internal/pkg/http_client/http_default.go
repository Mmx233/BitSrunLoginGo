//go:build !linux

package http_client

import (
	"crypto/tls"
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"github.com/Mmx233/BitSrunLoginGo/tools"
	"net"
	"net/http"
)

func CreateClientFromEth(eth *tools.Eth) *http.Client {
	var addr net.Addr
	if eth != nil {
		addr = eth.Addr
	}

	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   config.Timeout,
				LocalAddr: addr,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.Settings.Basic.SkipCertVerify,
			},
			Proxy: http.ProxyFromEnvironment,
		},
		Timeout: config.Timeout,
	}
}
