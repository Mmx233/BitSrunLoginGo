package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	dnsUtil "github.com/Mmx233/BitSrunLoginGo/dns/util"
	"github.com/Mmx233/tool"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"time"
)

type DnsProvider struct {
	TTL             uint       `mapstructure:"-"`
	Http            *tool.Http `mapstructure:"-"`
	AccessKeyId     string     `mapstructure:"access_key_id"`
	AccessKeySecret string     `mapstructure:"access_key_secret"`
}

func New(ttl uint, conf map[string]interface{}, Http *http.Client) (*DnsProvider, error) {
	var p = DnsProvider{
		Http: tool.NewHttpTool(Http),
	}
	return &p, dnsUtil.DecodeConfig(conf, &p)
}

func (a DnsProvider) SendRequest(Type, Action string, data map[string]interface{}) (*http.Response, error) {
	var reqOpt = tool.DoHttpReq{
		Url: "https://alidns.aliyuncs.com",
	}
	rand.Seed(time.Now().UnixNano())
	data["Format"] = "json"
	data["Version"] = "2015-01-09"
	data["SignatureMethod"] = "HMAC-SHA1"
	data["SignatureNonce"] = fmt.Sprint(tool.Rand.Num(10000000, 90000000))
	data["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	data["Action"] = Action

	signStr := Type + "&" + url.QueryEscape("/")
	var keys = make([]string, len(data))
	var i int
	for k := range data {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		signStr += "&" + k + "=" + url.QueryEscape(fmt.Sprint(data[k]))
	}

	mac := hmac.New(sha1.New, []byte(a.AccessKeySecret+"&"))
	_, e := mac.Write([]byte(signStr))
	if e != nil {
		return nil, e
	}
	data["Signature"] = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%x", mac.Sum(nil))))

	if Type == "GET" || Type == "DELETE" {
		reqOpt.Query = data
	} else {
		reqOpt.Body = data
	}

	resp, e := a.Http.Request(Type, &reqOpt)
	if e != nil {
		return nil, e
	}

	if resp.StatusCode > 299 {
		defer resp.Body.Close()
		var res Response
		if e = json.NewDecoder(resp.Body).Decode(&res); e != nil {
			return nil, e
		}

		return nil, errors.New(res.Message)
	}

	return resp, nil
}

func (a DnsProvider) DomainRecordStatus(subDomain, rootDomain string) (*DomainStatus, bool, error) {
	resp, e := a.SendRequest("GET", "DescribeDomainRecords", map[string]interface{}{
		"DomainName": rootDomain,
		"SearchMode": "EXACT",
		"KeyWord":    subDomain,
		"PageSize":   1,
		"Type":       "A",
	})
	if e != nil {
		return nil, false, e
	}
	defer resp.Body.Close()

	var res DomainStatusRes
	if e = json.NewDecoder(resp.Body).Decode(&res); e != nil {
		return nil, false, e
	}

	if res.TotalCount == 0 || len(res.DomainRecords.Record) == 0 {
		return nil, false, nil
	}

	return &res.DomainRecords.Record[0], true, nil
}

func (a DnsProvider) UpdateRecord(RecordId, subDomain, ip string) error {
	resp, e := a.SendRequest("POST", "UpdateDomainRecord", map[string]interface{}{
		"RecordId": RecordId,
		"RR":       subDomain,
		"Type":     "A",
		"Value":    ip,
		"TTL":      a.TTL,
	})
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	return nil
}

func (a DnsProvider) NewRecord(subDomain, rootDomain, ip string) error {
	resp, e := a.SendRequest("POST", "AddDomainRecord", map[string]interface{}{
		"DomainName": rootDomain,
		"RR":         subDomain,
		"Type":       "A",
		"Value":      ip,
		"TTL":        a.TTL,
	})
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	return nil
}

func (a DnsProvider) SetDomainRecord(domain, ip string) error {
	subDomain, rootDomain, e := dnsUtil.DecodeDomain(domain)
	if e != nil {
		return e
	}

	record, exist, e := a.DomainRecordStatus(subDomain, rootDomain)
	if e != nil {
		return e
	}

	if exist {
		return a.UpdateRecord(record.RecordId, subDomain, ip)
	} else {
		return a.NewRecord(subDomain, rootDomain, ip)
	}
}
