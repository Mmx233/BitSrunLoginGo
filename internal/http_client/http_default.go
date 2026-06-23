//go:build !linux

package http_client

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"github.com/Mmx233/BitSrunLoginGo/tools"
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
				Resolver:  newResolver(),
			}).DialContext,
			TLSHandshakeTimeout: config.Timeout,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.Settings.Basic.SkipCertVerify,
			},
			Proxy: http.ProxyFromEnvironment,
		},
		Timeout: config.Timeout,
	}
}
