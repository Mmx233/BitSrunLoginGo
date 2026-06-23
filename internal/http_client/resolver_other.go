//go:build !android

package http_client

import "net"

func platformResolver() *net.Resolver {
	return nil
}
