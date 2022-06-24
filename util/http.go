package util

import (
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/tool"
	"net"
)

var HttpTool *tool.Http

var httpTools map[net.Addr]*tool.Http

func init() {
	if global.Config.Settings.Basic.Interfaces == "" {
		HttpTool = tool.NewHttpTool(tool.GenHttpClient(&tool.HttpClientOptions{
			Transport: tool.GenHttpTransport(&tool.HttpTransportOptions{
				Timeout:           global.Timeout,
				SkipSslCertVerify: global.Config.Settings.Basic.SkipCertVerify,
			}),
			Timeout: global.Timeout,
		}))
	} else {
		httpTools = make(map[net.Addr]*tool.Http, 0)
	}
}

func HttpTools(addr net.Addr) *tool.Http {
	if HttpTool != nil {
		return HttpTool
	}
	if transport, ok := httpTools[addr]; ok {
		return transport
	} else {
		transport = tool.NewHttpTool(tool.GenHttpClient(&tool.HttpClientOptions{
			Transport: tool.GenHttpTransport(&tool.HttpTransportOptions{
				Timeout:           global.Timeout,
				LocalAddr:         addr,
				SkipSslCertVerify: global.Config.Settings.Basic.SkipCertVerify,
			}),
			Timeout: global.Timeout,
		}))
		httpTools[addr] = transport
		return transport
	}
}
