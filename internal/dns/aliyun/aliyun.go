package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Mmx233/BitSrunLoginGo/internal/dns/util"
	"github.com/Mmx233/tool"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"time"
)

type Aliyun struct {
	AccessKeyId     string `json:"access_key_id,omitempty" yaml:"access_key_id,omitempty"`
	AccessKeySecret string `json:"access_key_secret,omitempty" yaml:"access_key_secret,omitempty"`
}

type DnsProvider struct {
	TTL  uint
	Http *tool.Http
	Aliyun
}

func New(ttl uint, conf Aliyun, Http *http.Client) (*DnsProvider, error) {
	if conf.AccessKeyId == "" || conf.AccessKeySecret == "" {
		return nil, errors.New("aliyun AccessKey 不能为空")
	}
	return &DnsProvider{
		TTL:    ttl,
		Http:   tool.NewHttpTool(Http),
		Aliyun: conf,
	}, nil
}

func (a DnsProvider) SendRequest(Type, Action string, data map[string]interface{}) (*http.Response, error) {
	var reqOpt = tool.DoHttpReq{
		Url: "https://alidns.aliyuncs.com",
	}
	data["Format"] = "json"
	data["Version"] = "2015-01-09"
	data["SignatureMethod"] = "HMAC-SHA1"
	data["SignatureVersion"] = "1.0"
	data["SignatureNonce"] = fmt.Sprint(tool.RandMath(rand.NewSource(time.Now().UnixNano())).Num(10000000, 90000000))
	data["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	data["Action"] = Action
	data["AccessKeyId"] = a.AccessKeyId

	signStr := Type + "&" + url.QueryEscape("/") + "&"
	var keys = make([]string, len(data))
	var i int
	for k := range data {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for i, k := range keys {
		str := k + "=" + url.QueryEscape(fmt.Sprint(data[k]))
		if i == 0 {
			str = url.QueryEscape(str)
		} else {
			str = url.QueryEscape("&" + str)
		}
		signStr += str
	}

	mac := hmac.New(sha1.New, []byte(a.AccessKeySecret+"&"))
	_, err := mac.Write([]byte(signStr))
	if err != nil {
		return nil, err
	}
	data["Signature"] = base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if Type == "GET" || Type == "DELETE" {
		reqOpt.Query = data
	} else {
		reqOpt.Body = data
	}

	resp, err := a.Http.Request(Type, &reqOpt)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		defer resp.Body.Close()
		var res Response
		if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return nil, err
		}

		return nil, errors.New(res.Message)
	}

	return resp, nil
}

func (a DnsProvider) DomainRecordStatus(subDomain, rootDomain string) (*DomainStatus, bool, error) {
	resp, err := a.SendRequest("GET", "DescribeDomainRecords", map[string]interface{}{
		"DomainName": rootDomain,
		"SearchMode": "EXACT",
		"KeyWord":    subDomain,
		"PageSize":   1,
		"Type":       "A",
	})
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	var res DomainStatusRes
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, false, err
	}

	if res.TotalCount == 0 || len(res.DomainRecords.Record) == 0 {
		return nil, false, nil
	}

	return &res.DomainRecords.Record[0], true, nil
}

func (a DnsProvider) UpdateRecord(RecordId, subDomain, ip string) error {
	resp, err := a.SendRequest("POST", "UpdateDomainRecord", map[string]interface{}{
		"RecordId": RecordId,
		"RR":       subDomain,
		"Type":     "A",
		"Value":    ip,
		"TTL":      a.TTL,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (a DnsProvider) NewRecord(subDomain, rootDomain, ip string) error {
	resp, err := a.SendRequest("POST", "AddDomainRecord", map[string]interface{}{
		"DomainName": rootDomain,
		"RR":         subDomain,
		"Type":       "A",
		"Value":      ip,
		"TTL":        a.TTL,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (a DnsProvider) SetDomainRecord(domain, ip string) error {
	subDomain, rootDomain, err := dnsUtil.DecodeDomain(domain)
	if err != nil {
		return err
	}

	record, exist, err := a.DomainRecordStatus(subDomain, rootDomain)
	if err != nil {
		return err
	}

	if exist {
		if record.Value == ip {
			return nil
		}
		return a.UpdateRecord(record.RecordId, subDomain, ip)
	} else {
		return a.NewRecord(subDomain, rootDomain, ip)
	}
}
