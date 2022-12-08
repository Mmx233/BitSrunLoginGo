package dns

import (
	"github.com/Mmx233/BitSrunLoginGo/dns/aliyun"
	"github.com/Mmx233/BitSrunLoginGo/dns/models"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

func Run(provider, ip string, conf map[string]interface{}) error {
	var meta models.BasicConfig
	e := mapstructure.Decode(conf, &meta)
	if e != nil {
		log.Warnf("解析 DDNS 配置失败：%v", e)
		return e
	}

	// 配置解析

	var dns models.DnsProvider
	switch provider {
	case "aliyun":
		dns, e = aliyun.New(meta.Other)
	default:
		log.Warnf("DDNS 模块 dns 运营商 %s 不支持", provider)
		return nil
	}
	if e != nil {
		log.Warnf("解析 DDNS 配置失败：%v", e)
		return e
	}

	// 修改 dns 记录

	if e = dns.SetDomainRecord(meta.Domain, ip); e != nil {
		log.Warnf("设置 dns 解析记录失败：%v", e)
		return e
	}

	return nil
}
