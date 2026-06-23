package http_client

import (
	"net/http"

	"github.com/Mmx233/BitSrunLoginGo/tools"
)

// _DefaultClient 无网卡绑定的客户端，用于未指定网卡的情况
var _DefaultClient *http.Client

// ethClient 记录某个网卡当前绑定的 IP 及对应客户端
type ethClient struct {
	addr   string
	client *http.Client
}

// _EthClientMap 以网卡名为 key 缓存绑定了具体网卡的客户端。
// 每个网卡仅保留一个条目，IP 变化时重建，避免继续绑定已失效的旧地址。
var _EthClientMap map[string]*ethClient

func init() {
	_DefaultClient = CreateClientFromEth(nil)
	_EthClientMap = make(map[string]*ethClient)
}

func ClientSelect(eth *tools.Eth) *http.Client {
	if eth == nil {
		return _DefaultClient
	}
	addr := eth.Addr.String()
	if c, ok := _EthClientMap[eth.Name]; ok {
		if c.addr == addr {
			return c.client
		}
		// 网卡 IP 已变化，关闭旧连接并重建客户端
		c.client.CloseIdleConnections()
	}
	client := CreateClientFromEth(eth)
	_EthClientMap[eth.Name] = &ethClient{addr: addr, client: client}
	return client
}
