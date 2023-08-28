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
	err := dnsUtil.DecodeConfig(conf, &p)
	if err != nil {
		return nil, err
	}

	if p.Zone == "" {
		return nil, errors.New("cloudflare zone 不能为空")
	}
	p.ZoneResource = cloudflare.ZoneIdentifier(p.Zone)

	if p.Token == "" {
		return nil, errors.New("cloudflare token 不能为空")
	}

	p.Api, err = cloudflare.NewWithAPIToken(p.Token, cloudflare.HTTPClient(Http))
	return &p, err
}

func (a DnsProvider) SetDomainRecord(domain, ip string) error {
	records, _, err := a.Api.ListDNSRecords(context.Background(), a.ZoneResource, cloudflare.ListDNSRecordsParams{
		Type: "A",
		Name: domain,
	})
	if err != nil {
		return err
	}

	if len(records) == 0 {
		_, err = a.Api.CreateDNSRecord(context.Background(), a.ZoneResource, cloudflare.CreateDNSRecordParams{
			Type:    "A",
			Name:    domain,
			Content: ip,
			TTL:     a.TTL,
		})
		return err
	} else {
		record := records[0]
		if record.Content == ip {
			return nil
		}
		_, err = a.Api.UpdateDNSRecord(context.Background(), a.ZoneResource, cloudflare.UpdateDNSRecordParams{
			ID:      record.ID,
			Content: ip,
		})
		return err
	}
}
