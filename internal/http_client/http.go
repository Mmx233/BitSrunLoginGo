package http_client

import (
	"context"
	"net"
	"net/http"

	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"github.com/Mmx233/BitSrunLoginGo/tools"
)

// newResolver 选择 DNS 解析器，优先使用配置 dns_server
func newResolver() *net.Resolver {
	if server := config.Settings.Basic.DNSServer; server != "" {
		return resolverFromServers([]string{ensureDNSPort(server)})
	}
	return platformResolver()
}

// resolverFromServers 构造将查询定向到指定 DNS 服务器的解析器，
// 按顺序尝试直到某个可用。调用方需保证 servers 非空。
func resolverFromServers(servers []string) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
			dialer := net.Dialer{Timeout: config.Timeout}
			var err error
			for _, server := range servers {
				var conn net.Conn
				if conn, err = dialer.DialContext(ctx, network, server); err == nil {
					return conn, nil
				}
			}
			return nil, err
		},
	}
}

// ensureDNSPort 为缺少端口的地址补全默认 DNS 端口 53。
func ensureDNSPort(addr string) string {
	if _, _, err := net.SplitHostPort(addr); err != nil {
		return net.JoinHostPort(addr, "53")
	}
	return addr
}

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
