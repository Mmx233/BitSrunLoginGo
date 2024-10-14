package dnspod

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/dns/util"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"net/http"
	"strings"
)

type DnsPod struct {
	SecretId  string `json:"secret_id,omitempty" yaml:"secret_id,omitempty"`
	SecretKey string `json:"secret_key,omitempty" yaml:"secret_key,omitempty"`
}

type DnsProvider struct {
	Client *dnspod.Client
	TTL    uint64
	DnsPod
}

func New(ttl uint64, conf DnsPod, Http http.RoundTripper) (*DnsProvider, error) {
	var p = DnsProvider{TTL: ttl, DnsPod: conf}
	var err error
	p.Client, err = dnspod.NewClient(common.NewCredential(p.SecretId, p.SecretKey), regions.Guangzhou, profile.NewClientProfile())
	p.Client.WithHttpTransport(Http)
	return &p, err
}

func (a DnsProvider) SetDomainRecord(domain, ip string) error {
	subDomain, rootDomain, err := dnsUtil.DecodeDomain(domain)
	if err != nil {
		return err
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
	res, err := a.Client.DescribeRecordList(reqRecordList)
	if (err != nil && strings.Contains(err.Error(), dnspod.RESOURCENOTFOUND_NODATAOFRECORD)) || (err == nil && len(res.Response.RecordList) == 0) {
		reqNewRecord := dnspod.NewCreateRecordRequest()
		reqNewRecord.TTL = &a.TTL
		reqNewRecord.Domain = &rootDomain
		reqNewRecord.RecordType = &recordType
		reqNewRecord.RecordLine = &recordLine
		reqNewRecord.Value = &ip
		reqNewRecord.SubDomain = &subDomain
		_, err = a.Client.CreateRecord(reqNewRecord)
		return err
	} else if err != nil {
		return err
	}

	reqModifyRecord := dnspod.NewModifyRecordRequest()
	reqModifyRecord.Domain = &rootDomain
	reqModifyRecord.SubDomain = &subDomain
	reqModifyRecord.Value = &ip
	reqModifyRecord.RecordId = res.Response.RecordList[0].RecordId
	reqModifyRecord.RecordLine = &recordLine
	reqModifyRecord.RecordType = &recordType
	_, err = a.Client.ModifyRecord(reqModifyRecord)
	return err
}
