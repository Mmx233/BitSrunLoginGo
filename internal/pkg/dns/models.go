package dns

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Provider interface {
	SetDomainRecord(domain, ip string) error
}

type Config struct {
	Provider string
	IP       string
	Domain   string
	TTL      uint
	Conf     config.DdnsProviderConfigSum
	Http     *http.Client
	Logger   *log.Logger
}
