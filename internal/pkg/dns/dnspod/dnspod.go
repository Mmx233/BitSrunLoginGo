package dnspod

import (
	dnsUtil2 "github.com/Mmx233/BitSrunLoginGo/internal/pkg/dns/util"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"net/http"
	"strings"
)

type DnsProvider struct {
	Client    *dnspod.Client `mapstructure:"-"`
	TTL       uint64         `mapstructure:"-"`
	SecretId  string         `mapstructure:"secret_id"`
	SecretKey string         `mapstructure:"secret_key"`
}

func New(ttl uint64, conf map[string]interface{}, Http http.RoundTripper) (*DnsProvider, error) {
	var p = DnsProvider{TTL: ttl}
	e := dnsUtil2.DecodeConfig(conf, &p)
	if e != nil {
		return nil, e
	}
	p.Client, e = dnspod.NewClient(common.NewCredential(p.SecretId, p.SecretKey), regions.Guangzhou, profile.NewClientProfile())
	p.Client.WithHttpTransport(Http)
	return &p, e
}

func (a DnsProvider) SetDomainRecord(domain, ip string) error {
	subDomain, rootDomain, e := dnsUtil2.DecodeDomain(domain)
	if e != nil {
		return e
	}

	var (
		recordType        = "A"
		recordLine        = "默认"
		limit      uint64 = 1
	)

	reqRecordList := dnspod.NewDescribeRecordListRequest()
	reqRecordList.Domain = &rootDomain
	reqRecordList.Subdomain = &subDomain
	reqRecordList.Limit = &limit
	res, e := a.Client.DescribeRecordList(reqRecordList)
	if (e != nil && strings.Contains(e.Error(), dnspod.RESOURCENOTFOUND_NODATAOFRECORD)) || (e == nil && len(res.Response.RecordList) == 0) {
		reqNewRecord := dnspod.NewCreateRecordRequest()
		reqNewRecord.TTL = &a.TTL
		reqNewRecord.Domain = &rootDomain
		reqNewRecord.RecordType = &recordType
		reqNewRecord.RecordLine = &recordLine
		reqNewRecord.Value = &ip
		reqNewRecord.SubDomain = &subDomain
		_, e = a.Client.CreateRecord(reqNewRecord)
		return e
	} else if e != nil {
		return e
	}

	reqModifyRecord := dnspod.NewModifyRecordRequest()
	reqModifyRecord.Domain = &rootDomain
	reqModifyRecord.SubDomain = &subDomain
	reqModifyRecord.Value = &ip
	reqModifyRecord.RecordId = res.Response.RecordList[0].RecordId
	reqModifyRecord.RecordLine = &recordLine
	reqModifyRecord.RecordType = &recordType
	_, e = a.Client.ModifyRecord(reqModifyRecord)
	return e
}
