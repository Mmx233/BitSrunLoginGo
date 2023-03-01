package cloudflare

import (
	"context"
	"errors"
	"github.com/Mmx233/BitSrunLoginGo/internal/pkg/dns/util"
	"github.com/cloudflare/cloudflare-go"
	"net/http"
)

type DnsProvider struct {
	Api          *cloudflare.API               `mapstructure:"-"`
	TTL          int                           `mapstructure:"-"`
	Zone         string                        `mapstructure:"zone"`
	ZoneResource *cloudflare.ResourceContainer `mapstructure:"-"`
	Token        string                        `mapstructure:"token"`
}

func New(ttl int, conf map[string]interface{}, Http *http.Client) (*DnsProvider, error) {
	var p = DnsProvider{
		TTL: ttl,
	}
	e := dnsUtil.DecodeConfig(conf, &p)
	if e != nil {
		return nil, e
	}

	if p.Zone == "" {
		return nil, errors.New("cloudflare zone 不能为空")
	}
	p.ZoneResource = cloudflare.ZoneIdentifier(p.Zone)

	if p.Token == "" {
		return nil, errors.New("cloudflare token 不能为空")
	}

	p.Api, e = cloudflare.NewWithAPIToken(p.Token, cloudflare.HTTPClient(Http))
	return &p, e
}

func (a DnsProvider) SetDomainRecord(domain, ip string) error {
	records, _, e := a.Api.ListDNSRecords(context.Background(), a.ZoneResource, cloudflare.ListDNSRecordsParams{
		Type: "A",
		Name: domain,
	})
	if e != nil {
		return e
	}

	if len(records) == 0 {
		_, e = a.Api.CreateDNSRecord(context.Background(), a.ZoneResource, cloudflare.CreateDNSRecordParams{
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
		return a.Api.UpdateDNSRecord(context.Background(), a.ZoneResource, cloudflare.UpdateDNSRecordParams{
			ID:      record.ID,
			Content: ip,
		})
	}
}
