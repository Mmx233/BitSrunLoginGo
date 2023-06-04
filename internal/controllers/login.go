package controllers

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/global"
	"github.com/Mmx233/BitSrunLoginGo/internal/pkg/dns"
	"github.com/Mmx233/BitSrunLoginGo/pkg/srun"
	"github.com/Mmx233/BitSrunLoginGo/tools"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Login 登录逻辑
func Login(eth *tools.Eth, debugOutput bool) error {
	// 登录配置初始化
	httpClient := tools.HttpPackSelect(eth).Client
	conf := &srun.Conf{
		Https: global.Config.Settings.Basic.Https,
		LoginInfo: srun.LoginInfo{
			Form: global.Config.Form,
			Meta: global.Config.Meta,
		},
		Client: httpClient,
	}

	// 选择输出函数
	var output func(args ...interface{})
	if debugOutput {
		output = log.Debugln
	} else {
		output = log.Infoln
	}

	output("正在获取登录状态")

	online, ip, e := srun.LoginStatus(conf)
	if e != nil {
		return e
	}

	log.Debugln("认证客户端 ip: ", ip)

	// 登录执行

	if online {
		output("已登录~")

		if global.Config.Settings.DDNS.Enable && global.Config.Settings.Guardian.Enable && ipLast != ip {
			if ddns(ip, httpClient) == nil {
				ipLast = ip
			}
		}

		return nil
	} else {
		log.Infoln("检测到用户未登录，开始尝试登录...")

		if e = srun.DoLogin(ip, conf); e != nil {
			return e
		}

		log.Infoln("登录成功~")

		if global.Config.Settings.DDNS.Enable {
			_ = ddns(ip, httpClient)
		}
	}

	return nil
}

var ipLast string

func ddns(ip string, httpClient *http.Client) error {
	return dns.Run(&dns.Config{
		Provider: global.Config.Settings.DDNS.Provider,
		IP:       ip,
		Domain:   global.Config.Settings.DDNS.Domain,
		TTL:      global.Config.Settings.DDNS.TTL,
		Conf:     global.Config.Settings.DDNS.Config,
		Http:     httpClient,
	})
}
