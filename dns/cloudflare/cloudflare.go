package cloudflare

import (
	dnsUtil "github.com/Mmx233/BitSrunLoginGo/dns/util"
	"github.com/cloudflare/cloudflare-go"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type DnsProvider struct {
	Api   *cloudflare.API `mapstructure:"-"`
	TTL   uint            `mapstructure:"-"`
	Zone  string          `mapstructure:"zone"`
	Email string          `mapstructure:"email"`
	Token string          `mapstructure:"token"`
}

func New(ttl uint, conf map[string]interface{}, Http *http.Client) (*DnsProvider, error) {
	var p = DnsProvider{
		TTL: ttl,
	}
	e := dnsUtil.DecodeConfig(conf, &p)
	if e != nil {
		return nil, e
	}
	log.Debugln("cloudflare dns provider:", &p)

	p.Api, e = cloudflare.New(p.Token, p.Email, cloudflare.HTTPClient(Http))
	return &p, e
}

func (a DnsProvider) SetDomainRecord(domain, ip string) error {

}
