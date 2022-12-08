package dns

import "net/http"

type Provider interface {
	SetDomainRecord(domain, ip string) error
}

type Config struct {
	Provider string
	IP       string
	Conf     map[string]interface{}
	Http     *http.Client
}

type BasicConfig struct {
	Domain string                 `mapstructure:"domain"`
	Other  map[string]interface{} `mapstructure:",remain"`
}
