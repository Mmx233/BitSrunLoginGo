//go:build android

package http_client

import (
	"bufio"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// platformResolver 在 Android/Termux 上读取 Termux 的 resolv.conf 构造解析器。
// 静态构建（CGO disabled）的 Go 纯解析器默认读取 /etc/resolv.conf，
// 其在 Android 上通常指向不可达的 [::1]:53，导致 DNS 解析失败。
// Termux 真正的 resolv.conf 位于 $PREFIX/etc/resolv.conf。
func platformResolver() *net.Resolver {
	prefix := os.Getenv("PREFIX")
	if prefix == "" {
		return nil
	}
	servers := parseResolvConf(filepath.Join(prefix, "etc", "resolv.conf"))
	if len(servers) == 0 {
		return nil
	}
	return resolverFromServers(servers)
}

// parseResolvConf 解析 resolv.conf，提取 nameserver 条目并补全端口。
func parseResolvConf(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var servers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 && fields[0] == "nameserver" {
			servers = append(servers, ensureDNSPort(fields[1]))
		}
	}
	return servers
}
