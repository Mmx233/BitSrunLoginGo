package http_client

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"github.com/Mmx233/BitSrunLoginGo/internal/config/flags"
	"github.com/Mmx233/BitSrunLoginGo/tools"
	"net"
	"net/http"
)

type Http struct {
	Client *http.Client
}

var HttpPack *Http

var httpTools map[string]*Http

func init() {
	logger := config.Logger
	if config.Settings.Basic.Interfaces == "" {
		var eth *tools.Eth
		if flags.Interface != "" {
			netEth, err := net.InterfaceByName(flags.Interface)
			if err != nil {
				logger.Warnf("获取指定网卡 %s 失败，使用默认网卡: %v", flags.Interface, err)
			} else {
				eth, err = tools.ConvertInterface(logger, *netEth)
				if err != nil {
					logger.Warnf("获取指定网卡 %s ip 地址失败，使用默认网卡: %v", flags.Interface, err)
				} else if eth == nil {
					logger.Warnf("指定网卡 %s 无可用 ip 地址，使用默认网卡", flags.Interface)
				} else {
					logger.Debugf("使用指定网卡 %s ip: %s", eth.Name, eth.Addr.String())
				}
			}
		}

		HttpPack = genHttpPack(eth)
	} else {
		httpTools = make(map[string]*Http)
	}
}

func HttpPackSelect(eth *tools.Eth) *Http {
	if HttpPack != nil {
		return HttpPack
	}
	if addrHttp, ok := httpTools[eth.Name]; ok {
		return addrHttp
	} else {
		addrHttp = genHttpPack(eth)
		httpTools[eth.Name] = addrHttp
		return addrHttp
	}
}
