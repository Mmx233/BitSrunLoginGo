package global

import (
	"github.com/Mmx233/tool"
	"net"
	"net/http"
)

var transport *http.Transport

var transports map[net.Addr]*http.Transport

func Transports(addr net.Addr) *http.Transport {
	if transport != nil {
		return transport
	}
	if transport, ok := transports[addr]; ok {
		return transport
	} else {
		transport = tool.HTTP.GenTransport(&tool.GenTransport{
			Timeout:   Timeout,
			LocalAddr: addr,
		})
		transports[addr] = transport
		return transport
	}
}

func initTransport() {
	if Config.Settings.Basic.Interfaces == "" {
		transport = tool.HTTP.GenTransport(&tool.GenTransport{
			Timeout: Timeout,
		})
	} else {
		transports = make(map[net.Addr]*http.Transport, 0)
	}
}
