package cloudflare

import (
	"context"
	"errors"
	"github.com/cloudflare/cloudflare-go"
	"net/http"
)

type Cloudflare struct {
	Zone  string `json:"zone" yaml:"zone"`
	Token string `json:"token" yaml:"token"`
}

type DnsProvider struct {
	Api          *cloudflare.API
	TTL          int
	ZoneResource *cloudflare.ResourceContainer
	Cloudflare
}

func New(ttl int, conf Cloudflare, Http *http.Client) (*DnsProvider, error) {
	var p = DnsProvider{
		TTL:        ttl,
		Cloudflare: conf,
	}
	if p.Zone == "" {
		return nil, errors.New("cloudflare zone 不能为空")
	}
	if p.Token == "" {
		return nil, errors.New("cloudflare token 不能为空")
	}
	p.ZoneResource = cloudflare.ZoneIdentifier(p.Zone)
	var err error
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
