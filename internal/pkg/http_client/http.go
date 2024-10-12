package http_client

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"github.com/Mmx233/BitSrunLoginGo/internal/config/flags"
	"github.com/Mmx233/BitSrunLoginGo/internal/config/keys"
	"github.com/Mmx233/BitSrunLoginGo/tools"
	"net"
	"net/http"
)

var _DefaultClient *http.Client

var _EthClientMap map[string]*http.Client

func init() {
	logger := config.Logger.WithField(keys.LogComponent, "init http")
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

		_DefaultClient = CreateClientFromEth(eth)
	} else {
		_EthClientMap = make(map[string]*http.Client)
	}
}

func ClientSelect(eth *tools.Eth) *http.Client {
	if _DefaultClient != nil {
		return _DefaultClient
	}
	if client, ok := _EthClientMap[eth.Name]; ok {
		return client
	} else {
		client = CreateClientFromEth(eth)
		_EthClientMap[eth.Name] = client
		return client
	}
}
