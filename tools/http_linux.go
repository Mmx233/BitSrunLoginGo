package tools

import (
	"crypto/tls"
	"github.com/Mmx233/BitSrunLoginGo/internal/global"
	"github.com/Mmx233/tool"
	"net"
	"net/http"
	"syscall"
)

func genHttpPack(eth *Eth) *Http {
	dialer := net.Dialer{
		Timeout: global.Timeout,
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
		TLSHandshakeTimeout: global.Timeout,
		IdleConnTimeout:     global.Timeout,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: global.Config.Settings.Basic.SkipCertVerify},
	}
	tr.Proxy = http.ProxyFromEnvironment
	return &Http{
		Client: tool.GenHttpClient(&tool.HttpClientOptions{
			Transport: tr,
			Timeout:   global.Timeout,
		}),
	}
}
