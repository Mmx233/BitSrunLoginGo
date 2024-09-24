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
	if c.Logger == nil {
		c.Logger = log.New()
	}
	if c.TTL == 0 {
		c.TTL = 600
	}

	c.Logger.Infof("开始 %s DDNS 流程", c.Provider)

	var dns Provider
	var err error
	switch c.Provider {
	case "aliyun":
		dns, err = aliyun.New(c.TTL, c.Conf.Aliyun, c.Http)
	case "cloudflare":
		dns, err = cloudflare.New(int(c.TTL), c.Conf.Cloudflare, c.Http)
	case "dnspod":
		dns, err = dnspod.New(uint64(c.TTL), c.Conf.DnsPod, c.Http.Transport)
	default:
		var msg string
		if c.Provider == "" {
			msg = "DDNS 模块 dns 运营商不能为空"
		} else {
			msg = fmt.Sprintf("DDNS 模块 dns 运营商 %s 不支持", c.Provider)
		}
		c.Logger.Warnln(msg)
		return errors.New(msg)
	}
	if err != nil {
		c.Logger.Warnf("解析 DDNS config 失败：%v", err)
		return err
	}

	// 修改 dns 记录

	if err = dns.SetDomainRecord(c.Domain, c.IP); err != nil {
		c.Logger.Warnf("设置 dns 解析记录失败：%v", err)
		return err
	}

	c.Logger.Infof("DDNS 配置应用成功: %s | %s", c.Domain, c.IP)

	return nil
}
