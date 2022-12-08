package models

type DnsProvider interface {
	SetDomainRecord(domain, ip string) error
}
