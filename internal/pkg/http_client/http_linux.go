package http_client

import (
	"crypto/tls"
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"github.com/Mmx233/BitSrunLoginGo/tools"
	"net"
	"net/http"
	"syscall"
)

func CreateClientFromEth(eth *tools.Eth) *http.Client {
	dialer := net.Dialer{
		Timeout: config.Timeout,
	}
	if eth != nil {
		dialer.LocalAddr = eth.Addr
		ethName := eth.Name
		dialer.Control = func(network string, address string, c syscall.RawConn) error {
			var opErr error
			fn := func(fd uintptr) {
				opErr = syscall.SetsockoptString(int(fd), syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, ethName)
			}
			if err := c.Control(fn); err != nil {
				return err
			}
			if opErr != nil {
				return opErr
			}
			return nil
		}
	}

	return &http.Client{
		Transport: &http.Transport{
			DialContext:         dialer.DialContext,
			TLSHandshakeTimeout: config.Timeout,
			IdleConnTimeout:     config.Timeout,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: config.Settings.Basic.SkipCertVerify},
			Proxy:               http.ProxyFromEnvironment,
		},
		Timeout: config.Timeout,
	}
}
