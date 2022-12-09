package dns

import "net/http"

type Provider interface {
	SetDomainRecord(domain, ip string) error
}

type Config struct {
	Provider string
	IP       string
	Domain   string
	TTL      uint
	Conf     map[string]interface{}
	Http     *http.Client
}
