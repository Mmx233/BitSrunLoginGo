package http_client

import (
	"crypto/tls"
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"github.com/Mmx233/BitSrunLoginGo/tools"
	"github.com/Mmx233/tool"
	"net"
	"net/http"
	"syscall"
)

func genHttpPack(eth *tools.Eth) *Http {
	dialer := net.Dialer{
		Timeout: config.Timeout,
	}
	if eth != nil {
		dialer.LocalAddr = eth.Addr
		ethName := eth.Name
		dialer.Control = func(network string, address string, c syscall.RawConn) error {
			var operr error
			fn := func(fd uintptr) {
				operr = syscall.SetsockoptString(int(fd), syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, ethName)
			}
			if err := c.Control(fn); err != nil {
				return err
			}
			if operr != nil {
				return operr
			}
			return nil
		}
	}

	tr := &http.Transport{
		DialContext:         dialer.DialContext,
		TLSHandshakeTimeout: config.Timeout,
		IdleConnTimeout:     config.Timeout,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: config.Settings.Basic.SkipCertVerify},
	}
	tr.Proxy = http.ProxyFromEnvironment
	return &Http{
		Client: tool.GenHttpClient(&tool.HttpClientOptions{
			Transport: tr,
			Timeout:   config.Timeout,
		}),
	}
}
