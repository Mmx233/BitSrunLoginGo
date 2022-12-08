package cloudflare

import (
	"context"
	"errors"
	dnsUtil "github.com/Mmx233/BitSrunLoginGo/dns/util"
	"github.com/cloudflare/cloudflare-go"
	"net/http"
)

type DnsProvider struct {
	Api   *cloudflare.API `mapstructure:"-"`
	TTL   int             `mapstructure:"-"`
	Zone  string          `mapstructure:"zone"`
	Token string          `mapstructure:"token"`
}

func New(ttl uint, conf map[string]interface{}, Http *http.Client) (*DnsProvider, error) {
	var p = DnsProvider{
		TTL: int(ttl),
	}
	e := dnsUtil.DecodeConfig(conf, &p)
	if e != nil {
		return nil, e
	}

	if p.Zone == "" {
		return nil, errors.New("cloudflare zone 不能为空")
	} else if p.Token == "" {
		return nil, errors.New("cloudflare token 不能为空")
	}

	p.Api, e = cloudflare.NewWithAPIToken(p.Token, cloudflare.HTTPClient(Http))
	return &p, e
}

func (a DnsProvider) SetDomainRecord(domain, ip string) error {
	records, e := a.Api.DNSRecords(context.Background(), a.Zone, cloudflare.DNSRecord{
		Type: "A",
		Name: domain,
	})
	if e != nil {
		return e
	}

	if len(records) == 0 {
		_, e = a.Api.CreateDNSRecord(context.Background(), a.Zone, cloudflare.DNSRecord{
			Type:    "A",
			Name:    domain,
			Content: ip,
			TTL:     a.TTL,
		})
		return e
	} else {
		record := records[0]

		if record.Content == ip {
			return nil
		}
		record.Content = ip
		return a.Api.UpdateDNSRecord(context.Background(), a.Zone, record.ID, record)
	}
}
