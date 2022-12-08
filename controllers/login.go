package controllers

import (
	"github.com/Mmx233/BitSrunLoginGo/dns"
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/BitSrunLoginGo/util"
	BitSrun "github.com/Mmx233/BitSrunLoginGo/v1"
	log "github.com/sirupsen/logrus"
	"net"
)

// Login 登录逻辑
func Login(localAddr net.Addr, debugOutput bool) error {
	// 登录状态检查

	httpClient := util.HttpPackSelect(localAddr).Client
	conf := &BitSrun.Conf{
		Https: global.Config.Settings.Basic.Https,
		LoginInfo: BitSrun.LoginInfo{
			Form: &global.Config.Form,
			Meta: &global.Config.Meta,
		},
		Client: httpClient,
	}

	var output func(args ...interface{})
	if debugOutput {
		output = log.Debugln
	} else {
		output = log.Infoln
	}

	output("正在获取登录状态")

	online, ip, e := BitSrun.LoginStatus(conf)
	if e != nil {
		return e
	}

	if localAddr != nil && global.Config.Settings.Basic.UseDhcpIP {
		ip = localAddr.(*net.TCPAddr).IP.String()
	} else if global.Flags.ClientIP != "" {
		ip = global.Flags.ClientIP
	}

	log.Debugln("认证客户端 ip: ", ip)

	// 登录执行

	if online {
		output("已登录~")

		return nil
	} else {
		log.Infoln("检测到用户未登录，开始尝试登录...")

		if e = BitSrun.DoLogin(ip, conf); e != nil {
			return e
		}

		log.Infoln("登录成功~")
	}

	// DDNS
	enable, _ := global.Config.Settings.DDNS["enable"]
	if open, _ := enable.(bool); open {
		log.Debugln("开始 DDNS 设置流程")

		provider, _ := global.Config.Settings.DDNS["provider"]
		providerStr, _ := provider.(string)
		if providerStr == "" {
			log.Warnln("DDNS 模块 dns 运营商不能为空")
			return nil
		}

		delete(global.Config.Settings.DDNS, "provider")

		_ = dns.Run(&dns.Config{
			Provider: providerStr,
			IP:       ip,
			Conf:     global.Config.Settings.DDNS,
			Http:     httpClient,
		})
	}

	return nil
}
