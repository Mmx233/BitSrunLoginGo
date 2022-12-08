package dns

import (
	"github.com/Mmx233/BitSrunLoginGo/dns/aliyun"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

func Run(c *Config) error {
	var meta BasicConfig
	e := mapstructure.Decode(c.Conf, &meta)
	if e != nil {
		log.Warnf("解析 DDNS 配置失败：%v", e)
		return e
	}

	if meta.TTL == 0 {
		meta.TTL = 600
	}

	// 配置解析

	var dns Provider
	switch c.Provider {
	case "aliyun":
		dns, e = aliyun.New(meta.TTL, meta.Other, c.Http)
	default:
		log.Warnf("DDNS 模块 dns 运营商 %s 不支持", c.Provider)
		return nil
	}
	if e != nil {
		log.Warnf("解析 DDNS 配置失败：%v", e)
		return e
	}

	// 修改 dns 记录

	if e = dns.SetDomainRecord(meta.Domain, c.IP); e != nil {
		log.Warnf("设置 dns 解析记录失败：%v", e)
		return e
	}

	return nil
}
