package dns

import (
	"errors"
	"fmt"
	"github.com/Mmx233/BitSrunLoginGo/internal/pkg/dns/aliyun"
	"github.com/Mmx233/BitSrunLoginGo/internal/pkg/dns/cloudflare"
	"github.com/Mmx233/BitSrunLoginGo/internal/pkg/dns/dnspod"
	log "github.com/sirupsen/logrus"
)

func Run(c *Config) error {
	log.Debugf("开始 %s DDNS 流程", c.Provider)

	if c.TTL == 0 {
		c.TTL = 600
	}

	// 配置解析

	var dns Provider
	var err error
	switch c.Provider {
	case "aliyun":
		dns, err = aliyun.New(c.TTL, c.Conf, c.Http)
	case "cloudflare":
		dns, err = cloudflare.New(int(c.TTL), c.Conf, c.Http)
	case "dnspod":
		dns, err = dnspod.New(uint64(c.TTL), c.Conf, c.Http.Transport)
	default:
		var msg string
		if c.Provider == "" {
			msg = "DDNS 模块 dns 运营商不能为空"
		} else {
			msg = fmt.Sprintf("DDNS 模块 dns 运营商 %s 不支持", c.Provider)
		}
		log.Warnln(msg)
		return errors.New(msg)
	}
	if err != nil {
		log.Warnf("解析 DDNS config 失败：%v", err)
		return err
	}

	// 修改 dns 记录

	if err = dns.SetDomainRecord(c.Domain, c.IP); err != nil {
		log.Warnf("设置 dns 解析记录失败：%v", err)
		return err
	}

	log.Debugf("DDNS 配置应用成功: %s | %s", c.Domain, c.IP)

	return nil
}
